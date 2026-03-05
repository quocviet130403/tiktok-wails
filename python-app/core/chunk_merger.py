"""ChunkMerger - Merge overlapping ASR chunks.

Ported from VideoCaptioner. Uses sliding window algorithm
to find best alignment in overlap regions and merge at midpoint.
"""

import difflib
import logging
from typing import List, Optional, Tuple

from core.asr_data import ASRData, ASRDataSeg

logger = logging.getLogger("chunk_merger")


class ChunkMerger:
    """Merge overlapping ASR chunks using sliding window alignment."""

    def __init__(self, min_match_count: int = 2, fuzzy_threshold: float = 0.7):
        self.min_match_count = min_match_count
        self.fuzzy_threshold = fuzzy_threshold
        self._is_word_level = False

    def merge_chunks(
        self,
        chunks: List[ASRData],
        chunk_offsets: Optional[List[int]] = None,
        overlap_duration: int = 10000,
    ) -> ASRData:
        """Merge multiple ASR chunk results."""
        if not chunks:
            raise ValueError("chunks cannot be empty")

        if len(chunks) == 1:
            return chunks[0]

        self._is_word_level = any(chunk.is_word_timestamp() for chunk in chunks)

        if chunk_offsets is None:
            chunk_offsets = self._infer_chunk_offsets(chunks, overlap_duration)

        if len(chunks) != len(chunk_offsets):
            raise ValueError(f"chunks ({len(chunks)}) and offsets ({len(chunk_offsets)}) count mismatch")

        # Adjust timestamps to absolute time
        adjusted_chunks = [
            self._adjust_timestamps(chunk.segments, offset)
            for chunk, offset in zip(chunks, chunk_offsets)
        ]

        # Merge pairs sequentially
        merged_segments = adjusted_chunks[0]
        for i in range(1, len(adjusted_chunks)):
            merged_segments = self._merge_two_sequences(
                merged_segments, adjusted_chunks[i], overlap_duration
            )

        logger.info(f"Merge complete: {len(merged_segments)} segments")
        return ASRData(merged_segments)

    def _merge_two_sequences(
        self,
        left: List[ASRDataSeg],
        right: List[ASRDataSeg],
        overlap_duration: int,
    ) -> List[ASRDataSeg]:
        """Merge two segment sequences using sliding window."""
        if not left:
            return right
        if not right:
            return left

        left_len = len(left)
        left_overlap = self._extract_overlap(left, from_end=True, duration=overlap_duration)
        right_overlap = self._extract_overlap(right, from_end=False, duration=overlap_duration)

        if not left_overlap or not right_overlap:
            return left + right

        best_match = self._find_best_alignment(left_overlap, right_overlap)
        if best_match is None:
            # No valid match - use time boundary
            split_idx = left_len
            right_start = right[0].start_time
            for i in range(left_len - 1, -1, -1):
                if left[i].end_time <= right_start:
                    split_idx = i + 1
                    break
            return left[:split_idx] + right

        left_start_idx, left_end_idx, right_start_idx, right_end_idx, _ = best_match
        left_mid = (left_start_idx + left_end_idx) // 2
        right_mid = (right_start_idx + right_end_idx) // 2
        left_overlap_offset = left_len - len(left_overlap)
        left_cut = left_overlap_offset + left_mid

        return left[:left_cut] + right[right_mid:]

    def _find_best_alignment(
        self, left: List[ASRDataSeg], right: List[ASRDataSeg]
    ) -> Optional[Tuple[int, int, int, int, int]]:
        """Sliding window to find best alignment position."""
        left_len = len(left)
        right_len = len(right)
        best_score = 0.0
        best_result = None

        for i in range(1, left_len + right_len + 1):
            epsilon = float(i) / 10000.0
            left_start = max(0, left_len - i)
            left_end = min(left_len, left_len + right_len - i)
            right_start = max(0, i - left_len)
            right_end = min(right_len, i)

            left_slice = left[left_start:left_end]
            right_slice = right[right_start:right_end]

            if len(left_slice) != len(right_slice):
                continue

            if self._is_word_level:
                matches = sum(1 for l, r in zip(left_slice, right_slice) if l.text == r.text)
            else:
                matches = sum(
                    1 for l, r in zip(left_slice, right_slice)
                    if difflib.SequenceMatcher(None, l.text, r.text).ratio() > self.fuzzy_threshold
                )

            score = matches / float(i) + epsilon
            if matches >= self.min_match_count and score > best_score:
                best_score = score
                best_result = (left_start, left_end, right_start, right_end, matches)

        return best_result

    @staticmethod
    def _adjust_timestamps(segments: List[ASRDataSeg], offset: int) -> List[ASRDataSeg]:
        return [
            ASRDataSeg(
                text=seg.text,
                start_time=seg.start_time + offset,
                end_time=seg.end_time + offset,
                translated_text=seg.translated_text,
            )
            for seg in segments
        ]

    @staticmethod
    def _extract_overlap(
        segments: List[ASRDataSeg], from_end: bool, duration: int
    ) -> List[ASRDataSeg]:
        if not segments:
            return []
        overlap = []
        if from_end:
            threshold = segments[-1].end_time - duration
            for seg in reversed(segments):
                if seg.start_time >= threshold:
                    overlap.insert(0, seg)
                else:
                    break
        else:
            threshold = segments[0].start_time + duration
            for seg in segments:
                if seg.end_time <= threshold:
                    overlap.append(seg)
                else:
                    break
        return overlap

    @staticmethod
    def _infer_chunk_offsets(chunks: List[ASRData], overlap_duration: int) -> List[int]:
        offsets = [0]
        for i in range(1, len(chunks)):
            prev = chunks[i - 1]
            if prev.segments:
                prev_end = prev.segments[-1].end_time
                next_offset = offsets[-1] + prev_end - overlap_duration
                offsets.append(max(next_offset, offsets[-1]))
            else:
                offsets.append(offsets[-1])
        return offsets
