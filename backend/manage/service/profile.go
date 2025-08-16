package service

type Profiles struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Hashtag         string `json:"hashtag"`
	FirstComment    string `json:"first_comment"`
	IsAuthenticated bool   `json:"is_authenticated"`
	ProxyIP         string `json:"proxy_ip"`
	ProxyPort       string `json:"proxy_port"`
}

type ProfileManagerInterface interface {
	UpdateProfile(id int, name, hashtag, firstComment, proxy_ip, proxy_port string) error
	DeleteProfile(id int) error
	AddProfile(name, hashtag, firstComment, proxy_ip, proxy_port string) error
	GetAllProfiles() ([]Profiles, error)
}

var (
	localProfileManager ProfileManagerInterface
)

func ProfileManager() ProfileManagerInterface {
	if localProfileManager == nil {
		panic("ProfileManager is not initialized")
	}
	return localProfileManager
}

func InitProfileManager(i ProfileManagerInterface) {
	localProfileManager = i
}
