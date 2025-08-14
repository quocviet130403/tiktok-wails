package service

type Video struct {
	ID           int
	Title        string
	VideoURL     string
	ThumbnailURL string
	Duration     int
	LikeCount    int
	AccountID    int
	Status       string
}

type VideoManagerInterface interface {
	LoginTiktok(temdir string) error
	UploadVideo(profile, video, title string) error
	AddVideo(title, videoURL, thumbnailURL string, duration int, likeCount int, accountID int) error
	GetAllVideos(page int, pageSize int) ([]Video, error)
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
