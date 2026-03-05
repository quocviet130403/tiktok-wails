"""ASR Data structures for subtitle processing.

Ported from VideoCaptioner's asr_data.py with simplifications:
- Removed PyQt dependencies
- Removed YouTube VTT, LRC formats
- Kept: SRT, ASS parsing and generation
- Kept: Word-level timestamp detection and splitting
"""

import json
import math
import os
import platform
import re
from enum import Enum
from pathlib import Path
from typing import List, Optional, Tuple

from core.text_utils import is_mainly_cjk

# Multi-language word split pattern
_WORD_SPLIT_PATTERN = (
    r"[a-zA-Z\u00c0-\u00ff\u0100-\u017f']+"  # Latin (extended)
    r"|[\u0400-\u04ff]+"  # Cyrillic
    r"|[\u0370-\u03ff]+"  # Greek
    r"|[\u0600-\u06ff]+"  # Arabic
    r"|[\u0590-\u05ff]+"  # Hebrew
    r"|\d+"  # Numbers
    r"|[\u4e00-\u9fff]"  # Chinese
    r"|[\u3040-\u309f]"  # Hiragana
    r"|[\u30a0-\u30ff]"  # Katakana
    r"|[\uac00-\ud7af]"  # Korean
    r"|[\u0e00-\u0e7f][\u0e30-\u0e3a\u0e47-\u0e4e]*"  # Thai
    r"|[\u0900-\u097f]"  # Devanagari
)


class SubtitleLayoutEnum(Enum):
    """Subtitle layout options."""
    TRANSLATE_ON_TOP = "译文在上"
    ORIGINAL_ON_TOP = "原文在上"
    ONLY_ORIGINAL = "仅原文"
    ONLY_TRANSLATE = "仅译文"


def _handle_long_path(path: str) -> str:
    """Handle Windows long path limitation."""
    if platform.system() == "Windows" and len(path) > 260 and not path.startswith(r"\\?\ "):
        return rf"\\?\{os.path.abspath(path)}"
    return path


class ASRDataSeg:
    """Single subtitle segment with timing information."""

    def __init__(self, text: str, start_time: int, end_time: int, translated_text: str = ""):
        self.text = text
        self.translated_text = translated_text
        self.start_time = start_time  # milliseconds
        self.end_time = end_time  # milliseconds

    def to_srt_ts(self) -> str:
        """Convert to SRT timestamp format."""
        return f"{self._ms_to_srt_time(self.start_time)} --> {self._ms_to_srt_time(self.end_time)}"

    def to_ass_ts(self) -> Tuple[str, str]:
        """Convert to ASS timestamp format."""
        return self._ms_to_ass_ts(self.start_time), self._ms_to_ass_ts(self.end_time)

    @staticmethod
    def _ms_to_srt_time(ms: int) -> str:
        """Convert milliseconds to SRT time format (HH:MM:SS,mmm)."""
        total_seconds, milliseconds = divmod(ms, 1000)
        minutes, seconds = divmod(total_seconds, 60)
        hours, minutes = divmod(minutes, 60)
        return f"{int(hours):02}:{int(minutes):02}:{int(seconds):02},{int(milliseconds):03}"

    @staticmethod
    def _ms_to_ass_ts(ms: int) -> str:
        """Convert milliseconds to ASS timestamp format (H:MM:SS.cc)."""
        total_seconds, milliseconds = divmod(ms, 1000)
        minutes, seconds = divmod(total_seconds, 60)
        hours, minutes = divmod(minutes, 60)
        centiseconds = int(milliseconds / 10)
        return f"{int(hours):01}:{int(minutes):02}:{int(seconds):02}.{centiseconds:02}"

    @property
    def transcript(self) -> str:
        return self.text

    def __str__(self) -> str:
        return f"ASRDataSeg({self.text}, {self.start_time}, {self.end_time})"


