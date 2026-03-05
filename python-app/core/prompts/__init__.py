"""Prompt management module.

All prompts are stored as Markdown files, supporting template variable substitution.

Usage:
    from core.prompts import get_prompt

    prompt = get_prompt("split/sentence")
    prompt = get_prompt("translate/standard", target_language="越南语")
"""

import functools
from pathlib import Path
from string import Template

PROMPTS_DIR = Path(__file__).parent


@functools.lru_cache(maxsize=32)
def _load_prompt_file(prompt_path: str) -> str:
    """Load prompt from file (with LRU cache)."""
    file_path = PROMPTS_DIR / f"{prompt_path}.md"
    if not file_path.exists():
        raise FileNotFoundError(f"Prompt file not found: {prompt_path}.md at {file_path}")
    return file_path.read_text(encoding="utf-8")


def get_prompt(prompt_path: str, **kwargs) -> str:
    """Get prompt with template variable substitution.

    Args:
        prompt_path: Prompt path, e.g. "split/sentence", "optimize/subtitle"
        **kwargs: Template variables to replace $variable or ${variable}

    Returns:
        Processed prompt text
    """
    raw_prompt = _load_prompt_file(prompt_path)
    if not kwargs:
        return raw_prompt
    template = Template(raw_prompt)
    return template.safe_substitute(**kwargs)


def reload_cache():
    """Clear prompt cache (for development hot-reload)."""
    _load_prompt_file.cache_clear()


__all__ = ["get_prompt", "reload_cache"]
