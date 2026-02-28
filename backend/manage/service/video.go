package service

type Video struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	VideoURL        string `json:"video_url"`
	ThumbnailURL    string `json:"thumbnail_url"`
	Duration        int    `json:"duration"`
	LikeCount       int    `json:"like_count"`
	ProfileDouyinID int    `json:"profile_douyin_id"`
	Status          string `json:"status"`
	IsDeleted       bool   `json:"is_deleted"`
}

type VideoManagerInterface interface {
	LoginTiktok(temdir string) error
	UploadVideo(profile, video, title string) error
	AddVideo(title, videoURL, thumbnailURL string, duration int, likeCount int, profileDouyinID int) (Video, error)
	GetAllVideos(page int, pageSize int) ([]Video, error)
	GetAllVideosNP() ([]Video, error)
	GetVideoReup(profile_id int) ([]Video, error)
	UpdateStatusReup(video_id, profile_id int, status string) error
	CreateConnectWithProfile(profileID int, videoID int) error
	DeleteVideo(video Video) error
	GetCompleteProfileVideos(video_id int) (int, error)
}

var (
	localVideoManager VideoManagerInterface
)

func VideoManager() VideoManagerInterface {
	if localVideoManager == nil {
		return nil
	}
	return localVideoManager
}

func InitVideoManager(i VideoManagerInterface) {
	localVideoManager = i
}
