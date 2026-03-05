"""CLI entry point for subtitle pipeline.

Called by Go backend via exec.Command:
    python translate_cli.py --video_path "path/to/video.mp4" --config '{"llm_model": "deepseek-chat"}'

Outputs JSON to stdout for Go to parse.
Logs go to stderr so they don't interfere with JSON output.
"""

import argparse
import json
import logging
import os
import sys
import time

# Setup logging to stderr (stdout is reserved for JSON output to Go)
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(name)s] %(levelname)s: %(message)s",
    datefmt="%H:%M:%S",
    stream=sys.stderr,
)
logger = logging.getLogger("translate_cli")


def run_pipeline(
    video_path: str,
    base_url: str = "https://api.deepseek.com/v1",
    api_key: str = "sk-f6e9fbcd490845cb863e7bd660677c86",
    llm_model: str = "deepseek-chat",
    whisper_model: str = "large-v3-turbo",
    source_language: str = "zh",
    target_language: str = "越南语",
    use_reflect: bool = False,
    skip_optimize: bool = False,
    skip_split: bool = False,
) -> dict:
    """Run the full subtitle pipeline.

    Returns dict with status, output_path, etc.
    """
    # Init LLM client
    from core.llm_client import init_client
    init_client(base_url, api_key)

    # Output paths
    base_name = os.path.splitext(video_path)[0]
    ass_path = f"{base_name}-sub.ass"
    output_path = f"{base_name}-sub.mp4"

    start_time = time.time()

    # ========== Step 1: ASR ==========
    logger.info(f"🎧 Step 1/5: Transcribing ({whisper_model})...")
    from core.asr_engine import transcribe_video
    asr_data = transcribe_video(
        video_path=video_path,
        model_size=whisper_model,
        language=source_language,
    )
    logger.info(f"   → {len(asr_data)} word segments")

    # ========== Step 2: Split ==========
    if not skip_split:
        logger.info("✂️  Step 2/5: Splitting sentences...")
        from core.split import SubtitleSplitter
        splitter = SubtitleSplitter(model=llm_model)
        asr_data = splitter.split_subtitle(asr_data)
        logger.info(f"   → {len(asr_data)} sentences")
    else:
        logger.info("⏭️  Step 2/5: Skipped")

    # ========== Step 3: Optimize ==========
    if not skip_optimize:
        logger.info("🔧 Step 3/5: Fixing ASR errors...")
        from core.optimizer import SubtitleOptimizer
        optimizer = SubtitleOptimizer(model=llm_model)
        asr_data = optimizer.optimize_subtitle(asr_data)
        logger.info(f"   → {len(asr_data)} segments optimized")
    else:
        logger.info("⏭️  Step 3/5: Skipped")

    # ========== Step 4: Translate ==========
    logger.info(f"🌐 Step 4/5: Translating to {target_language}...")
    from core.translator import LLMTranslator
    translator = LLMTranslator(
        model=llm_model,
        target_language=target_language,
        is_reflect=use_reflect,
    )
    asr_data = translator.translate_subtitle(asr_data)
    logger.info(f"   → {len(asr_data)} segments translated")

    # ========== Step 5: Render ASS ==========
    logger.info("📝 Step 5/5: Generating ASS subtitle...")
    # Optimize timing to reduce flicker between adjacent segments
    asr_data.optimize_timing()
    from core.asr_data import SubtitleLayoutEnum
    asr_data.to_ass(save_path=ass_path, layout=SubtitleLayoutEnum.ONLY_TRANSLATE)
    logger.info(f"   → Saved {ass_path}")

    # ========== Step 6: FFmpeg burn-in ==========
    logger.info("🎬 Burning subtitles into video...")
    import subprocess
    ass_escaped = ass_path.replace("\\", "/").replace(":", "\\:")
    cmd = ["ffmpeg", "-i", video_path, "-vf", f"subtitles='{ass_escaped}'", "-c:a", "copy", "-y", output_path]
    result = subprocess.run(cmd, capture_output=True, text=True)
    if result.returncode != 0:
        # Retry without quotes
        cmd = ["ffmpeg", "-i", video_path, "-vf", f"subtitles={ass_escaped}", "-c:a", "copy", "-y", output_path]
        subprocess.run(cmd, capture_output=True, text=True)

    elapsed = time.time() - start_time
    logger.info(f"✅ Done in {elapsed:.1f}s → {output_path}")

    return {
        "status": "success",
        "output_path": output_path,
        "ass_path": ass_path,
        "segments": len(asr_data),
        "elapsed_seconds": round(elapsed, 1),
    }


def main():
    parser = argparse.ArgumentParser(description="Subtitle pipeline CLI")
    parser.add_argument("--video_path", required=True, help="Path to video file")
    parser.add_argument("--config", default="{}", help="JSON config string")
    args = parser.parse_args()

    if not os.path.exists(args.video_path):
        result = {"status": "error", "error": f"Video not found: {args.video_path}"}
        print(json.dumps(result))
        sys.exit(1)

    # Parse config
    try:
        config = json.loads(args.config)
    except json.JSONDecodeError:
        config = {}

    try:
        result = run_pipeline(
            video_path=args.video_path,
            base_url=config.get("base_url", "https://api.deepseek.com/v1"),
            api_key=config.get("api_key", "sk-f6e9fbcd490845cb863e7bd660677c86"),
            llm_model=config.get("llm_model", "deepseek-chat"),
            whisper_model=config.get("whisper_model", "large-v3-turbo"),
            source_language=config.get("source_language", "zh"),
            target_language=config.get("target_language", "越南语"),
            use_reflect=config.get("use_reflect", False),
            skip_optimize=config.get("skip_optimize", False),
            skip_split=config.get("skip_split", False),
        )
        # Output JSON to stdout for Go to parse
        print(json.dumps(result))
    except Exception as e:
        logger.error(f"Pipeline failed: {e}", exc_info=True)
        result = {"status": "error", "error": str(e)}
        print(json.dumps(result))
        sys.exit(1)


if __name__ == "__main__":
    main()
