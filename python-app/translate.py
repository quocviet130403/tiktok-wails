"""Full subtitle pipeline Flask API.

Pipeline: ASR → Split → Optimize → Translate → ASS → FFmpeg

API: POST /translate-video
Body: {
    "video_path": "/path/to/video.mp4",
    "base_url": "https://api.deepseek.com/v1",   (optional, or use env)
    "api_key": "sk-xxx",                          (optional, or use env)
    "llm_model": "deepseek-chat",                 (optional, default deepseek-chat)
    "whisper_model": "large-v3-turbo",            (optional, default large-v3-turbo)
    "source_language": "zh",                      (optional, default zh)
    "target_language": "越南语",                   (optional, default 越南语)
    "use_reflect": false,                          (optional, reflect translation)
    "skip_optimize": false,                        (optional, skip ASR correction)
    "skip_split": false,                           (optional, skip LLM splitting)
}
"""

import logging
import os
import subprocess
import sys
import time

from flask import Flask, jsonify, request

# Setup logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(name)s] %(levelname)s: %(message)s",
    datefmt="%H:%M:%S",
)
logger = logging.getLogger("translate_api")

app = Flask(__name__)


@app.route("/health", methods=["GET"])
def health():
    """Health check endpoint."""
    return jsonify({"status": "ok", "version": "2.0"})


@app.route("/translate-video", methods=["POST"])
def translate_video():
    """Main API endpoint - full subtitle pipeline."""
    data = request.get_json()
    if not data or "video_path" not in data:
        return jsonify({"error": "Missing video_path parameter"}), 400

    video_path = data["video_path"]
    if not os.path.exists(video_path):
        return jsonify({"error": f"Video not found: {video_path}"}), 404

    # Config (request body > env vars > defaults)
    base_url = data.get("base_url", os.getenv("OPENAI_BASE_URL", "https://api.deepseek.com/v1"))
    api_key = data.get("api_key", os.getenv("OPENAI_API_KEY", "sk-f6e9fbcd490845cb863e7bd660677c86"))
    llm_model = data.get("llm_model", "deepseek-chat")
    whisper_model = data.get("whisper_model", "large-v3-turbo")
    source_language = data.get("source_language", "zh")
    target_language = data.get("target_language", "越南语")
    use_reflect = data.get("use_reflect", False)
    skip_optimize = data.get("skip_optimize", False)
    skip_split = data.get("skip_split", False)

    # Init LLM client if credentials provided
    if base_url and api_key:
        from core.llm_client import init_client
        init_client(base_url, api_key)

    # Output paths
    base_name = os.path.splitext(video_path)[0]
    ass_path = f"{base_name}-sub.ass"
    output_path = f"{base_name}-sub.mp4"

    start_time = time.time()

    try:
        # ========== Step 1: ASR ==========
        logger.info(f"🎧 Step 1/5: Transcribing video ({whisper_model})...")
        from core.asr_engine import transcribe_video
        asr_data = transcribe_video(
            video_path=video_path,
            model_size=whisper_model,
            language=source_language,
        )
        logger.info(f"   → {len(asr_data)} word segments, is_word={asr_data.is_word_timestamp()}")

        # ========== Step 2: Split ==========
        if not skip_split:
            logger.info("✂️  Step 2/5: Intelligent sentence splitting...")
            from core.split import SubtitleSplitter
            splitter = SubtitleSplitter(model=llm_model)
            asr_data = splitter.split_subtitle(asr_data)
            logger.info(f"   → {len(asr_data)} sentences after split")
        else:
            logger.info("⏭️  Step 2/5: Skipping split (skip_split=true)")

        # ========== Step 3: Optimize ==========
        if not skip_optimize:
            logger.info("🔧 Step 3/5: Fixing ASR errors with LLM...")
            from core.optimizer import SubtitleOptimizer
            optimizer = SubtitleOptimizer(model=llm_model)
            asr_data = optimizer.optimize_subtitle(asr_data)
            logger.info(f"   → {len(asr_data)} segments optimized")
        else:
            logger.info("⏭️  Step 3/5: Skipping optimize (skip_optimize=true)")

        # ========== Step 4: Translate ==========
        logger.info(f"🌐 Step 4/5: Translating to {target_language} (reflect={use_reflect})...")
        from core.translator import LLMTranslator
        translator = LLMTranslator(
            model=llm_model,
            target_language=target_language,
            is_reflect=use_reflect,
        )
        asr_data = translator.translate_subtitle(asr_data)
        logger.info(f"   → {len(asr_data)} segments translated")

        # ========== Step 5: Render ==========
        logger.info("📝 Step 5/5: Generating ASS subtitle file...")
        from core.asr_data import SubtitleLayoutEnum
        asr_data.to_ass(save_path=ass_path, layout=SubtitleLayoutEnum.ONLY_TRANSLATE)
        logger.info(f"   → Saved to {ass_path}")

        # ========== Step 6: FFmpeg burn-in ==========
        logger.info("🎬 Burning subtitles into video...")
        # Escape paths for ffmpeg subtitles filter
        ass_path_escaped = ass_path.replace("\\", "/").replace(":", "\\:")
        cmd = [
            "ffmpeg", "-i", video_path,
            "-vf", f"subtitles='{ass_path_escaped}'",
            "-c:a", "copy", "-y", output_path,
        ]
        result = subprocess.run(cmd, capture_output=True, text=True)
        if result.returncode != 0:
            logger.warning(f"FFmpeg failed, trying without quotes: {result.stderr[-200:]}")
            # Retry without quotes
            cmd = [
                "ffmpeg", "-i", video_path,
                "-vf", f"subtitles={ass_path_escaped}",
                "-c:a", "copy", "-y", output_path,
            ]
            result = subprocess.run(cmd, capture_output=True, text=True)

        elapsed = time.time() - start_time
        logger.info(f"✅ Done in {elapsed:.1f}s → {output_path}")

        return jsonify({
            "status": "success",
            "output_path": output_path,
            "ass_path": ass_path,
            "segments": len(asr_data),
            "elapsed_seconds": round(elapsed, 1),
        })

    except Exception as e:
        elapsed = time.time() - start_time
        logger.error(f"❌ Pipeline failed after {elapsed:.1f}s: {e}", exc_info=True)
        return jsonify({
            "status": "error",
            "error": str(e),
            "elapsed_seconds": round(elapsed, 1),
        }), 500


if __name__ == "__main__":
    logger.info("Starting subtitle pipeline API on port 9230...")
    app.run(host="0.0.0.0", port=9230, debug=True)