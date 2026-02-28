package backend

import (
	"context"
	"fmt"
	"tiktok-wails/backend/manage/service"
)

// App struct
type App struct {
	ctx   context.Context
	ready bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

var errNotReady = fmt.Errorf("server đang khởi tạo, vui lòng chờ...")

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.ready = true // InitServer() đã chạy xong trong main()
}

// Profiles
func (a *App) GetAllProfiles() []service.Profiles {
	if !a.ready {
		return nil
	}
	profiles, err := service.ProfileManager().GetAllProfiles()
	if err != nil {
		return nil
	}
	return profiles
}

func (a *App) AddProfile(name, hashtag, firstComment string, proxyIP string, proxyPort string) error {
	if !a.ready {
		return errNotReady
	}
	return service.ProfileManager().AddProfile(name, hashtag, firstComment, proxyIP, proxyPort)
}

func (a *App) UpdateProfile(id int, name, hashtag, firstComment string, proxyIP string, proxyPort string) error {
	if !a.ready {
		return errNotReady
	}
	return service.ProfileManager().UpdateProfile(id, name, hashtag, firstComment, proxyIP, proxyPort)
}

func (a *App) DeleteProfile(id int) error {
	if !a.ready {
		return errNotReady
	}
	return service.ProfileManager().DeleteProfile(id)
}

// Profile Douyin
func (a *App) GetAllDouyinProfiles() ([]service.ProfileDouyin, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return service.ProfileDouyinManager().GetAllProfiles()
}

func (a *App) AddDouyinProfile(nickname, url string) error {
	if !a.ready {
		return errNotReady
	}
	return service.ProfileDouyinManager().AddProfile(nickname, url)
}

func (a *App) UpdateDouyinProfile(id int, nickname, url string) error {
	if !a.ready {
		return errNotReady
	}
	return service.ProfileDouyinManager().UpdateProfile(id, nickname, url)
}

func (a *App) DeleteDouyinProfile(id int) error {
	if !a.ready {
		return errNotReady
	}
	return service.ProfileDouyinManager().DeleteProfile(id)
}

func (a *App) ToggleHasTranslate(id int) error {
	if !a.ready {
		return errNotReady
	}
	return service.ProfileDouyinManager().ToggleHasTranslate(id)
}

// Videos
func (a *App) GetAllVideos(page int, pageSize int) ([]service.Video, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return service.VideoManager().GetAllVideos(page, pageSize)
}

func (a *App) AddVideo(title string, videoURL string, thumbnailURL string, duration int, likeCount int, profileDouyinID int) (service.Video, error) {
	if !a.ready {
		return service.Video{}, errNotReady
	}
	return service.VideoManager().AddVideo(title, videoURL, thumbnailURL, duration, likeCount, profileDouyinID)
}

// Settings
func (a *App) GetAllSettings() (map[string]string, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return service.SettingManager().GetAllSettings()
}

func (a *App) GetSetting(key string) (string, error) {
	if !a.ready {
		return "", errNotReady
	}
	return service.SettingManager().GetSetting(key)
}

func (a *App) SetSetting(key, value string) error {
	if !a.ready {
		return errNotReady
	}
	return service.SettingManager().SetSetting(key, value)
}

func (a *App) SeederSetting() error {
	if !a.ready {
		return errNotReady
	}
	return service.SettingManager().SetSetting("path_chrome", "C:/Program Files/Google/Chrome/Application/chrome.exe")
}

// Connect profile with douyin profiles
func (a *App) ConnectWithProfileDouyin(profileId int, listProfileDouyinId []int) error {
	if !a.ready {
		return errNotReady
	}
	return service.ProfileManager().ConnectWithProfileDouyin(profileId, listProfileDouyinId)
}

func (a *App) GetAllDouyinProfilesFromProfile(profileId int) ([]service.ProfileDouyin, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return service.ProfileDouyinManager().GetAllProfileDouyinFromProfile(profileId)
}
