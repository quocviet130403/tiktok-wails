"""ASR Engine - Wrapper for faster-whisper transcription.

Provides a simple interface to transcribe video/audio files
to ASRData with word-level timestamps.
"""

import logging
import subprocess
import tempfile
from pathlib import Path

from core.asr_data import ASRData, ASRDataSeg

logger = logging.getLogger("asr_engine")


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


def transcribe_video(
    video_path: str,
    model_size: str = "large-v3-turbo",
    language: str = "zh",
    device: str = "auto",
    compute_type: str = "auto",
    vad_filter: bool = True,
    word_timestamps: bool = True,
) -> ASRData:
    """Transcribe video/audio file to ASRData.

    Args:
        video_path: Path to video or audio file
        model_size: Whisper model size (tiny, base, small, medium, large-v3, large-v3-turbo)
        language: Source language code (zh, en, ja, etc.)
        device: Device to use (auto, cuda, cpu)
        compute_type: Compute type (auto, float16, int8, etc.)
        vad_filter: Enable Voice Activity Detection filter
        word_timestamps: Enable word-level timestamps

    Returns:
        ASRData with transcribed segments
    """
    from faster_whisper import WhisperModel

    logger.info(f"Loading faster-whisper model: {model_size} on {device}")

    # Auto-detect compute type based on device
    if compute_type == "auto":
        if device == "cuda" or (device == "auto"):
            compute_type = "float16"
        else:
            compute_type = "int8"

    if device == "auto":
        try:
            import torch
            device = "cuda" if torch.cuda.is_available() else "cpu"
        except ImportError:
            device = "cpu"

    if device == "cpu":
        compute_type = "int8"

    model = WhisperModel(model_size, device=device, compute_type=compute_type)

    # Extract audio if input is video
    video_ext = Path(video_path).suffix.lower()
    audio_extensions = {".wav", ".mp3", ".flac", ".ogg", ".m4a", ".aac"}

    if video_ext not in audio_extensions:
        # Need to extract audio first
        with tempfile.NamedTemporaryFile(suffix=".wav", delete=False) as tmp:
            audio_path = tmp.name

        try:
            logger.info("Extracting audio from video...")
            _extract_audio(video_path, audio_path)
            transcribe_path = audio_path
        except Exception as e:
            # Try direct transcription as fallback
            logger.warning(f"Audio extraction failed, trying direct: {e}")
            transcribe_path = video_path
    else:
        transcribe_path = video_path

    logger.info(f"Transcribing with language={language}, vad_filter={vad_filter}")

    segments_gen, info = model.transcribe(
        transcribe_path,
        language=language,
        vad_filter=vad_filter,
        word_timestamps=word_timestamps,
    )

    logger.info(f"Detected language: {info.language} (prob={info.language_probability:.2f})")

    # Convert to ASRData
    asr_segments = []

    for segment in segments_gen:
        if word_timestamps and segment.words:
            # Word-level timestamps
            for word in segment.words:
                start_ms = int(word.start * 1000)
                end_ms = int(word.end * 1000)
                text = word.word.strip()
                if text:
                    asr_segments.append(ASRDataSeg(text=text, start_time=start_ms, end_time=end_ms))
        else:
            # Segment-level timestamps
            start_ms = int(segment.start * 1000)
            end_ms = int(segment.end * 1000)
            text = segment.text.strip()
            if text:
                asr_segments.append(ASRDataSeg(text=text, start_time=start_ms, end_time=end_ms))

    # Clean up temp audio file
    if video_ext not in audio_extensions:
        try:
            Path(audio_path).unlink(missing_ok=True)
        except Exception:
            pass

    logger.info(f"Transcription complete: {len(asr_segments)} segments")
    return ASRData(asr_segments)
