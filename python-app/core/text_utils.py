"""Text utility functions for subtitle processing."""

import re
import unicodedata

# CJK Unicode ranges
CJK_PATTERN = re.compile(
    r"[\u4e00-\u9fff]"        # CJK Unified Ideographs
    r"|[\u3400-\u4dbf]"       # CJK Extension A
    r"|[\u3040-\u309f]"       # Hiragana
    r"|[\u30a0-\u30ff]"       # Katakana
    r"|[\uac00-\ud7af]"       # Korean Syllables
    r"|[\u0e00-\u0e7f]"       # Thai
)

# Punctuation pattern
PUNCTUATION_PATTERN = re.compile(
    r'^[\s\u3000-\u303f\uff00-\uffef\u2000-\u206f'
    r'!"#$%&\'()*+,\-./:;<=>?@\[\\\]^_`{|}~]+$'
)

# Trailing CJK punctuation
TRAILING_PUNCT = re.compile(r'[，。、！？；：,\.!?;:]+$')


def count_words(text: str) -> int:
    """Count words/characters based on language.

    For CJK text: counts characters.
    For space-separated languages: counts words.
    """
    if not text or not text.strip():
        return 0

    if is_mainly_cjk(text):
        # Count CJK characters
        return sum(1 for char in text if CJK_PATTERN.match(char))
    else:
        # Count space-separated words
        return len(text.split())


def is_mainly_cjk(text: str) -> bool:
    """Check if text is mainly CJK (>30% CJK characters)."""
    if not text:
        return False

    total_chars = len(text.replace(" ", ""))
    if total_chars == 0:
        return False

    cjk_count = sum(1 for char in text if CJK_PATTERN.match(char))
    return cjk_count / total_chars > 0.3


def is_pure_punctuation(text: str) -> bool:
    """Check if text contains only punctuation characters."""
    if not text or not text.strip():
        return True
    return bool(PUNCTUATION_PATTERN.match(text.strip()))


def is_space_separated_language(text: str) -> bool:
    """Check if text uses space-separated words (Latin, Cyrillic, etc.)."""
    if not text:
        return False
    # If mainly CJK, not space-separated
    if is_mainly_cjk(text):
        return False
    # Check for Latin or Cyrillic characters
    for char in text:
        cat = unicodedata.category(char)
        if cat.startswith('L'):  # Letter category
            script = unicodedata.name(char, '').split()[0] if unicodedata.name(char, '') else ''
            if script in ('LATIN', 'CYRILLIC', 'GREEK', 'ARMENIAN', 'GEORGIAN'):
                return True
    return False


def remove_trailing_punctuation(text: str) -> str:
    """Remove trailing Chinese/English punctuation from text."""
    return TRAILING_PUNCT.sub('', text).strip()
