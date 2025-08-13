package service

type AccountManagerInterface interface {
	UpdateAccount(id int, name, urlReup, hashtag, firstComment string) error
	DeleteAccount(id int) error
	AddAccount(name, urlReup, hashtag, firstComment, lastVideoReup string, retryCount int, isAuthenticated bool) error
}

var (
	localAccountManager AccountManagerInterface
)

func AccountManager() AccountManagerInterface {
	if localAccountManager != nil {
		panic("AccountManager is already initialized")
	}
	return localAccountManager
}

func InitAccountManager(i AccountManagerInterface) {
	localAccountManager = i
}
