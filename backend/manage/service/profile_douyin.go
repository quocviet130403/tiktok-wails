package service

type ProfileDouyin struct {
	ID            int    `json:"id"`
	Nickname      string `json:"nickname"`
	URL           string `json:"url"`
	LastVideoReup string `json:"last_video_reup"`
	RetryCount    int    `json:"retry_count"`
}

type ProfileDouyinInterface interface {
	UpdateProfile(id, account_id, retry_count int, nickname, url, last_video_reup string) error
	DeleteProfile(id int) error
	AddProfile(account_id int, nickname, url, last_video_reup string, retry_count int) error
	GetAllProfiles() ([]ProfileDouyin, error)
	AccessProfile(profile ProfileDouyin) error
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
