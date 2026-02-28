package initialize

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

func InitLogger() {
	// Tạo thư mục logs nếu chưa tồn tại
	home, _ := os.UserHomeDir()
	logDir := filepath.Join(home, "TiktokReupVM", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("Không thể tạo thư mục logs: %v", err)
		return
	}

	// Tạo file log với tên theo ngày
	logFileName := time.Now().Format("2006-01-02") + ".log"
	logFilePath := filepath.Join(logDir, logFileName)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Không thể mở file log: %v", err)
		return
	}

	// Ghi log ra cả file và stdout
	multiWriter := os.Stdout // Go standard log doesn't support io.MultiWriter directly with log.SetOutput
	_ = logFile              // Keep reference

	// Dùng MultiWriter để output ra cả console và file
	log.SetOutput(&dualWriter{stdout: multiWriter, file: logFile})
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Printf("[Logger] Đã khởi tạo — ghi log vào: %s", logFilePath)
}

// dualWriter ghi output ra cả stdout và file
type dualWriter struct {
	stdout *os.File
	file   *os.File
}

func (w *dualWriter) Write(p []byte) (n int, err error) {
	// Ghi ra stdout
	n, err = w.stdout.Write(p)
	if err != nil {
		return n, err
	}
	// Ghi ra file
	_, err = w.file.Write(p)
	return n, err
}
