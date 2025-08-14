package service

type Accounts struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	UrlReup         string `json:"url_reup"`
	Hashtag         string `json:"hashtag"`
	FirstComment    string `json:"first_comment"`
	LastVideoReup   string `json:"last_video_reup"`
	RetryCount      int    `json:"retry_count"`
	IsAuthenticated bool   `json:"is_authenticated"`
}

type AccountManagerInterface interface {
	UpdateAccount(id int, name, urlReup, hashtag, firstComment string) error
	DeleteAccount(id int) error
	AddAccount(name, urlReup, hashtag, firstComment string) error
	GetAllAccounts() ([]Accounts, error)
}

var (
	localAccountManager AccountManagerInterface
)

func AccountManager() AccountManagerInterface {
	if localAccountManager == nil {
		panic("AccountManager is not initialized")
	}
	return localAccountManager
}

func InitAccountManager(i AccountManagerInterface) {
	localAccountManager = i
}
