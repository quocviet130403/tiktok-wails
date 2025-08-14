package service

type Profiles struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Hashtag         string `json:"hashtag"`
	FirstComment    string `json:"first_comment"`
	IsAuthenticated bool   `json:"is_authenticated"`
}

type ProfileManagerInterface interface {
	UpdateProfile(id int, name, hashtag, firstComment string) error
	DeleteProfile(id int) error
	AddProfile(name, hashtag, firstComment string) error
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
