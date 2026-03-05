package implement

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"runtime"
	"tiktok-wails/backend/global"
)

type PythonManager struct {
	db *sql.DB
}

func NewPythonManager(db *sql.DB) *PythonManager {
	return &PythonManager{
		db: db,
	}
}

// getLLMSetting lấy setting từ database, trả về empty string nếu không có
func (pm *PythonManager) getLLMSetting(key string) string {
	var value string
	err := pm.db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err != nil {
		return ""
	}
	return value
}

// getPythonCmd trả về python command phù hợp với OS
func getPythonCmd() string {
	if runtime.GOOS == "windows" {
		return "python"
	}
	return "python3"
}

// getScriptPath trả về đường dẫn đến translate_cli.py
func getScriptPath() string {
	return filepath.Join("python-app", "translate_cli.py")
}

func (pm *PythonManager) TranslateVideo(id int) error {
	// Lấy thông tin video từ database
	var title string
	err := pm.db.QueryRow("SELECT title FROM videos WHERE id = ?", id).Scan(&title)
	if err != nil {
		return fmt.Errorf("không tìm thấy video với ID %d: %w", id, err)
	}

	// Tạo đường dẫn video
	videoPath := fmt.Sprintf("%s/%s.mp4", global.PathVideoReup, title)

	// Lấy LLM config từ settings
	config := map[string]interface{}{
		"base_url":  pm.getLLMSetting("llm_base_url"),
		"api_key":   pm.getLLMSetting("llm_api_key"),
		"llm_model": pm.getLLMSetting("llm_model"),
	}
	if config["llm_model"] == "" {
		config["llm_model"] = "deepseek-chat"
	}

	configJSON, _ := json.Marshal(config)

	// Gọi Python CLI trực tiếp (không cần Flask server)
	log.Printf("[Python] Bắt đầu dịch video: %s", videoPath)

	cmd := exec.Command(
		getPythonCmd(),
		getScriptPath(),
		"--video_path", videoPath,
		"--config", string(configJSON),
	)

	// Capture stdout (JSON result) và stderr (logs)
	stdout, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// stderr chứa logs
			return fmt.Errorf("Python pipeline lỗi: %s\nLogs: %s", err, string(exitErr.Stderr))
		}
		return fmt.Errorf("không thể chạy Python: %w", err)
	}

	// Parse JSON result từ stdout
	var result struct {
		Status         string  `json:"status"`
		OutputPath     string  `json:"output_path"`
		ASSPath        string  `json:"ass_path"`
		Segments       int     `json:"segments"`
		ElapsedSeconds float64 `json:"elapsed_seconds"`
		Error          string  `json:"error"`
	}

	if err := json.Unmarshal(stdout, &result); err != nil {
		return fmt.Errorf("không thể parse kết quả Python: %w\nOutput: %s", err, string(stdout))
	}

	if result.Status == "error" {
		return fmt.Errorf("pipeline lỗi: %s", result.Error)
	}

	log.Printf("[Python] Đã dịch video thành công: %s → %s (%d segments, %.1fs)",
		videoPath, result.OutputPath, result.Segments, result.ElapsedSeconds)
	return nil
}
