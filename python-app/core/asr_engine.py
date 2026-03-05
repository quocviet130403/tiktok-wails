"""ASR Engine - Wrapper for faster-whisper transcription with chunked support.

Provides a simple interface to transcribe video/audio files
to ASRData with word-level timestamps.
Supports ChunkedASR for long videos (>10 minutes).
"""

import io
import logging
import subprocess
import tempfile
import threading
from concurrent.futures import ThreadPoolExecutor, as_completed
from pathlib import Path
from typing import Callable, List, Optional, Tuple

from core.asr_data import ASRData, ASRDataSeg
from core.chunk_merger import ChunkMerger

logger = logging.getLogger("asr_engine")

# Chunked ASR constants
MS_PER_SECOND = 1000
DEFAULT_CHUNK_LENGTH_SEC = 10 * 60  # 10 minutes
DEFAULT_CHUNK_OVERLAP_SEC = 10  # 10 seconds overlap


def _extract_audio(video_path: str, audio_path: str) -> None:
    """Extract audio from video using ffmpeg."""
    cmd = [
        "ffmpeg", "-i", video_path,
        "-vn", "-acodec", "pcm_s16le",
        "-ar", "16000", "-ac", "1",
        "-y", audio_path,
    ]
    result = subprocess.run(cmd, capture_output=True, text=True)
    if result.returncode != 0:
        raise RuntimeError(f"FFmpeg audio extraction failed: {result.stderr}")


def _transcribe_audio(
    audio_path: str,
    model_size: str = "large-v3-turbo",
    language: str = "zh",
    device: str = "auto",
    compute_type: str = "auto",
    vad_filter: bool = True,
    word_timestamps: bool = True,
) -> ASRData:
    """Transcribe a single audio file (no chunking)."""
    from faster_whisper import WhisperModel

    # Auto-detect compute type and device
    if device == "auto":
        try:
            import torch
            device = "cuda" if torch.cuda.is_available() else "cpu"
        except ImportError:
            device = "cpu"

    if compute_type == "auto":
        compute_type = "float16" if device == "cuda" else "int8"
    if device == "cpu":
        compute_type = "int8"

    logger.info(f"Loading faster-whisper model: {model_size} on {device}")
    model = WhisperModel(model_size, device=device, compute_type=compute_type)

    segments_gen, info = model.transcribe(
        audio_path,
        language=language,
        vad_filter=vad_filter,
        word_timestamps=word_timestamps,
    )

    logger.info(f"Detected language: {info.language} (prob={info.language_probability:.2f})")

    asr_segments = []
    for segment in segments_gen:
        if word_timestamps and segment.words:
            for word in segment.words:
                start_ms = int(word.start * 1000)
                end_ms = int(word.end * 1000)
                text = word.word.strip()
                if text:
                    asr_segments.append(ASRDataSeg(text=text, start_time=start_ms, end_time=end_ms))
        else:
            start_ms = int(segment.start * 1000)
            end_ms = int(segment.end * 1000)
            text = segment.text.strip()
            if text:
                asr_segments.append(ASRDataSeg(text=text, start_time=start_ms, end_time=end_ms))

    logger.info(f"Transcription complete: {len(asr_segments)} segments")
    return ASRData(asr_segments)


def _get_audio_duration_ms(audio_path: str) -> int:
    """Get audio duration in milliseconds using ffprobe."""
    cmd = [
        "ffprobe", "-i", audio_path,
        "-show_entries", "format=duration",
        "-v", "quiet", "-of", "csv=p=0",
    ]
    try:
        result = subprocess.run(cmd, capture_output=True, text=True)
        return int(float(result.stdout.strip()) * 1000)
    except Exception:
        return 0


