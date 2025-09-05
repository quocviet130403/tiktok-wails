package implement

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"tiktok-wails/backend/global"
	"tiktok-wails/backend/manage/service"
	"time"

	"github.com/chromedp/cdproto/network"
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
	rows, err := pm.db.Query(`SELECT id, nickname, url, last_video_reup, retry_count, has_translate FROM profile_douyin`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []service.ProfileDouyin
	for rows.Next() {
		var p service.ProfileDouyin
		if err := rows.Scan(&p.ID, &p.Nickname, &p.URL, &p.LastVideoReup, &p.RetryCount, &p.HasTranslate); err != nil {
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

func (pm *ProfileDouyinManager) UpdateLastVideoReup(id int, lastVideoReup any) error {
	_, err := pm.db.Query("UPDATE profile_douyin SET last_video_reup = ? WHERE id = ?", lastVideoReup, id)
	if err != nil {
		return fmt.Errorf("lỗi khi cập nhật last_video_reup: %w", err)
	}
	return nil
}

func (pm *ProfileDouyinManager) GetVideoFromProfile(profile service.ProfileDouyin) error {

	profile.Nickname = time.Now().Format("20060102150405")

	defer func() {
		captchaFiles := []string{
			`./captcha/bg-` + profile.Nickname + `.png`,
			`./captcha/slide-` + profile.Nickname + `.png`,
			`./captcha/result-` + profile.Nickname + `.png`,
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
		// Thêm các flag để debugging tốt hơn trong headless mode
		chromedp.Flag("enable-logging", true),
		chromedp.Flag("log-level", "0"),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-features", "VizDisplayCompositor"),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.UserDataDir(`C:\Users\viet1\AppData\Local\Temp\`+profile.Nickname),
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
	var networkEnabled bool
	var networkMutex sync.Mutex

	// --- BƯỚC 1: Enable network domain trước khi bắt đầu ---
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			if err := network.Enable().Do(ctx); err != nil {
				return err
			}
			networkMutex.Lock()
			networkEnabled = true
			networkMutex.Unlock()
			log.Println("Network domain đã được enable")
			return nil
		}),
	)
	if err != nil {
		return fmt.Errorf("lỗi khi enable network: %w", err)
	}

	// --- BƯỚC 2: Thiết lập network listener với retry mechanism ---
	log.Println("Bắt đầu lắng nghe network requests...")

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventResponseReceived:
			if ev.Response.URL != "" && strings.Contains(ev.Response.URL, "https://www.douyin.com/aweme/v1/web/aweme/post/") {
				// log.Printf("Phát hiện API call: %s", ev.Response.URL)
				wg.Add(1)
				go func(requestID network.RequestID) {
					defer wg.Done()

					// Đảm bảo network đã được enable
					networkMutex.Lock()
					enabled := networkEnabled
					networkMutex.Unlock()

					if !enabled {
						log.Println("Network chưa được enable, bỏ qua request này")
						return
					}

					// Chờ một chút để đảm bảo response đã hoàn thành
					time.Sleep(500 * time.Millisecond)

					// Retry mechanism cho việc lấy response body
					var data []byte
					var err error
					maxRetries := 3

					for i := 0; i < maxRetries; i++ {
						// Tạo context riêng với timeout ngắn hơn cho mỗi lần thử
						requestCtx, requestCancel := context.WithTimeout(ctx, 10*time.Second)

						err = chromedp.Run(requestCtx, chromedp.ActionFunc(func(ctx context.Context) error {
							data, err = network.GetResponseBody(requestID).Do(ctx)
							return err
						}))

						requestCancel()

						if err == nil && len(data) > 0 {
							break
						}

						log.Printf("Lần thử %d/%d lấy response body thất bại: %v", i+1, maxRetries, err)

						// Chờ một chút trước khi thử lại
						if i < maxRetries-1 {
							time.Sleep(time.Duration(i+1) * time.Second)
						}
					}

					if err != nil {
						log.Printf("Không thể lấy response body sau %d lần thử: %v", maxRetries, err)
						return
					}

					if len(data) == 0 {
						log.Println("Response body rỗng")
						return
					}

					var videoList ListVideos
					if err := json.Unmarshal(data, &videoList); err != nil {
						log.Printf("Lỗi khi parse JSON: %v", err)
						return
					}

					log.Printf("Đã nhận được %d video từ API", len(videoList.AwemeLists))
					var newVideoList []AwemeLists
					for _, video := range videoList.AwemeLists {
						if video.Desc == profile.LastVideoReup {
							break
						}
						newVideoList = append(newVideoList, video)
					}

					for i, video := range newVideoList {
						videoAdded, err := service.VideoManager().AddVideo(video.Desc, video.Video.PlayAddr.URLList[0], video.Video.Cover.URLList[0], video.Duration, video.Statistic.LikeCount, profile.ID)
						if err != nil {
							log.Printf("Lỗi khi thêm video %d: %v", i+1, err)
						}

						err = service.VideoManager().CreateConnectWithProfile(profile.ID, videoAdded.ID)
						if err != nil {
							log.Printf("Lỗi khi tạo kết nối giữa video và profile: %v", err)
							return
						}
					}

					err = service.ProfileDouyinManager().UpdateLastVideoReup(profile.ID, newVideoList[0].Desc)
					if err != nil {
						log.Printf("Lỗi khi cập nhật video đã reup: %v", err)
						return
					}
				}(ev.RequestID)
			}
		case *network.EventLoadingFinished:
			// Log khi request hoàn thành để debugging
			// log.Printf("Request hoàn thành: %s", ev.RequestID)
		}
	})

	// --- BƯỚC 3: Điều hướng đến trang ---
	err = chromedp.Run(ctx,
		chromedp.Navigate(profile.URL),
		chromedp.Sleep(5*time.Second), // Tăng thời gian chờ cho headless mode
	)
	if err != nil {
		return fmt.Errorf("lỗi khi điều hướng: %w", err)
	}

	// --- BƯỚC 4: Kiểm tra và xử lý CAPTCHA song song ---
	wg.Add(1)
	go func() {
		defer wg.Done()
		handleCaptchaAsync(ctx, profile.Nickname)
	}()

	// --- BƯỚC 5: Giữ kết nối để lắng nghe network requests ---
	log.Println("Đang lắng nghe network requests trong 60 giây...")

	// Tạo một goroutine để scroll trang định kỳ nhằm trigger load video
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for i := 0; i < 12; i++ { // Scroll 12 lần trong 60 giây
			select {
			case <-ctx.Done():
				log.Println("Context đã bị hủy, dừng scroll")
				return
			case <-ticker.C:
				err := chromedp.Run(ctx,
					chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight);`, nil),
					chromedp.Sleep(1*time.Second),
				)
				if err != nil {
					log.Printf("Lỗi khi scroll: %v", err)
				} else {
					log.Printf("Scroll %d/12", i+1)
				}
			}
		}
		log.Println("Hoàn thành scroll trang")
	}()

	// --- BƯỚC 6: Thêm một goroutine để kiểm tra trang đã load xong chưa ---
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Chờ trang load xong
		err := chromedp.Run(ctx,
			chromedp.WaitReady("body", chromedp.ByQuery),
			chromedp.Sleep(2*time.Second),
		)
		if err != nil {
			log.Printf("Lỗi khi chờ trang load: %v", err)
		} else {
			log.Println("Trang đã load xong")
		}
	}()

	// Chờ tất cả goroutine hoàn thành hoặc timeout
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
	case <-ctx.Done():
		log.Println("Context bị hủy")
	}

	log.Println("Kết thúc quá trình crawl video từ Douyin")
	return nil
}

func (pm *ProfileDouyinManager) GetAllProfileDouyinFromProfile(profileId int) ([]service.ProfileDouyin, error) {
	query := `
	SELECT pd.id, pd.nickname, pd.url
	FROM profile_douyin pd
	JOIN profiles_profile_douyin ppd ON pd.id = ppd.profile_douyin_id
	WHERE ppd.profile_id = ?`
	rows, err := pm.db.Query(query, profileId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []service.ProfileDouyin
	for rows.Next() {
		var profile service.ProfileDouyin
		if err := rows.Scan(&profile.ID, &profile.Nickname, &profile.URL); err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func (pm *ProfileDouyinManager) ToggleHasTranslate(id int) error {
	var currentValue bool
	err := pm.db.QueryRow("SELECT has_translate FROM profile_douyin WHERE id = ?", id).Scan(&currentValue)
	if err != nil {
		return err
	}

	newValue := !currentValue
	_, err = pm.db.Exec("UPDATE profile_douyin SET has_translate = ? WHERE id = ?", newValue, id)
	return err
}
