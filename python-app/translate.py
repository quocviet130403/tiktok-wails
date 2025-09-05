from flask import Flask, request, jsonify
import whisper_timestamped as whisper
from transformers import MarianTokenizer, MarianMTModel
import subprocess
import time
import os

app = Flask(__name__)

# --- Load mô hình dịch một lần khi khởi động ---
print("🔁 Đang tải mô hình dịch (zh → vi)...")
model_name = "Helsinki-NLP/opus-mt-zh-vi"
tokenizer = MarianTokenizer.from_pretrained(model_name)
model = MarianMTModel.from_pretrained(model_name)
print("✓ Đã tải xong mô hình dịch")

# --- Hàm dịch ---
def translate_text(text):
    inputs = tokenizer(text, return_tensors="pt", truncation=True, padding=True)
    translated = model.generate(**inputs, max_length=200)
    output = tokenizer.decode(translated[0], skip_special_tokens=True)
    return output

# --- Hàm tạo thời gian theo format ASS ---
def format_time(seconds):
    h = int(seconds // 3600)
    m = int((seconds % 3600) // 60)
    s = int(seconds % 60)
    cs = int((seconds - int(seconds)) * 100)
    return f"{h}:{m:02d}:{s:02d}.{cs:02d}"

# --- Hàm tạo file .ass ---
def create_ass_file(segments, ass_file):
    header = """[Script Info]
Title: Vietnamese Sub
ScriptType: v4.00+
PlayResX: 1280
PlayResY: 720

[V4+ Styles]
Format: Name, Fontname, Fontsize, PrimaryColour, BackColour, OutlineColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding
Style: Default,Arial,28,&H00FFFFFF,&H80000000,&H00000000,-1,0,0,0,100,100,0,0,3,1,0,2,10,10,40,1

[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
"""
    with open(ass_file, "w", encoding="utf-8") as f:
        f.write(header)
        for seg in segments:
            start = format_time(seg['start'])
            end = format_time(seg['end'])
            text = seg['translated'].replace("\n", " ")
            f.write(f"Dialogue: 0,{start},{end},Default,,0,0,0,,{text}\n")

# --- API xử lý video ---
@app.route("/translate-video", methods=["POST"])
def translate_video():
    data = request.get_json()
    if not data or "video_path" not in data:
        return jsonify({"error": "Thiếu tham số video_path"}), 400

    video_path = data["video_path"]
    if not os.path.exists(video_path):
        return jsonify({"error": "Video không tồn tại"}), 404

    output_path = video_path.replace(".mp4", "-sub.mp4")
    ass_path = video_path.replace(".mp4", "-sub.ass")

    # --- Nhận diện tiếng Trung ---
    print(f"🎧 Đang nhận diện tiếng Trung từ video {video_path}...")
    audio = whisper.load_audio(video_path)
    model_whisper = whisper.load_model("base", device="cpu")
    result = whisper.transcribe(model_whisper, audio, language="zh")
    segments = result['segments']

    # --- Dịch ---
    print("🌐 Đang dịch sang tiếng Việt...")
    for i, seg in enumerate(segments):
        print(f"  [{i+1}/{len(segments)}] {seg['text'][:30]}...")
        seg['translated'] = translate_text(seg['text'])
        time.sleep(0.2)

    # --- Ghi file ASS ---
    print("📝 Tạo file phụ đề ASS...")
    create_ass_file(segments, ass_path)

    # --- Ghép phụ đề vào video ---
    print("🎬 Ghép phụ đề vào video...")
    subprocess.run([
        "ffmpeg", "-i", video_path, "-vf", f"subtitles={ass_path}", "-c:a", "copy", "-y", output_path
    ])

    return jsonify({
        "status": "success",
        "output_path": output_path
    })


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=9230, debug=True)