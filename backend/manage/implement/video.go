package implement

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(global.PathAppChrome),                         // chỉnh lại nếu cần
		chromedp.Flag("headless", false),                                // hiển thị trình duyệt
		chromedp.Flag("disable-blink-features", "AutomationControlled"), // ẩn thuộc tính "navigator.webdriver"
		chromedp.Flag("disable-infobars", true),                         // tắt thanh thông tin "Chrome is being controlled by..."
		chromedp.Flag("start-maximized", true),                          // mở trình duyệt với kích thước tối đa
		chromedp.Flag("disable-dev-shm-usage", true),                    // tránh crash trong môi trường ít tài nguyên
		chromedp.Flag("no-sandbox", true),                               // tránh sandbox errors (nên cân nhắc với bảo mật)
		chromedp.Flag("disable-extensions", true),                       // tắt extension mặc định
		chromedp.Flag("disable-gpu", true),                              // tắt GPU (tùy máy)
		chromedp.UserDataDir(global.PathTempProfile+temdir),             // Thư mục tạm để lưu dữ liệu người dùng
	)

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
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(global.PathAppChrome),                         // chỉnh lại nếu cần
		chromedp.Flag("headless", false),                                // hiển thị trình duyệt
		chromedp.Flag("disable-blink-features", "AutomationControlled"), // ẩn thuộc tính "navigator.webdriver"
		chromedp.Flag("disable-infobars", true),                         // tắt thanh thông tin "Chrome is being controlled by..."
		chromedp.Flag("start-maximized", true),                          // mở trình duyệt với kích thước tối đa
		chromedp.Flag("disable-dev-shm-usage", true),                    // tránh crash trong môi trường ít tài nguyên
		chromedp.Flag("no-sandbox", true),                               // tránh sandbox errors (nên cân nhắc với bảo mật)
		chromedp.Flag("disable-extensions", true),                       // tắt extension mặc định
		chromedp.Flag("disable-gpu", true),                              // tắt GPU (tùy máy)
		chromedp.Flag("lang", "vi-VN"),                                  // đặt ngôn ngữ tiếng Việt
		chromedp.Flag("accept-language", "vi-VN,vi;q=0.9,en;q=0.8"),     // thứ tự ưu tiên ngôn ngữ
		chromedp.UserDataDir(global.PathTempProfile+profile),
	)

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
			// Đọc file thành base64
			fileData, err := ioutil.ReadFile(`videos/` + video) // Đường dẫn đến file video
			if err != nil {
				return err
			}

			base64Data := base64.StdEncoding.EncodeToString(fileData)

			// Inject file data thông qua CDP
			js := fmt.Sprintf(`
				(function() {
					const input = document.querySelector('input[type="file"]');
					if (input) {
						const blob = new Blob([Uint8Array.from(atob('%s'), c => c.charCodeAt(0))], {type: 'video/mp4'});
						const file = new File([blob], 'video.mp4', {type: 'video/mp4'});
						const dt = new DataTransfer();
						dt.items.add(file);
						input.files = dt.files;
						input.dispatchEvent(new Event('change', {bubbles: true}));
					}
				})();
			`, base64Data)

			return chromedp.Evaluate(js, nil).Do(ctx)
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
		log.Fatalf("Lỗi khi điều hướng đến trang chính: %v", err)
		return err
	}

	return nil

}

func (vm *VideoManager) AddVideo(title, videoURL, thumbnailURL string, duration int, likeCount int, profileDouyinID int) (service.Video, error) {
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
	allVideos, err := global.DB.Query(`
		SELECT v.* FROM video AS v
		WHERE v.profile_douyin_id IN (
			SELECT pd.id FROM profile_douyin AS pd
			LEFT JOIN profile AS p ON pd.profile_id = p.id
			WHERE p.id = ?
			AND v.status = 'done'
		)
		LIMIT 10
	`, profile_id)
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

	// xóa file video
	videoPath := fmt.Sprintf("%s/%s.mp4", global.PathVideoReup, video.Title)
	err = os.Remove(videoPath)
	if err != nil {
		log.Printf("Lỗi khi xóa file video: %v", err)
		return err
	}

	// xóa file video sub
	videoPathSub := fmt.Sprintf("%s/%s-sub.mp4", global.PathVideoReup, video.Title)
	err = os.Remove(videoPathSub)
	if err != nil {
		log.Printf("Lỗi khi xóa file video sub: %v", err)
		return err
	}

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
