"""Unified LLM client with retry and thread-safe singleton."""

import logging
import os
import threading
from typing import Any, List, Optional
from urllib.parse import urlparse, urlunparse

import openai
from openai import OpenAI
from tenacity import (
    RetryCallState,
    retry,
    retry_if_exception_type,
    stop_after_attempt,
    wait_random_exponential,
)

logger = logging.getLogger("llm_client")

_global_client: Optional[OpenAI] = None
_client_lock = threading.Lock()


def normalize_base_url(base_url: str) -> str:
    """Normalize API base URL by ensuring /v1 suffix when needed."""
    url = base_url.strip()
    parsed = urlparse(url)
    path = parsed.path.rstrip("/")

    if not path:
        path = "/v1"

    return urlunparse(
        (parsed.scheme, parsed.netloc, path, parsed.params, parsed.query, parsed.fragment)
    )


def reset_client():
    """Reset the global client (used when config changes)."""
    global _global_client
    with _client_lock:
        _global_client = None


def init_client(base_url: str, api_key: str):
    """Initialize the global LLM client with given credentials."""
    global _global_client
    with _client_lock:
        normalized_url = normalize_base_url(base_url)
        _global_client = OpenAI(base_url=normalized_url, api_key=api_key)
        logger.info(f"LLM client initialized with base_url={normalized_url}")


def get_llm_client() -> OpenAI:
    """Get global LLM client instance (thread-safe singleton)."""
    global _global_client

    if _global_client is None:
        with _client_lock:
            if _global_client is None:
                base_url = os.getenv("OPENAI_BASE_URL", "").strip()
                api_key = os.getenv("OPENAI_API_KEY", "").strip()

                if not base_url or not api_key:
                    raise ValueError(
                        "LLM not configured. Set OPENAI_BASE_URL and OPENAI_API_KEY "
                        "or call init_client() first."
                    )

                base_url = normalize_base_url(base_url)
                _global_client = OpenAI(base_url=base_url, api_key=api_key)

    return _global_client


def _before_sleep_log(retry_state: RetryCallState) -> None:
    logger.warning("Rate limit hit, retrying with exponential backoff...")


@retry(
    stop=stop_after_attempt(10),
    wait=wait_random_exponential(multiplier=1, min=5, max=60),
    retry=retry_if_exception_type(openai.RateLimitError),
    before_sleep=_before_sleep_log,
)
def call_llm(
    messages: List[dict],
    model: str,
    temperature: float = 1.0,
    **kwargs: Any,
) -> Any:
    """Call LLM API with automatic retry on rate limits.

    Args:
        messages: Chat messages list
        model: Model name (e.g. 'deepseek-chat', 'gpt-4o-mini')
        temperature: Sampling temperature
        **kwargs: Additional OpenAI API params

    Returns:
        OpenAI ChatCompletion response

    Raises:
        ValueError: If response is empty or invalid
    """
    client = get_llm_client()

    response = client.chat.completions.create(
        model=model,
        messages=messages,
        temperature=temperature,
        **kwargs,
    )

    if not (
        response
        and hasattr(response, "choices")
        and response.choices
        and len(response.choices) > 0
        and hasattr(response.choices[0], "message")
        and response.choices[0].message.content
    ):
        raise ValueError("Invalid LLM API response: empty choices or content")

    return response
