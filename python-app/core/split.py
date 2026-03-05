"""Subtitle splitter - Intelligent sentence splitting using LLM.

Ported from VideoCaptioner's split module. Pipeline:
1. Send continuous text to LLM → LLM inserts <br> at semantic breakpoints
2. Match LLM sentences back to word-level segments using difflib
3. Fall back to rule-based splitting if LLM fails
"""

import difflib
import logging
import re
from typing import List, Tuple, Union

from core.asr_data import ASRData, ASRDataSeg
from core.llm_client import call_llm
from core.prompts import get_prompt
from core.text_utils import count_words, is_mainly_cjk, is_pure_punctuation, is_space_separated_language

logger = logging.getLogger("subtitle_splitter")

# Configuration constants
MAX_WORD_COUNT_CJK = 25
MAX_WORD_COUNT_ENGLISH = 18
SEGMENT_WORD_THRESHOLD = 500
MAX_GAP = 1500  # ms
SPLIT_SEARCH_RANGE = 30
TIME_GAP_WINDOW_SIZE = 5
TIME_GAP_MULTIPLIER = 3
MIN_GROUP_SIZE = 5
RULE_SPLIT_GAP = 500  # ms
RULE_MIN_SEGMENT_SIZE = 4
PREFIX_WORD_RATIO = 0.6
SUFFIX_WORD_RATIO = 0.4
MATCH_SIMILARITY_THRESHOLD = 0.5
MATCH_MAX_SHIFT = 30
MATCH_MAX_UNMATCHED = 5
MATCH_LARGE_SHIFT = 100
LLM_MAX_STEPS = 2


def _preprocess_segments(segments: List[ASRDataSeg]) -> List[ASRDataSeg]:
    """Remove punctuation-only segments and normalize spacing."""
    new_segments = []
    for seg in segments:
        if not is_pure_punctuation(seg.text):
            text = seg.text.strip()
            if is_space_separated_language(text):
                seg.text = text + " "
            new_segments.append(seg)
    return new_segments


def _split_by_llm(
    text: str,
    model: str,
    max_word_count_cjk: int = 18,
    max_word_count_english: int = 12,
) -> List[str]:
    """Use LLM to split text into sentences with agent loop validation."""
    system_prompt = get_prompt(
        "split/sentence",
        max_word_count_cjk=max_word_count_cjk,
        max_word_count_english=max_word_count_english,
    )
    user_prompt = f"Please use multiple <br> tags to separate the following sentence:\n{text}"
    messages = [
        {"role": "system", "content": system_prompt},
        {"role": "user", "content": user_prompt},
    ]

    last_result = None
    for step in range(LLM_MAX_STEPS):
        response = call_llm(messages=messages, model=model, temperature=0.1)
        result_text = response.choices[0].message.content
        result_text_cleaned = re.sub(r"\n+", "", result_text)
        split_result = [s.strip() for s in result_text_cleaned.split("<br>") if s.strip()]
        last_result = split_result

        # Validate
        is_valid, error_message = _validate_split_result(
            text, split_result, max_word_count_cjk, max_word_count_english
        )
        if is_valid:
            return split_result

        logger.warning(f"Split validation failed (attempt {step + 1}): {error_message}")
        messages.append({"role": "assistant", "content": result_text})
        messages.append({
            "role": "user",
            "content": f"Error: {error_message}\nFix the errors and output the COMPLETE corrected text with <br> tags.",
        })

    return last_result if last_result else [text]


def _validate_split_result(
    original_text: str,
    split_result: List[str],
    max_word_count_cjk: int,
    max_word_count_english: int,
) -> Tuple[bool, str]:
    """Validate split result: content similarity and length limits."""
    if not split_result:
        return False, "No segments found."

    # Check content similarity
    original_cleaned = re.sub(r"\s+", " ", original_text)
    text_is_cjk = is_mainly_cjk(original_cleaned)
    merged_char = "" if text_is_cjk else " "
    merged = merged_char.join(split_result)
    merged_cleaned = re.sub(r"\s+", " ", merged)

    matcher = difflib.SequenceMatcher(None, original_cleaned, merged_cleaned)
    similarity = matcher.ratio()

    if similarity < 0.96:
        return False, f"Content modified (similarity: {similarity:.1%}). Keep original text unchanged, only insert <br>."

    # Check length limits
    violations = []
    for i, segment in enumerate(split_result, 1):
        word_count = count_words(segment)
        max_allowed = max_word_count_cjk if text_is_cjk else max_word_count_english
        if word_count > max_allowed:
            preview = segment[:40] + "..." if len(segment) > 40 else segment
            violations.append(f"Segment {i} '{preview}': {word_count} > {max_allowed} limit")

    if violations:
        return False, "Length violations:\n" + "\n".join(f"- {v}" for v in violations)

    return True, ""