class ASRData:
    """Collection of ASR subtitle segments."""

    def __init__(self, segments: List[ASRDataSeg]):
        filtered = [seg for seg in segments if seg.text and seg.text.strip()]
        filtered.sort(key=lambda x: x.start_time)
        self.segments = filtered

    def __iter__(self):
        return iter(self.segments)

    def __len__(self) -> int:
        return len(self.segments)

    def has_data(self) -> bool:
        return len(self.segments) > 0

    def _is_word_level_segment(self, segment: ASRDataSeg) -> bool:
        """Check if a single segment is word-level."""
        text = segment.text.strip()
        if is_mainly_cjk(text):
            return len(text) <= 2
        return len(text.split()) == 1

    def is_word_timestamp(self) -> bool:
        """Check if timestamps are word-level (not sentence-level).

        Returns True if 80%+ segments are single words/characters.
        """
        if not self.segments:
            return False
        word_level_count = sum(1 for seg in self.segments if self._is_word_level_segment(seg))
        return (word_level_count / len(self.segments)) >= 0.8

    def split_to_word_segments(self) -> "ASRData":
        """Split sentence-level subtitles to word-level with estimated timestamps."""
        CHARS_PER_PHONEME = 4
        new_segments = []

        for seg in self.segments:
            text = seg.text
            duration = seg.end_time - seg.start_time
            words_list = list(re.finditer(_WORD_SPLIT_PATTERN, text))

            if not words_list:
                continue

            total_phonemes = sum(math.ceil(len(w.group()) / CHARS_PER_PHONEME) for w in words_list)
            time_per_phoneme = duration / max(total_phonemes, 1)

            current_time = seg.start_time
            for word_match in words_list:
                word = word_match.group()
                word_phonemes = math.ceil(len(word) / CHARS_PER_PHONEME)
                word_duration = int(time_per_phoneme * word_phonemes)
                word_end_time = min(current_time + word_duration, seg.end_time)
                new_segments.append(ASRDataSeg(text=word, start_time=current_time, end_time=word_end_time))
                current_time = word_end_time

        self.segments = new_segments
        return self

    def remove_punctuation(self) -> "ASRData":
        """Remove trailing Chinese punctuation from segments."""
        punct = r"[，。]"
        for seg in self.segments:
            seg.text = re.sub(f"{punct}+$", "", seg.text.strip())
            seg.translated_text = re.sub(f"{punct}+$", "", seg.translated_text.strip())
        return self

    def to_txt(self, save_path=None, layout: SubtitleLayoutEnum = SubtitleLayoutEnum.ORIGINAL_ON_TOP) -> str:
        """Convert to plain text."""
        result = []
        for seg in self.segments:
            original = seg.text
            translated = seg.translated_text
            if layout == SubtitleLayoutEnum.ORIGINAL_ON_TOP:
                text = f"{original}\n{translated}" if translated else original
            elif layout == SubtitleLayoutEnum.TRANSLATE_ON_TOP:
                text = f"{translated}\n{original}" if translated else original
            elif layout == SubtitleLayoutEnum.ONLY_ORIGINAL:
                text = original
            else:
                text = translated if translated else original
            result.append(text)
        text = "\n".join(result)
        if save_path:
            Path(save_path).parent.mkdir(parents=True, exist_ok=True)
            with open(_handle_long_path(save_path), "w", encoding="utf-8") as f:
                f.write(text)
        return text

    def to_srt(self, layout: SubtitleLayoutEnum = SubtitleLayoutEnum.ORIGINAL_ON_TOP, save_path=None) -> str:
        """Convert to SRT subtitle format."""
        srt_lines = []
        for n, seg in enumerate(self.segments, 1):
            original = seg.text
            translated = seg.translated_text
            if layout == SubtitleLayoutEnum.ORIGINAL_ON_TOP:
                text = f"{original}\n{translated}" if translated else original
            elif layout == SubtitleLayoutEnum.TRANSLATE_ON_TOP:
                text = f"{translated}\n{original}" if translated else original
            elif layout == SubtitleLayoutEnum.ONLY_ORIGINAL:
                text = original
            else:
                text = translated if translated else original
            srt_lines.append(f"{n}\n{seg.to_srt_ts()}\n{text}\n")

        srt_text = "\n".join(srt_lines)
        if save_path:
            Path(save_path).parent.mkdir(parents=True, exist_ok=True)
            with open(_handle_long_path(save_path), "w", encoding="utf-8") as f:
                f.write(srt_text)
        return srt_text

    def to_ass(
        self,
        style_str: Optional[str] = None,
        layout: SubtitleLayoutEnum = SubtitleLayoutEnum.ONLY_TRANSLATE,
        save_path: Optional[str] = None,
        video_width: int = 1280,
        video_height: int = 720,
    ) -> str:
        """Convert to ASS subtitle format."""
        if not style_str:
            style_str = (
                "[V4+ Styles]\n"
                "Format: Name,Fontname,Fontsize,PrimaryColour,SecondaryColour,OutlineColour,BackColour,"
                "Bold,Italic,Underline,StrikeOut,ScaleX,ScaleY,Spacing,Angle,BorderStyle,Outline,Shadow,"
                "Alignment,MarginL,MarginR,MarginV,Encoding\n"
                "Style: Default,Arial,28,&H00FFFFFF,&H000000FF,&H00000000,&H80000000,-1,0,0,0,100,100,"
                "0,0,3,1,0,2,10,10,40,1\n"
                "Style: Secondary,Arial,22,&H00FFFFFF,&H000000FF,&H00000000,&H80000000,-1,0,0,0,100,100,"
                "0,0,3,1,0,2,10,10,40,1"
            )

        ass_content = (
            "[Script Info]\n"
            "; Generated by tiktok-wails subtitle pipeline\n"
            "ScriptType: v4.00+\n"
            f"PlayResX: {video_width}\n"
            f"PlayResY: {video_height}\n\n"
            f"{style_str}\n\n"
            "[Events]\n"
            "Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text\n"
        )

        tmpl = "Dialogue: 0,{},{},{},,0,0,0,,{}\n"
        for seg in self.segments:
            start_time, end_time = seg.to_ass_ts()
            original = seg.text
            translated = seg.translated_text
            has_translation = bool(translated and translated.strip())

            if layout == SubtitleLayoutEnum.TRANSLATE_ON_TOP:
                if has_translation:
                    ass_content += tmpl.format(start_time, end_time, "Default", translated)
                    ass_content += tmpl.format(start_time, end_time, "Secondary", original)
                else:
                    ass_content += tmpl.format(start_time, end_time, "Default", original)
            elif layout == SubtitleLayoutEnum.ORIGINAL_ON_TOP:
                if has_translation:
                    ass_content += tmpl.format(start_time, end_time, "Default", original)
                    ass_content += tmpl.format(start_time, end_time, "Secondary", translated)
                else:
                    ass_content += tmpl.format(start_time, end_time, "Default", original)
            elif layout == SubtitleLayoutEnum.ONLY_ORIGINAL:
                ass_content += tmpl.format(start_time, end_time, "Default", original)
            else:  # ONLY_TRANSLATE
                text = translated if has_translation else original
                ass_content += tmpl.format(start_time, end_time, "Default", text)

        if save_path:
            Path(save_path).parent.mkdir(parents=True, exist_ok=True)
            with open(_handle_long_path(save_path), "w", encoding="utf-8") as f:
                f.write(ass_content)
        return ass_content

    def to_json(self) -> dict:
        """Convert to JSON format."""
        result = {}
        for i, seg in enumerate(self.segments, 1):
            result[str(i)] = {
                "start_time": seg.start_time,
                "end_time": seg.end_time,
                "original_subtitle": seg.text,
                "translated_subtitle": seg.translated_text,
            }
        return result

    @staticmethod
    def from_subtitle_file(file_path: str) -> "ASRData":
        """Load ASRData from subtitle file (.srt, .ass)."""
        path = Path(file_path)
        if not path.exists():
            raise FileNotFoundError(f"File not found: {path}")

        try:
            content = path.read_text(encoding="utf-8")
        except UnicodeDecodeError:
            content = path.read_text(encoding="gbk")

        suffix = path.suffix.lower()
        if suffix == ".srt":
            return ASRData.from_srt(content)
        elif suffix == ".ass":
            return ASRData.from_ass(content)
        elif suffix == ".json":
            return ASRData.from_json(json.loads(content))
        else:
            raise ValueError(f"Unsupported file format: {suffix}")

    @staticmethod
    def from_json(json_data: dict) -> "ASRData":
        """Create ASRData from JSON data."""
        segments = []
        for i in sorted(json_data.keys(), key=int):
            d = json_data[i]
            segments.append(ASRDataSeg(
                text=d["original_subtitle"],
                translated_text=d.get("translated_subtitle", ""),
                start_time=d["start_time"],
                end_time=d["end_time"],
            ))
        return ASRData(segments)

    @staticmethod
    def from_srt(srt_str: str) -> "ASRData":
        """Create ASRData from SRT format string."""
        segments = []
        srt_time_pattern = re.compile(
            r"(\d{2}):(\d{2}):(\d{1,2})[.,](\d{3})\s-->\s(\d{2}):(\d{2}):(\d{1,2})[.,](\d{3})"
        )
        blocks = re.split(r"\n\s*\n", srt_str.strip())

        for block in blocks:
            lines = block.splitlines()
            if len(lines) < 3:
                continue

            match = srt_time_pattern.match(lines[1])
            if not match:
                continue

            tp = list(map(int, match.groups()))
            start_time = tp[0] * 3600000 + tp[1] * 60000 + tp[2] * 1000 + tp[3]
            end_time = tp[4] * 3600000 + tp[5] * 60000 + tp[6] * 1000 + tp[7]
            segments.append(ASRDataSeg(" ".join(lines[2:]), start_time, end_time))

        return ASRData(segments)

    @staticmethod
    def from_ass(ass_str: str) -> "ASRData":
        """Create ASRData from ASS format string."""
        segments = []
        ass_time_pattern = re.compile(
            r"Dialogue: \d+,(\d+:\d{2}:\d{2}\.\d{2}),(\d+:\d{2}:\d{2}\.\d{2}),(.*?),.*?,\d+,\d+,\d+,.*?,(.*?)$"
        )

        def parse_ass_time(time_str: str) -> int:
            hours, minutes, seconds = time_str.split(":")
            seconds, centiseconds = seconds.split(".")
            return int(hours) * 3600000 + int(minutes) * 60000 + int(seconds) * 1000 + int(centiseconds) * 10

        for line in ass_str.splitlines():
            if line.startswith("Dialogue:"):
                match = ass_time_pattern.match(line)
                if match:
                    start_time = parse_ass_time(match.group(1))
                    end_time = parse_ass_time(match.group(2))
                    text = match.group(4)
                    text = re.sub(r"\{[^}]*\}", "", text)
                    text = text.replace("\\N", "\n").strip()
                    if text:
                        segments.append(ASRDataSeg(text, start_time, end_time))

        return ASRData(segments)

    def __str__(self):
        return self.to_txt()
