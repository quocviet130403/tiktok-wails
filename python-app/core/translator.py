"""Subtitle translator - LLM-based translation with multi-threading and cache.

Ported from VideoCaptioner's translate module with full feature parity:
- Multi-threaded batch translation (ThreadPoolExecutor)
- Translation cache (disk-based, 7 day expiry)
- Standard and reflect (3-stage) translation modes
- Single-line fallback on batch failure
"""

import hashlib
import json
import logging
import os
import pickle
import time
from concurrent.futures import ThreadPoolExecutor, as_completed
from pathlib import Path
from typing import Any, Callable, Dict, List, Optional, Tuple, Union

import json_repair
import openai

from core.asr_data import ASRData, ASRDataSeg
from core.llm_client import call_llm
from core.prompts import get_prompt

logger = logging.getLogger("subtitle_translator")

MAX_STEPS = 3
DEFAULT_BATCH_SIZE = 10
DEFAULT_THREAD_NUM = 4
CACHE_EXPIRY_SECONDS = 7 * 24 * 3600  # 7 days


class TranslationCache:
    """Simple disk-based translation cache."""

    def __init__(self, cache_dir: Optional[str] = None):
        if cache_dir is None:
            cache_dir = os.path.join(os.path.dirname(__file__), "..", ".cache", "translate")
        self.cache_dir = Path(cache_dir)
        self.cache_dir.mkdir(parents=True, exist_ok=True)

    def _key_to_path(self, key: str) -> Path:
        hashed = hashlib.md5(key.encode()).hexdigest()
        return self.cache_dir / f"{hashed}.pkl"

    def get(self, key: str) -> Optional[Any]:
        path = self._key_to_path(key)
        if not path.exists():
            return None
        try:
            with open(path, "rb") as f:
                data = pickle.load(f)
            if time.time() - data.get("timestamp", 0) > CACHE_EXPIRY_SECONDS:
                path.unlink(missing_ok=True)
                return None
            return data.get("value")
        except Exception:
            return None

    def set(self, key: str, value: Any) -> None:
        path = self._key_to_path(key)
        try:
            with open(path, "wb") as f:
                pickle.dump({"value": value, "timestamp": time.time()}, f)
        except Exception:
            pass


class LLMTranslator:
    """LLM-based subtitle translator with multi-threading, cache, and agent loop."""

    def __init__(
        self,
        model: str,
        target_language: str = "越南语",
        batch_num: int = DEFAULT_BATCH_SIZE,
        thread_num: int = DEFAULT_THREAD_NUM,
        custom_prompt: str = "",
        is_reflect: bool = False,
        update_callback: Optional[Callable] = None,
    ):
        self.model = model
        self.target_language = target_language
        self.batch_num = batch_num
        self.thread_num = thread_num
        self.custom_prompt = custom_prompt
        self.is_reflect = is_reflect
        self.update_callback = update_callback
        self._cache = TranslationCache()

    def translate_subtitle(self, subtitle_data: Union[str, ASRData]) -> ASRData:
        """Translate all subtitles with multi-threading (main entry)."""
        try:
            if isinstance(subtitle_data, str):
                asr_data = ASRData.from_subtitle_file(subtitle_data)
            else:
                asr_data = subtitle_data

            # Build chunks
            translate_items = [
                {"index": i, "text": seg.text}
                for i, seg in enumerate(asr_data.segments, 1)
            ]

            chunks = [
                translate_items[i:i + self.batch_num]
                for i in range(0, len(translate_items), self.batch_num)
            ]

            logger.info(f"Translating {len(translate_items)} segments in {len(chunks)} batches "
                        f"with {self.thread_num} threads")

            # Multi-threaded translation
            translation_map: Dict[int, str] = {}

            with ThreadPoolExecutor(max_workers=self.thread_num) as executor:
                futures = {
                    executor.submit(self._safe_translate_chunk, chunk): chunk
                    for chunk in chunks
                }

                for future in as_completed(futures):
                    chunk = futures[future]
                    try:
                        result = future.result()
                        translation_map.update(result)
                    except Exception as e:
                        logger.error(f"Chunk translation failed: {e}")
                        # Fallback: add original text
                        for item in chunk:
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

    def _get_cache_key(self, chunk: List[dict]) -> str:
        """Generate cache key for a chunk."""
        text = "|".join(f"{item['index']}:{item['text']}" for item in chunk)
        return f"LLMTranslator:{self.model}:{self.target_language}:{hashlib.md5(text.encode()).hexdigest()}"

    def _safe_translate_chunk(self, chunk: List[dict]) -> Dict[int, str]:
        """Translate with cache check."""
        cache_key = self._get_cache_key(chunk)
        cached = self._cache.get(cache_key)
        if cached is not None:
            logger.info(f"Cache hit for segments {chunk[0]['index']}-{chunk[-1]['index']}")
            return cached

        result = self._translate_chunk(chunk)

        # Cache result
        self._cache.set(cache_key, result)
        return result

    def _translate_chunk(self, chunk: List[dict]) -> Dict[int, str]:
        """Translate a batch of subtitles using agent loop."""
        subtitle_dict = {str(item["index"]): item["text"] for item in chunk}

        logger.info(f"Translating subtitles: {chunk[0]['index']} - {chunk[-1]['index']}")

        # Get prompt
        prompt_name = "translate/reflect" if self.is_reflect else "translate/standard"
        prompt = get_prompt(
            prompt_name,
            target_language=self.target_language,
            custom_prompt=self.custom_prompt,
        )

        try:
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

        except openai.RateLimitError as e:
            logger.error(f"Rate limit: {e}")
            raise
        except openai.AuthenticationError as e:
            logger.error(f"Auth error: {e}")
            raise
        except Exception as e:
            logger.warning(f"Batch translate failed, trying single mode: {e}")
            return self._translate_chunk_single(chunk)

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
                sort_key = lambda x: int(x) if x.isdigit() else x
                parts.append(f"Missing keys {sorted(missing, key=sort_key)} - you must translate these items")
            if extra:
                parts.append(f"Extra keys {sorted(extra, key=sort_key)} - remove them")
            return False, "; ".join(parts)

        if self.is_reflect:
            for key, value in response_dict.items():
                if not isinstance(value, dict):
                    return False, f"Key '{key}': value must be a dict with 'native_translation' field"
                if "native_translation" not in value:
                    return False, f"Key '{key}': missing 'native_translation' field"

        return True, ""

    def _translate_chunk_single(self, chunk: List[dict]) -> Dict[int, str]:
        """Fallback: translate one line at a time."""
        single_prompt = get_prompt("translate/single", target_language=self.target_language)
        result = {}

        for item in chunk:
            try:
                response = call_llm(
                    messages=[
                        {"role": "system", "content": single_prompt},
                        {"role": "user", "content": item["text"]},
                    ],
                    model=self.model,
                    temperature=0.7,
                )
                result[item["index"]] = response.choices[0].message.content.strip()
            except Exception as e:
                logger.error(f"Single translate failed for {item['index']}: {e}")
                result[item["index"]] = item["text"]

        return result
