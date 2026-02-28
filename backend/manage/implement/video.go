package implement

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"tiktok-wails/backend/global"
	"tiktok-wails/backend/manage/service"
	"time"

	"github.com/chromedp/chromedp"
)

type VideoManager struct {
	db *sql.DB
}

func NewVideoManager(db *sql.DB) *VideoManager {
	return &VideoManager{db: db}
}

func (vm *VideoManager) LoginTiktok(temdir string) error {
	opts := DefaultChromeOpts(global.PathTempProfile+temdir, false)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 300*time.Second) // Tăng timeout một chút
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.tiktok.com/login"),
		chromedp.WaitVisible(`.TUXButton-iconContainer`, chromedp.ByQuery),
	)

	if err != nil {
		return err
	}

	return nil
}

func (vm *VideoManager) UploadVideo(profile, video, title string) error {
	opts := VisibleChromeOpts(global.PathTempProfile + profile)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 90*time.Second) // Tăng timeout một chút
	defer cancel()

	// --- BƯỚC 1: Điều hướng đến trang chính ---
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.tiktok.com/tiktokstudio/upload?from=webapp"),
		chromedp.Sleep(5*time.Second),
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	log.Println("=== start ===")
		// 	return nil
		// }),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Sử dụng SetUploadFiles thay vì base64 để tránh tốn RAM
			videoPath, err := filepath.Abs(filepath.Join("videos", video))
			if err != nil {
				return fmt.Errorf("lỗi khi lấy đường dẫn video: %w", err)
			}

			// Kiểm tra file tồn tại
			if _, err := os.Stat(videoPath); os.IsNotExist(err) {
				return fmt.Errorf("file video không tồn tại: %s", videoPath)
			}

			return chromedp.SetUploadFiles(`input[type="file"]`, []string{videoPath}, chromedp.ByQuery).Do(ctx)
		}),

		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	log.Println("=== done ===")
		// 	return nil
		// }),

		chromedp.ActionFunc(func(ctx context.Context) error {
			// log.Println("=== Điền caption ===")

			// Tìm và click vào editor
			return chromedp.Click(`.public-DraftEditor-content`, chromedp.ByQuery).Do(ctx)
		}),

		chromedp.Sleep(1*time.Second),

		chromedp.ActionFunc(func(ctx context.Context) error {
			// Xóa nội dung cũ và điền nội dung mới
			js := fmt.Sprintf(`
				(function() {
					const editor = document.querySelector('.public-DraftEditor-content');
					if (editor) {
						// Focus vào editor
						editor.focus();
						
						// Xóa tất cả nội dung
						document.execCommand('selectAll');
						document.execCommand('delete');
						
						// Điền nội dung mới
						document.execCommand('insertText', false, '%s');
						
						// Trigger events
						editor.dispatchEvent(new Event('input', {bubbles: true}));
						editor.dispatchEvent(new Event('change', {bubbles: true}));
					}
				})();
			`, title)

			return chromedp.Evaluate(js, nil).Do(ctx)
		}),

		chromedp.Sleep(2*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Thiếu click button submit để đăng video
			return chromedp.Click(`.button-group > button`, chromedp.ByQuery).Do(ctx)
		}),

		chromedp.Sleep(3*time.Second),
	)

	if err != nil {
		log.Printf("Lỗi khi upload video: %v", err)
		return err
	}

	return nil

}

