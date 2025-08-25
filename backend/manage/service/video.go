package service

type Video struct {
	ID              int
	Title           string
	VideoURL        string
	ThumbnailURL    string
	Duration        int
	LikeCount       int
	ProfileDouyinID int
	Status          string
	IsDeleted       bool
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
		panic("VideoManager is not initialized")
	}
	return localVideoManager
}

func InitVideoManager(i VideoManagerInterface) {
	localVideoManager = i
}
