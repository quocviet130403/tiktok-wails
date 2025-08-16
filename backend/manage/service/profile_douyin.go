package service

type ProfileDouyin struct {
	ID            int    `json:"id"`
	Nickname      string `json:"nickname"`
	URL           string `json:"url"`
	LastVideoReup any    `json:"last_video_reup"`
	RetryCount    int    `json:"retry_count"`
}

type ProfileDouyinInterface interface {
	UpdateProfile(id int, nickname, url string) error
	DeleteProfile(id int) error
	AddProfile(nickname, url string) error
	GetAllProfiles() ([]ProfileDouyin, error)
	AccessProfile(profile ProfileDouyin) error
	GetVideoFromProfile(profile ProfileDouyin) error
	UpdateLastVideoReup(id int, lastVideoReup any) error
}

var (
	localProfileDouyinManager ProfileDouyinInterface
)

func ProfileDouyinManager() ProfileDouyinInterface {
	if localProfileDouyinManager == nil {
		panic("ProfileDouyinManager is not initialized")
	}
	return localProfileDouyinManager
}

func InitProfileDouyinManager(i ProfileDouyinInterface) {
	localProfileDouyinManager = i
}
