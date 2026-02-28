package initialize

import (
	"log"
	"strconv"
	"sync"
	"tiktok-wails/backend/global"
	"tiktok-wails/backend/manage/service"
	"time"
)

// waitUntilHour tính thời gian chờ đến giờ chạy tiếp theo
func waitUntilHour(hour int) time.Duration {
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())
	if now.After(next) {
		next = next.Add(24 * time.Hour)
	}
	return next.Sub(now)
}

// Job định nghĩa một scheduled job
type Job struct {
	Name string
	Fn   func() error
}

// runJobPipeline chạy danh sách jobs tuần tự — job trước xong thì job sau chạy ngay
func runJobPipeline(jobs []Job) {
	for _, job := range jobs {
		start := time.Now()
		log.Printf("[Scheduler] ▶ Bắt đầu job: %s", job.Name)

		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[Scheduler] ✖ Job '%s' bị panic: %v", job.Name, r)
				}
			}()
			if err := job.Fn(); err != nil {
				log.Printf("[Scheduler] ✖ Job '%s' lỗi: %v", job.Name, err)
			}
		}()

		elapsed := time.Since(start)
		log.Printf("[Scheduler] ✔ Hoàn thành job: %s (mất %v)", job.Name, elapsed.Round(time.Second))
	}
}

func InitSchedule() {
	parsedHour, err := strconv.Atoi(global.ScheduleSetting.RunAtTime)
	if err != nil {
		log.Printf("[Scheduler] Failed to parse VALUE_RUN_AT_TIME: %v", err)
		return
	}
	hour := parsedHour

	if hour < 0 || hour > 23 {
		log.Printf("[Scheduler] Giờ chạy không hợp lệ: %d (phải từ 0-23). Scheduler bị tắt.", hour)
		return
	}

	// Pipeline: 4 jobs chạy tuần tự — job trước xong → job sau chạy ngay
	jobs := []Job{
		{Name: "1. Scrape Douyin Videos", Fn: jobScrapeDouyin},
		{Name: "2. Upload TikTok Videos", Fn: jobUploadTikTok},
		{Name: "3. Delete Completed Videos", Fn: jobDeleteVideos},
		{Name: "4. Check Profile Authentication", Fn: jobCheckAuth},
	}

	go func() {
		for {
			wait := waitUntilHour(hour)
			log.Printf("[Scheduler] Pipeline sẽ chạy sau %v (lúc %d:00)", wait.Round(time.Minute), hour)
			time.Sleep(wait)

			pipelineStart := time.Now()
			log.Printf("[Scheduler] ═══════ BẮT ĐẦU PIPELINE ═══════")

			runJobPipeline(jobs)

			pipelineElapsed := time.Since(pipelineStart)
			log.Printf("[Scheduler] ═══════ KẾT THÚC PIPELINE (tổng: %v) ═══════", pipelineElapsed.Round(time.Second))

			// Chờ ít nhất 1 giờ trước khi check lại, tránh chạy lại cùng ngày
			time.Sleep(1 * time.Hour)
		}
	}()

	log.Printf("[Scheduler] Khởi tạo thành công — pipeline chạy hàng ngày lúc %d:00", hour)
}

// ──────────────────────────────────────────────────────────────
// Job 1: Scrape Video từ Douyin
// ──────────────────────────────────────────────────────────────
func jobScrapeDouyin() error {
	profileDouyins, err := service.ProfileDouyinManager().GetAllProfiles()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, profile := range profileDouyins {
		wg.Add(1)
		go func(profile service.ProfileDouyin) {
			defer wg.Done()
			retry := 0
			for {
				err := service.ProfileDouyinManager().GetVideoFromProfile(profile)
				if err != nil {
					log.Printf("  Lỗi khi lấy video từ profile %s: %v", profile.Nickname, err)
					retry++
					if retry >= 3 {
						log.Printf("  Đã thử lại %d lần, bỏ qua profile %s", retry, profile.Nickname)
						break
					}
					time.Sleep(2 * time.Second)
					continue
				}
				break
			}
		}(profile)
	}
	wg.Wait()
	return nil
}

// ──────────────────────────────────────────────────────────────
// Job 2: Upload video lên TikTok
// ──────────────────────────────────────────────────────────────
func jobUploadTikTok() error {
	profiles, err := service.ProfileManager().GetAllProfiles()
	if err != nil {
		return err
	}

	for _, profile := range profiles {
		allVideos, err := service.VideoManager().GetVideoReup(profile.ID)
		if err != nil {
			log.Printf("  Lỗi khi lấy video từ profile %s: %v", profile.Name, err)
			continue
		}

		for _, video := range allVideos {
			err := service.VideoManager().UploadVideo(profile.Name, video.Title, video.Title+" "+profile.Hashtag)
			if err != nil {
				log.Printf("  Lỗi khi upload video '%s': %v", video.Title, err)
				continue
			}

			err = service.VideoManager().UpdateStatusReup(video.ID, profile.ID, "done")
			if err != nil {
				log.Printf("  Lỗi khi cập nhật trạng thái video: %v", err)
			}
		}
	}
	return nil
}

// ──────────────────────────────────────────────────────────────
// Job 3: Xóa video đã hoàn thành reup
// ──────────────────────────────────────────────────────────────
func jobDeleteVideos() error {
	getAllVideos, err := service.VideoManager().GetAllVideosNP()
	if err != nil {
		return err
	}

	for _, video := range getAllVideos {
		count, err := service.VideoManager().GetCompleteProfileVideos(video.ID)
		if err != nil {
			log.Printf("  Lỗi khi đếm video đã upload: %v", err)
			continue
		}

		if count == 0 {
			err := service.VideoManager().DeleteVideo(video)
			if err != nil {
				log.Printf("  Lỗi khi xóa video: %v", err)
			}
		}
	}
	return nil
}

// ──────────────────────────────────────────────────────────────
// Job 4: Kiểm tra xác thực TikTok profile
// ──────────────────────────────────────────────────────────────
func jobCheckAuth() error {
	profiles, err := service.ProfileManager().GetAllProfileCheckAuthenticated()
	if err != nil {
		return err
	}

	for _, profile := range profiles {
		err := service.VideoManager().LoginTiktok(profile.Name)
		if err != nil {
			log.Printf("  Lỗi khi kiểm tra xác thực profile %s: %v", profile.Name, err)
			updateErr := service.ProfileManager().UpdateAuthenticatedStatus(profile.ID, false)
			if updateErr != nil {
				log.Printf("  Lỗi khi cập nhật trạng thái: %v", updateErr)
			}
		}
	}
	return nil
}