def _split_audio_chunks(
    audio_path: str,
    chunk_length_ms: int = DEFAULT_CHUNK_LENGTH_SEC * MS_PER_SECOND,
    overlap_ms: int = DEFAULT_CHUNK_OVERLAP_SEC * MS_PER_SECOND,
) -> List[Tuple[str, int]]:
    """Split audio into overlapping chunks using ffmpeg.

    Returns list of (chunk_file_path, offset_ms).
    """
    duration_ms = _get_audio_duration_ms(audio_path)
    if duration_ms <= 0 or duration_ms <= chunk_length_ms:
        return [(audio_path, 0)]

    chunks = []
    start_ms = 0
    idx = 0

    while start_ms < duration_ms:
        end_ms = min(start_ms + chunk_length_ms, duration_ms)
        chunk_path = tempfile.mktemp(suffix=f"_chunk{idx}.wav")

        cmd = [
            "ffmpeg",
            "-ss", f"{start_ms / 1000:.3f}",
            "-i", audio_path,
            "-t", f"{(end_ms - start_ms) / 1000:.3f}",
            "-acodec", "pcm_s16le", "-ar", "16000", "-ac", "1",
            "-y", chunk_path,
        ]
        subprocess.run(cmd, capture_output=True, text=True)
        chunks.append((chunk_path, start_ms))
        logger.info(f"Chunk {idx}: {start_ms/1000:.1f}s - {end_ms/1000:.1f}s")

        start_ms += chunk_length_ms - overlap_ms
        idx += 1
        if end_ms >= duration_ms:
            break

    return chunks


def transcribe_video(
    video_path: str,
    model_size: str = "large-v3-turbo",
    language: str = "zh",
    device: str = "auto",
    compute_type: str = "auto",
    vad_filter: bool = True,
    word_timestamps: bool = True,
    chunk_length_sec: int = DEFAULT_CHUNK_LENGTH_SEC,
) -> ASRData:
    """Transcribe video/audio file to ASRData with chunked support for long files.

    Args:
        video_path: Path to video or audio file
        model_size: Whisper model size
        language: Source language code
        device: Device (auto/cuda/cpu)
        compute_type: Compute type (auto/float16/int8)
        vad_filter: Enable VAD filter
        word_timestamps: Enable word-level timestamps
        chunk_length_sec: Chunk length for long audio (seconds)

    Returns:
        ASRData with transcribed segments
    """
    # Extract audio if input is video
    video_ext = Path(video_path).suffix.lower()
    audio_extensions = {".wav", ".mp3", ".flac", ".ogg", ".m4a", ".aac"}

    if video_ext not in audio_extensions:
        audio_path = tempfile.mktemp(suffix=".wav")
        logger.info("Extracting audio from video...")
        _extract_audio(video_path, audio_path)
    else:
        audio_path = video_path

    try:
        # Check duration to decide if chunking is needed
        duration_ms = _get_audio_duration_ms(audio_path)
        chunk_length_ms = chunk_length_sec * MS_PER_SECOND

        if duration_ms > 0 and duration_ms > chunk_length_ms:
            # Long audio - use chunked ASR
            logger.info(f"Long audio ({duration_ms/1000:.0f}s), using chunked ASR")
            chunks = _split_audio_chunks(audio_path, chunk_length_ms)

            if len(chunks) <= 1:
                return _transcribe_audio(
                    audio_path, model_size, language, device, compute_type,
                    vad_filter, word_timestamps,
                )

            # Transcribe each chunk
            chunk_results = []
            for i, (chunk_path, offset_ms) in enumerate(chunks):
                logger.info(f"Transcribing chunk {i+1}/{len(chunks)}")
                result = _transcribe_audio(
                    chunk_path, model_size, language, device, compute_type,
                    vad_filter, word_timestamps,
                )
                chunk_results.append(result)
                # Clean up chunk file
                if chunk_path != audio_path:
                    try:
                        Path(chunk_path).unlink(missing_ok=True)
                    except Exception:
                        pass

            # Merge chunks
            merger = ChunkMerger(min_match_count=2, fuzzy_threshold=0.7)
            chunk_offsets = [offset for _, offset in chunks]
            merged = merger.merge_chunks(
                chunk_results, chunk_offsets,
                overlap_duration=DEFAULT_CHUNK_OVERLAP_SEC * MS_PER_SECOND,
            )
            logger.info(f"Chunked ASR complete: {len(merged)} segments from {len(chunks)} chunks")
            return merged
        else:
            # Short audio - direct transcription
            return _transcribe_audio(
                audio_path, model_size, language, device, compute_type,
                vad_filter, word_timestamps,
            )
    finally:
        # Clean up temp audio
        if video_ext not in audio_extensions:
            try:
                Path(audio_path).unlink(missing_ok=True)
            except Exception:
                pass
