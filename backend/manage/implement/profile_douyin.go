package implement

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"tiktok-wails/backend/global"
	"tiktok-wails/backend/manage/service"
	"time"

	"github.com/chromedp/chromedp"
)

type ProfileDouyinManager struct {
	db *sql.DB
}

func NewProfileDouyinManager(db *sql.DB) *ProfileDouyinManager {
	return &ProfileDouyinManager{db: db}
}

func (pm *ProfileDouyinManager) UpdateProfile(id int, nickname, url string) error {
	updateSQL := `UPDATE profile_douyin SET nickname = ?, url = ? WHERE id = ?`
	_, err := pm.db.Exec(updateSQL, nickname, url, id)
	return err
}

func (pm *ProfileDouyinManager) DeleteProfile(id int) error {
	deleteSQL := `DELETE FROM profile_douyin WHERE id = ?`
	_, err := pm.db.Exec(deleteSQL, id)
	return err
}

func (pm *ProfileDouyinManager) AddProfile(nickname, url string) error {
	tx, err := pm.db.Begin()
	if err != nil {
		return fmt.Errorf("lỗi khi bắt đầu transaction: %w", err)
	}

	// Defer rollback để đảm bảo rollback nếu có lỗi
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	// Thử đăng nhập TikTok
	err = pm.AccessProfile(service.ProfileDouyin{
		Nickname: nickname,
		URL:      url,
	})
	if err != nil {
		return err
	}

	// Thực hiện insert trong transaction
	insertSQL := `INSERT INTO profile_douyin (nickname, url) VALUES (?, ?)`
	_, err = tx.Exec(insertSQL, nickname, url)
	if err != nil {
		return fmt.Errorf("lỗi khi thêm tài khoản: %w", err)
	}

	// Nếu đăng nhập thành công, commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("lỗi khi commit transaction: %w", err)
	}

	return nil
}

func (pm *ProfileDouyinManager) GetAllProfiles() ([]service.ProfileDouyin, error) {
	rows, err := pm.db.Query(`SELECT id, nickname, url, last_video_reup, retry_count FROM profile_douyin`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []service.ProfileDouyin
	for rows.Next() {
		var p service.ProfileDouyin
		if err := rows.Scan(&p.ID, &p.Nickname, &p.URL, &p.LastVideoReup, &p.RetryCount); err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

func (pm *ProfileDouyinManager) AccessProfile(profile service.ProfileDouyin) error {

	profile.Nickname = time.Now().Format("20060102150405")

	defer func() {
		captchaFiles := []string{
			global.PathHandleCaptcha + `/bg-` + profile.Nickname + `.png`,
			global.PathHandleCaptcha + `/slide-` + profile.Nickname + `.png`,
			global.PathHandleCaptcha + `/result-` + profile.Nickname + `.png`,
		}

		for _, file := range captchaFiles {
			if err := os.Remove(file); err != nil {
				// Only log if file exists but couldn't be removed
				if !os.IsNotExist(err) {
					log.Printf("Không thể xóa file %s: %v", file, err)
				}
			} else {
				log.Printf("Đã xóa file: %s", file)
			}
		}
	}()

	tempDir, err := os.MkdirTemp("", "douyin-test-access-*")
	if err != nil {
		return fmt.Errorf("lỗi khi tạo thư mục tạm: %w", err)
	}
	defer os.RemoveAll(tempDir)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath("C:/Program Files/Google/Chrome/Application/chrome.exe"),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("start-maximized", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("enable-logging", true),
		chromedp.Flag("log-level", "0"),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-features", "VizDisplayCompositor"),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.UserDataDir(tempDir),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	log.Println("Bắt đầu tự động hóa Douyin...")

	// Sử dụng WaitGroup để quản lý các goroutine
	var wg sync.WaitGroup

	err = chromedp.Run(ctx,
		chromedp.Navigate(profile.URL),
		chromedp.Sleep(5*time.Second), // Tăng thời gian chờ cho headless mode
	)
	if err != nil {
		return fmt.Errorf("lỗi khi điều hướng: %w", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		handleCaptchaAsync(ctx, profile.Nickname)
	}()

	// Tạo một goroutine để scroll trang định kỳ nhằm trigger load video
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	ticker := time.NewTicker(5 * time.Second)
	// 	defer ticker.Stop()

	// 	for i := 0; i < 12; i++ { // Scroll 12 lần trong 60 giây
	// 		select {
	// 		case <-ctx.Done():
	// 			log.Println("Context đã bị hủy, dừng scroll")
	// 			return
	// 		case <-ticker.C:
	// 			err := chromedp.Run(ctx,
	// 				chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight);`, nil),
	// 				chromedp.Sleep(1*time.Second),
	// 			)
	// 			if err != nil {
	// 				log.Printf("Lỗi khi scroll: %v", err)
	// 			} else {
	// 				log.Printf("Scroll %d/12", i+1)
	// 			}
	// 		}
	// 	}
	// 	log.Println("Hoàn thành scroll trang")
	// }()

	// --- BƯỚC 6: Thêm một goroutine để kiểm tra trang đã load xong chưa ---
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Chờ trang load xong
		err := chromedp.Run(ctx,
			chromedp.WaitVisible("#user_detail_element", chromedp.ByQuery),
			chromedp.Sleep(2*time.Second),
		)
		if err != nil {
			log.Printf("Lỗi khi chờ trang load: %v", err)
			return
		} else {
			log.Println("Trang đã load xong")
		}
	}()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Tất cả các task đã hoàn thành")
	case <-time.After(60 * time.Second):
		log.Println("Timeout sau 60 giây")
		return fmt.Errorf("timeout khi crawl Douyin")
	case <-ctx.Done():
		log.Println("Context bị hủy")
		return ctx.Err()
	}

	log.Println("Kết thúc quá trình crawl video từ Douyin")
	return nil
}
