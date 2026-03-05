"""Subtitle translator - LLM-based translation with agent loop.

Ported from VideoCaptioner's translate module.
Supports standard and reflect (3-stage) translation modes.
"""

import json
import logging
from typing import Any, Callable, Dict, List, Optional, Tuple, Union

import json_repair
import openai

from core.asr_data import ASRData, ASRDataSeg
from core.llm_client import call_llm
from core.prompts import get_prompt

logger = logging.getLogger("subtitle_translator")

MAX_STEPS = 3
DEFAULT_BATCH_SIZE = 10


class LLMTranslator:
    """LLM-based subtitle translator with agent loop validation."""

    def __init__(
        self,
        model: str,
        target_language: str = "越南语",
        batch_num: int = DEFAULT_BATCH_SIZE,
        custom_prompt: str = "",
        is_reflect: bool = False,
        update_callback: Optional[Callable] = None,
    ):
        self.model = model
        self.target_language = target_language
        self.batch_num = batch_num
        self.custom_prompt = custom_prompt
        self.is_reflect = is_reflect
        self.update_callback = update_callback

    def translate_subtitle(self, subtitle_data: Union[str, ASRData]) -> ASRData:
        """Translate all subtitles (main entry)."""
        try:
            if isinstance(subtitle_data, str):
                asr_data = ASRData.from_subtitle_file(subtitle_data)
            else:
                asr_data = subtitle_data

            # Build translation data
            translate_items = [
                {"index": i, "text": seg.text}
                for i, seg in enumerate(asr_data.segments, 1)
            ]

            # Split into chunks
            chunks = [
                translate_items[i:i + self.batch_num]
                for i in range(0, len(translate_items), self.batch_num)
            ]

            # Translate all chunks
            translation_map: Dict[int, str] = {}
            for chunk in chunks:
                try:
                    result = self._translate_chunk(chunk)
                    translation_map.update(result)
                except Exception as e:
                    logger.error(f"Chunk translation failed: {e}")
                    # Try single-line fallback
                    for item in chunk:
                        try:
                            text = self._translate_single(item["text"])
                            translation_map[item["index"]] = text
                        except Exception as e2:
                            logger.error(f"Single translation failed for {item['index']}: {e2}")
                            translation_map[item["index"]] = item["text"]

            # Apply translations
            for i, seg in enumerate(asr_data.segments, 1):
                seg.translated_text = translation_map.get(i, seg.text)

            if self.update_callback:
                self.update_callback(len(translation_map))

            return asr_data

        except Exception as e:
            logger.error(f"Translation failed: {e}")
            raise RuntimeError(f"Translation failed: {e}")

    def _translate_chunk(self, chunk: List[dict]) -> Dict[int, str]:
        """Translate a batch of subtitles using agent loop."""
        subtitle_dict = {str(item["index"]): item["text"] for item in chunk}

        logger.info(f"Translating subtitles: {chunk[0]['index']} - {chunk[-1]['index']}")

        # Get prompt
        if self.is_reflect:
            prompt = get_prompt(
                "translate/reflect",
                target_language=self.target_language,
                custom_prompt=self.custom_prompt,
            )
        else:
            prompt = get_prompt(
                "translate/standard",
                target_language=self.target_language,
                custom_prompt=self.custom_prompt,
            )

        # Agent loop
        result_dict = self._agent_loop(prompt, subtitle_dict)

        # Process reflect mode results
        if self.is_reflect and isinstance(result_dict, dict):
            processed = {}
            for k, v in result_dict.items():
                if isinstance(v, dict):
                    processed[int(k)] = str(v.get("native_translation", v))
                else:
                    processed[int(k)] = str(v)
        else:
            processed = {int(k): str(v) for k, v in result_dict.items()}

        return processed

    def _agent_loop(self, system_prompt: str, subtitle_dict: Dict[str, str]) -> Dict[str, str]:
        """LLM → Validate → Feedback → Retry."""
        messages = [
            {"role": "system", "content": system_prompt},
            {"role": "user", "content": json.dumps(subtitle_dict, ensure_ascii=False)},
        ]

        last_result = None
        for _ in range(MAX_STEPS):
            response = call_llm(messages=messages, model=self.model)
            content = response.choices[0].message.content.strip()
            response_dict = json_repair.loads(content)
            last_result = response_dict

            is_valid, error_msg = self._validate(response_dict, subtitle_dict)
            if is_valid:
                return response_dict

            messages.append({"role": "assistant", "content": json.dumps(response_dict, ensure_ascii=False)})
            messages.append({
                "role": "user",
                "content": f"Error: {error_msg}\n\nFix and output ONLY a valid JSON dictionary with ALL {len(subtitle_dict)} keys",
            })

        return last_result if last_result else subtitle_dict

    def _validate(self, response_dict: Any, subtitle_dict: Dict[str, str]) -> Tuple[bool, str]:
        """Validate translation result."""
        if not isinstance(response_dict, dict):
            return False, f"Output must be a dict, got {type(response_dict).__name__}"

        expected = set(subtitle_dict.keys())
        actual = set(response_dict.keys())

        if expected != actual:
            missing = expected - actual
            extra = actual - expected
            parts = []
            if missing:
                parts.append(f"Missing keys: {sorted(missing, key=lambda x: int(x) if x.isdigit() else x)}")
            if extra:
                parts.append(f"Extra keys: {sorted(extra, key=lambda x: int(x) if x.isdigit() else x)}")
            return False, "; ".join(parts)

        # Validate reflect mode structure
        if self.is_reflect:
            for key, value in response_dict.items():
                if not isinstance(value, dict):
                    return False, f"Key '{key}': value must be a dict with 'native_translation' field"
                if "native_translation" not in value:
                    return False, f"Key '{key}': missing 'native_translation' field"

        return True, ""

    def _translate_single(self, text: str) -> str:
        """Fallback: translate a single line."""
        single_prompt = get_prompt("translate/single", target_language=self.target_language)
        response = call_llm(
            messages=[
                {"role": "system", "content": single_prompt},
                {"role": "user", "content": text},
            ],
            model=self.model,
            temperature=0.7,
        )
        return response.choices[0].message.content.strip()
