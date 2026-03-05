"""Subtitle optimizer - Fix ASR errors using LLM agent loop.

Ported from VideoCaptioner's optimize module.
Pipeline: LLM corrects text → Validate (keys match + similarity ≥70%) → Retry up to 3x
"""

import difflib
import logging
import re
from typing import Callable, Dict, List, Optional, Tuple, Union

import json_repair

from core.asr_data import ASRData, ASRDataSeg
from core.llm_client import call_llm
from core.prompts import get_prompt
from core.text_utils import count_words

logger = logging.getLogger("subtitle_optimizer")

MAX_STEPS = 3
BATCH_SIZE = 10


class SubtitleAligner:
    """Align two text sequences using difflib for subtitle repair."""

    def __init__(self):
        self.line_numbers = [0, 0]

    def align_texts(self, source_text: List[str], target_text: List[str]) -> Tuple[List[str], List[str]]:
        """Align source and target text lists."""
        diff_iterator = difflib.ndiff(source_text, target_text)
        return self._pair_lines(diff_iterator)

    def _pair_lines(self, diff_iterator):
        source_lines = []
        target_lines = []
        flag = 0

        for source_line, target_line, _ in self._line_iterator(diff_iterator):
            if source_line is not None:
                if source_line[1] == "\n":
                    flag += 1
                    continue
                source_lines.append(source_line[1])
            if target_line is not None:
                if flag > 0:
                    flag -= 1
                    continue
                target_lines.append(target_line[1])

        for i in range(1, len(target_lines)):
            if target_lines[i] == "\n":
                target_lines[i] = target_lines[i - 1]

        return source_lines, target_lines

    def _line_iterator(self, diff_iterator):
        lines = []
        blank_lines_pending = 0
        blank_lines_to_yield = 0

        while True:
            while len(lines) < 4:
                lines.append(next(diff_iterator, "X"))

            diff_type = "".join([line[0] for line in lines])

            if diff_type.startswith("X"):
                blank_lines_to_yield = blank_lines_pending
            elif diff_type.startswith("-?+?"):
                yield self._format_line(lines, "?", 0), self._format_line(lines, "?", 1), True
                continue
            elif diff_type.startswith("--++"):
                blank_lines_pending -= 1
                yield self._format_line(lines, "-", 0), None, True
                continue
            elif diff_type.startswith(("--?+", "--+", "- ")):
                source_line, target_line = self._format_line(lines, "-", 0), None
                blank_lines_to_yield, blank_lines_pending = blank_lines_pending - 1, 0
            elif diff_type.startswith("-+?"):
                yield self._format_line(lines, None, 0), self._format_line(lines, "?", 1), True
                continue
            elif diff_type.startswith("-?+"):
                yield self._format_line(lines, "?", 0), self._format_line(lines, None, 1), True
                continue
            elif diff_type.startswith("-"):
                blank_lines_pending -= 1
                yield self._format_line(lines, "-", 0), None, True
                continue
            elif diff_type.startswith("+--"):
                blank_lines_pending += 1
                yield None, self._format_line(lines, "+", 1), True
                continue
            elif diff_type.startswith(("+ ", "+-")):
                source_line, target_line = None, self._format_line(lines, "+", 1)
                blank_lines_to_yield, blank_lines_pending = blank_lines_pending + 1, 0
            elif diff_type.startswith("+"):
                blank_lines_pending += 1
                yield None, self._format_line(lines, "+", 1), True
                continue
            elif diff_type.startswith(" "):
                yield self._format_line(lines[:], None, 0), self._format_line(lines, None, 1), False
                continue

            while blank_lines_to_yield < 0:
                blank_lines_to_yield += 1
                yield None, ("", "\n"), True
            while blank_lines_to_yield > 0:
                blank_lines_to_yield -= 1
                yield ("", "\n"), None, True

            if diff_type.startswith("X"):
                return
            else:
                yield source_line, target_line, True

    def _format_line(self, lines, format_key, side):
        self.line_numbers[side] += 1
        if format_key is None:
            return self.line_numbers[side], lines.pop(0)[2:]
        if format_key == "?":
            text = lines.pop(0)
            lines.pop(0)
            text = text[2:]
        else:
            text = lines.pop(0)[2:]
            if not text:
                text = ""
        return self.line_numbers[side], text


