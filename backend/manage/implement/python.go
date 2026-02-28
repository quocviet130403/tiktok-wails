package implement

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"tiktok-wails/backend/global"
	"time"
)

type PythonManager struct {
	db *sql.DB
}

func NewPythonManager(db *sql.DB) *PythonManager {
	return &PythonManager{
		db: db,
	}
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

	// Gọi Python Flask API
	requestBody, _ := json.Marshal(map[string]string{
		"video_path": videoPath,
	})

	client := &http.Client{Timeout: 10 * time.Minute} // Translation có thể mất lâu
	resp, err := client.Post("http://localhost:9230/translate-video", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("lỗi khi gọi API translate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API translate trả về status: %d", resp.StatusCode)
	}

	var result struct {
		Status     string `json:"status"`
		OutputPath string `json:"output_path"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("lỗi khi parse response: %w", err)
	}

	log.Printf("Đã dịch video thành công: %s → %s", videoPath, result.OutputPath)
	return nil
}