class SubtitleSplitter:
    """Intelligent subtitle splitter using LLM with rule-based fallback."""

    def __init__(
        self,
        model: str,
        max_word_count_cjk: int = MAX_WORD_COUNT_CJK,
        max_word_count_english: int = MAX_WORD_COUNT_ENGLISH,
    ):
        self.model = model
        self.max_word_count_cjk = max_word_count_cjk
        self.max_word_count_english = max_word_count_english

    def split_subtitle(self, subtitle_data: Union[str, ASRData]) -> ASRData:
        """Main entry: split subtitle into proper sentences."""
        try:
            if isinstance(subtitle_data, str):
                asr_data = ASRData.from_subtitle_file(subtitle_data)
            else:
                asr_data = subtitle_data

            if not asr_data.is_word_timestamp():
                asr_data = asr_data.split_to_word_segments()

            asr_data.segments = _preprocess_segments(asr_data.segments)
            txt = asr_data.to_txt().replace("\n", "")

            total_word_count = count_words(txt)
            num_segments = max(1, total_word_count // SEGMENT_WORD_THRESHOLD + (1 if total_word_count % SEGMENT_WORD_THRESHOLD > 0 else 0))
            logger.info(f"Word count: {total_word_count}, splitting into {num_segments} segment(s)")

            asr_data_list = self._split_asr_data(asr_data, num_segments)
            processed_segments = []
            for part in asr_data_list:
                if not part.segments:
                    continue
                try:
                    result = self._process_by_llm(part.segments)
                except Exception as e:
                    logger.warning(f"LLM split failed, using rules: {e}")
                    result = self._process_by_rules(part.segments)
                processed_segments.extend(result)

            processed_segments.sort(key=lambda seg: seg.start_time)
            return ASRData(processed_segments)

        except Exception as e:
            logger.error(f"Split failed: {e}")
            raise RuntimeError(f"Split failed: {e}")

    def _split_asr_data(self, asr_data: ASRData, num_segments: int) -> List[ASRData]:
        """Split ASR data into chunks at time gaps."""
        total_segs = len(asr_data.segments)
        if num_segments <= 1 or total_segs <= num_segments:
            return [asr_data]

        total_word_count = count_words(asr_data.to_txt())
        words_per_segment = total_word_count // num_segments
        split_indices = [i * words_per_segment for i in range(1, num_segments)]

        adjusted = []
        for sp in split_indices:
            start = max(0, sp - SPLIT_SEARCH_RANGE)
            end = min(total_segs - 1, sp + SPLIT_SEARCH_RANGE)
            max_gap = -1
            best = sp
            for j in range(start, end):
                gap = asr_data.segments[j + 1].start_time - asr_data.segments[j].end_time
                if gap > max_gap:
                    max_gap = gap
                    best = j
            adjusted.append(best)

        adjusted = sorted(set(adjusted))
        segments = []
        prev = 0
        for idx in adjusted:
            segments.append(ASRData(asr_data.segments[prev:idx + 1]))
            prev = idx + 1
        if prev < total_segs:
            segments.append(ASRData(asr_data.segments[prev:]))
        return segments

    def _process_by_llm(self, segments: List[ASRDataSeg]) -> List[ASRDataSeg]:
        """Use LLM for intelligent splitting."""
        txt = "".join(seg.text for seg in segments)
        logger.info(f"LLM split: text length={count_words(txt)}")

        sentences = _split_by_llm(
            text=txt, model=self.model,
            max_word_count_cjk=self.max_word_count_cjk,
            max_word_count_english=self.max_word_count_english,
        )
        return self._merge_segments_based_on_sentences(segments, sentences)

    def _process_by_rules(self, segments: List[ASRDataSeg]) -> List[ASRDataSeg]:
        """Rule-based fallback splitting."""
        # Group by time gaps
        groups = self._group_by_time_gaps(segments, max_gap=RULE_SPLIT_GAP, check_large_gaps=True)

        # Split long groups at common words
        result_groups = []
        for group in groups:
            max_wc = self.max_word_count_cjk if is_mainly_cjk("".join(s.text for s in group)) else self.max_word_count_english
            if count_words("".join(s.text for s in group)) > max_wc:
                result_groups.extend(self._split_by_common_words(group))
            else:
                result_groups.append(group)

        # Split remaining long segments
        result = []
        for group in result_groups:
            result.extend(self._split_long_segment(group))
        return result

    def _group_by_time_gaps(
        self, segments: List[ASRDataSeg], max_gap: int = MAX_GAP, check_large_gaps: bool = False
    ) -> List[List[ASRDataSeg]]:
        """Group segments by time gaps."""
        if not segments:
            return []
        result = []
        current = [segments[0]]
        recent_gaps = []

        for i in range(1, len(segments)):
            time_gap = segments[i].start_time - segments[i - 1].end_time
            if check_large_gaps:
                recent_gaps.append(time_gap)
                if len(recent_gaps) > TIME_GAP_WINDOW_SIZE:
                    recent_gaps.pop(0)
                if len(recent_gaps) == TIME_GAP_WINDOW_SIZE:
                    avg = sum(recent_gaps) / len(recent_gaps)
                    if time_gap > avg * TIME_GAP_MULTIPLIER and len(current) > MIN_GROUP_SIZE:
                        result.append(current)
                        current = []
                        recent_gaps = []
            if time_gap > max_gap:
                result.append(current)
                current = []
                recent_gaps = []
            current.append(segments[i])

        if current:
            result.append(current)
        return result

    def _split_by_common_words(self, segments: List[ASRDataSeg]) -> List[List[ASRDataSeg]]:
        """Split at common connecting words."""
        prefix_words = {"and", "or", "but", "if", "then", "because", "while", "when", "however",
                        "和", "及", "但", "而", "或", "因", "我", "你", "他", "她", "这", "那"}
        suffix_words = {".", ",", "!", "?", "。", "，", "！", "？", "的", "了", "着", "过", "吗", "呢", "吧"}

        result = []
        current = []
        for i, seg in enumerate(segments):
            max_wc = self.max_word_count_cjk if is_mainly_cjk(seg.text) else self.max_word_count_english
            if any(seg.text.lower().startswith(w) for w in prefix_words) and len(current) >= int(max_wc * PREFIX_WORD_RATIO):
                result.append(current)
                current = []
            if i > 0 and any(segments[i - 1].text.lower().endswith(w) for w in suffix_words) and len(current) >= int(max_wc * SUFFIX_WORD_RATIO):
                result.append(current)
                current = []
            current.append(seg)
        if current:
            result.append(current)
        return result

    def _split_long_segment(self, segments: List[ASRDataSeg]) -> List[ASRDataSeg]:
        """Split segments that are too long by finding max time gap."""
        result_segs = []
        to_process = [segments]

        while to_process:
            current = to_process.pop(0)
            if not current:
                continue

            merged_text = "".join(s.text for s in current)
            max_wc = self.max_word_count_cjk if is_mainly_cjk(merged_text) else self.max_word_count_english
            n = len(current)

            if count_words(merged_text) <= max_wc or n < RULE_MIN_SEGMENT_SIZE:
                result_segs.append(ASRDataSeg(merged_text.strip(), current[0].start_time, current[-1].end_time))
                continue

            gaps = [current[i + 1].start_time - current[i].end_time for i in range(n - 1)]
            if all(abs(g - gaps[0]) < 1e-6 for g in gaps):
                split_idx = n // 2
            else:
                start_idx = max(n // 6, 1)
                end_idx = min((5 * n) // 6, n - 2)
                split_idx = max(range(start_idx, end_idx), key=lambda i: gaps[i], default=n // 2)

            to_process.extend([current[:split_idx + 1], current[split_idx + 1:]])

        result_segs.sort(key=lambda s: s.start_time)
        return result_segs

    def _merge_segments_based_on_sentences(
        self, segments: List[ASRDataSeg], sentences: List[str]
    ) -> List[ASRDataSeg]:
        """Match LLM sentences back to word-level ASR segments using sliding window."""
        def preprocess(s):
            return " ".join(s.lower().split())

        asr_texts = [seg.text for seg in segments]
        asr_len = len(asr_texts)
        asr_index = 0
        max_shift = MATCH_MAX_SHIFT
        unmatched = 0
        new_segments = []

        for sentence in sentences:
            sentence_proc = preprocess(sentence)
            word_count = count_words(sentence_proc)
            best_ratio = 0.0
            best_pos = None
            best_ws = 0

            max_ws = min(word_count * 2, asr_len - asr_index)
            min_ws = max(1, word_count // 2)
            window_sizes = sorted(range(min_ws, max_ws + 1), key=lambda x: abs(x - word_count))

            for ws in window_sizes:
                max_start = min(asr_index + max_shift + 1, asr_len - ws + 1)
                for start in range(asr_index, max_start):
                    substr = "".join(asr_texts[start:start + ws])
                    ratio = difflib.SequenceMatcher(None, sentence_proc, preprocess(substr)).ratio()
                    if ratio > best_ratio:
                        best_ratio = ratio
                        best_pos = start
                        best_ws = ws
                    if ratio == 1.0:
                        break
                if best_ratio == 1.0:
                    break

            if best_ratio >= MATCH_SIMILARITY_THRESHOLD and best_pos is not None:
                segs = segments[best_pos:best_pos + best_ws]
                seg_groups = self._group_by_time_gaps(segs, max_gap=MAX_GAP)
                for group in seg_groups:
                    new_segments.extend(self._split_long_segment(group))
                max_shift = MATCH_MAX_SHIFT
                asr_index = best_pos + best_ws
            else:
                logger.warning(f"Could not match sentence: {sentence[:50]}...")
                unmatched += 1
                if unmatched > MATCH_MAX_UNMATCHED:
                    raise ValueError(f"Too many unmatched sentences ({unmatched})")
                max_shift = MATCH_LARGE_SHIFT
                asr_index = min(asr_index + 1, asr_len - 1)

        return new_segments