class SubtitleOptimizer:
    """Optimize subtitles using LLM with agent loop validation."""

    def __init__(
        self,
        model: str,
        batch_num: int = BATCH_SIZE,
        custom_prompt: str = "",
        update_callback: Optional[Callable] = None,
    ):
        self.model = model
        self.batch_num = batch_num
        self.custom_prompt = custom_prompt
        self.update_callback = update_callback

    def optimize_subtitle(self, subtitle_data: Union[str, ASRData]) -> ASRData:
        """Optimize subtitle text (main entry)."""
        try:
            if isinstance(subtitle_data, str):
                asr_data = ASRData.from_subtitle_file(subtitle_data)
            else:
                asr_data = subtitle_data

            subtitle_dict = {str(i): seg.text for i, seg in enumerate(asr_data.segments, 1)}
            chunks = self._split_chunks(subtitle_dict)

            optimized_dict = {}
            for chunk in chunks:
                try:
                    result = self._optimize_chunk(chunk)
                    optimized_dict.update(result)
                except Exception as e:
                    logger.error(f"Optimize chunk failed: {e}")
                    optimized_dict.update(chunk)

            new_segments = [
                ASRDataSeg(
                    text=optimized_dict.get(str(i), seg.text),
                    start_time=seg.start_time,
                    end_time=seg.end_time,
                )
                for i, seg in enumerate(asr_data.segments, 1)
            ]
            return ASRData(new_segments)

        except Exception as e:
            logger.error(f"Optimization failed: {e}")
            raise RuntimeError(f"Optimization failed: {e}")

    def _split_chunks(self, subtitle_dict: Dict[str, str]) -> List[Dict[str, str]]:
        items = list(subtitle_dict.items())
        return [dict(items[i:i + self.batch_num]) for i in range(0, len(items), self.batch_num)]

    def _optimize_chunk(self, subtitle_chunk: Dict[str, str]) -> Dict[str, str]:
        start_idx = next(iter(subtitle_chunk))
        end_idx = next(reversed(subtitle_chunk))
        logger.info(f"Optimizing subtitles: {start_idx} - {end_idx}")

        result = self._agent_loop(subtitle_chunk)

        if self.update_callback:
            self.update_callback(list(result.keys()))

        return result

    def _agent_loop(self, subtitle_chunk: Dict[str, str]) -> Dict[str, str]:
        """LLM → Validate → Feedback → Retry (max 3 times)."""
        user_prompt = (
            f"Correct the following subtitles. Keep the original language, do not translate:\n"
            f"<input_subtitle>{str(subtitle_chunk)}</input_subtitle>"
        )
        if self.custom_prompt:
            user_prompt += f"\nReference content:\n<reference>{self.custom_prompt}</reference>"

        messages = [
            {"role": "system", "content": get_prompt("optimize/subtitle")},
            {"role": "user", "content": user_prompt},
        ]

        last_result = None
        for step in range(MAX_STEPS):
            response = call_llm(messages=messages, model=self.model, temperature=0.2)
            result_text = response.choices[0].message.content
            if not result_text:
                raise ValueError("LLM returned empty result")

            parsed = json_repair.loads(result_text)
            if not isinstance(parsed, dict):
                raise ValueError(f"Expected dict, got {type(parsed)}")

            result_dict: Dict[str, str] = parsed
            last_result = result_dict

            is_valid, error_msg = self._validate(subtitle_chunk, result_dict)
            if is_valid:
                return self._repair_subtitle(subtitle_chunk, result_dict)

            logger.warning(f"Optimization validation failed (attempt {step + 1}): {error_msg}")
            messages.append({"role": "assistant", "content": result_text})
            messages.append({
                "role": "user",
                "content": f"Validation failed: {error_msg}\nPlease fix and output ONLY a valid JSON dictionary.",
            })

        logger.warning(f"Max attempts ({MAX_STEPS}) reached, using last result")
        return self._repair_subtitle(subtitle_chunk, last_result) if last_result else subtitle_chunk

    def _validate(self, original: Dict[str, str], optimized: Dict[str, str]) -> Tuple[bool, str]:
        """Validate: keys match + similarity ≥ 70%."""
        expected = set(original.keys())
        actual = set(optimized.keys())

        if expected != actual:
            missing = expected - actual
            extra = actual - expected
            parts = []
            if missing:
                parts.append(f"Missing keys: {sorted(missing)}")
            if extra:
                parts.append(f"Extra keys: {sorted(extra)}")
            return False, "\n".join(parts) + f"\nRequired keys: {sorted(expected)}"

        # Check similarity
        excessive = []
        for key in expected:
            orig = re.sub(r"\s+", " ", original[key]).strip()
            opt = re.sub(r"\s+", " ", optimized[key]).strip()
            similarity = difflib.SequenceMatcher(None, orig, opt).ratio()
            threshold = 0.3 if count_words(original[key]) <= 10 else 0.7
            if similarity < threshold:
                excessive.append(f"Key '{key}': similarity {similarity:.1%} < {threshold:.0%}")

        if excessive:
            return False, ";\n".join(excessive) + "\nMake MINIMAL changes only."

        return True, ""

    @staticmethod
    def _repair_subtitle(original: Dict[str, str], optimized: Dict[str, str]) -> Dict[str, str]:
        """Repair alignment between original and optimized."""
        try:
            aligner = SubtitleAligner()
            aligned_src, aligned_tgt = aligner.align_texts(list(original.values()), list(optimized.values()))
            if len(aligned_src) != len(aligned_tgt):
                return optimized
            start_id = next(iter(original.keys()))
            return {str(int(start_id) + i): text for i, text in enumerate(aligned_tgt)}
        except Exception as e:
            logger.error(f"Alignment failed: {e}")
            return optimized