func (vm *VideoManager) AddVideo(title, videoURL, thumbnailURL string, duration int, likeCount int, profileDouyinID int) (service.Video, error) {
	// Kiểm tra video đã tồn tại chưa (tránh duplicate)
	var existingID int
	err := vm.db.QueryRow(`SELECT id FROM videos WHERE video_url = ?`, videoURL).Scan(&existingID)
	if err == nil {
		// Video đã tồn tại, trả về video hiện có
		log.Printf("Video đã tồn tại với ID: %d, bỏ qua", existingID)
		return service.Video{
			ID:              existingID,
			Title:           title,
			VideoURL:        videoURL,
			ThumbnailURL:    thumbnailURL,
			Duration:        duration,
			LikeCount:       likeCount,
			ProfileDouyinID: profileDouyinID,
		}, nil
	}

	insertSQL := `INSERT INTO videos (title, video_url, thumbnail_url, duration, like_count, profile_douyin_id) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := vm.db.Exec(insertSQL, title, videoURL, thumbnailURL, duration, likeCount, profileDouyinID)
	if err != nil {
		return service.Video{}, err
	}

	videoID, err := result.LastInsertId()
	if err != nil {
		return service.Video{}, err
	}

	return service.Video{
		ID:              int(videoID),
		Title:           title,
		VideoURL:        videoURL,
		ThumbnailURL:    thumbnailURL,
		Duration:        duration,
		LikeCount:       likeCount,
		ProfileDouyinID: profileDouyinID,
	}, nil
}

func (vm *VideoManager) GetAllVideos(page int, pageSize int) ([]service.Video, error) {
	offset := (page - 1) * pageSize
	query := `SELECT id, title, video_url, thumbnail_url, duration, like_count, profile_douyin_id, status FROM videos LIMIT ? OFFSET ?`
	rows, err := vm.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []service.Video
	for rows.Next() {
		var video service.Video
		if err := rows.Scan(&video.ID, &video.Title, &video.VideoURL, &video.ThumbnailURL,
			&video.Duration, &video.LikeCount, &video.ProfileDouyinID, &video.Status); err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (vm *VideoManager) GetAllVideosNP() ([]service.Video, error) {
	query := `SELECT id, title, video_url, thumbnail_url, duration, like_count, profile_douyin_id, status, is_deleted_video FROM videos WHERE is_deleted_video = false AND status = 'done'`
	rows, err := vm.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []service.Video
	for rows.Next() {
		var video service.Video
		if err := rows.Scan(&video.ID, &video.Title, &video.VideoURL, &video.ThumbnailURL,
			&video.Duration, &video.LikeCount, &video.ProfileDouyinID, &video.Status, &video.IsDeleted); err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (vm *VideoManager) GetVideoReup(profile_id int) ([]service.Video, error) {
	allVideos, err := vm.db.Query(`
		SELECT v.id, v.title, v.video_url, v.thumbnail_url, v.duration, v.like_count, v.profile_douyin_id, v.status
		FROM videos AS v
		WHERE v.status = 'pending'
		AND v.is_deleted_video = false
		AND v.profile_douyin_id IN (
			SELECT ppd.profile_douyin_id FROM profiles_profile_douyin AS ppd
			WHERE ppd.profile_id = ?
		)
		AND v.id NOT IN (
			SELECT vp.video_id FROM videos_profiles AS vp
			WHERE vp.profile_id = ? AND vp.status = 'done'
		)
		LIMIT 10
	`, profile_id, profile_id)
	if err != nil {
		return nil, err
	}
	defer allVideos.Close()

	var videos []service.Video
	for allVideos.Next() {
		var video service.Video
		if err := allVideos.Scan(&video.ID, &video.Title, &video.VideoURL, &video.ThumbnailURL,
			&video.Duration, &video.LikeCount, &video.ProfileDouyinID, &video.Status); err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (vm *VideoManager) UpdateStatusReup(video_id, profile_id int, status string) error {
	updateSQL := `UPDATE videos_profiles SET status = ? WHERE video_id = ? AND profile_id = ?`
	_, err := vm.db.Exec(updateSQL, status, video_id, profile_id)
	return err
}

func (vm *VideoManager) CreateConnectWithProfile(profileID int, videoID int) error {
	createSQL := `INSERT INTO videos_profiles (video_id, profile_id) VALUES (?, ?)`
	_, err := vm.db.Exec(createSQL, videoID, profileID)
	return err
}

func (vm *VideoManager) DeleteVideo(video service.Video) error {
	deleteSQL := `UPDATE videos SET is_deleted_video = true WHERE id = ?`
	_, err := vm.db.Exec(deleteSQL, video.ID)
	if err != nil {
		return err
	}

	// xóa file video
	videoPath := fmt.Sprintf("%s/%s.mp4", global.PathVideoReup, video.Title)
	err = os.Remove(videoPath)
	if err != nil {
		log.Printf("Lỗi khi xóa file video: %v", err)
		return err
	}

	// // xóa file video sub
	// videoPathSub := fmt.Sprintf("%s/%s-sub.mp4", global.PathVideoReup, video.Title)
	// err = os.Remove(videoPathSub)
	// if err != nil {
	// 	log.Printf("Lỗi khi xóa file video sub: %v", err)
	// 	return err
	// }

	return err
}

func (vm *VideoManager) GetCompleteProfileVideos(video_id int) (int, error) {
	query := `
		SELECT count(v.id) FROM videos v
		LEFT JOIN videos_profiles vp ON v.id = vp.video_id
		WHERE vp.status = 'pending'
		AND v.id = ?;
	`
	rows, err := vm.db.Query(query, video_id)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
	}

	return count, nil
}
