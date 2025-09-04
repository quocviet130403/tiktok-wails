package backend

import (
	"context"
	"tiktok-wails/backend/manage/service"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Profiles
func (a *App) GetAllProfiles() []service.Profiles {
	profiles, err := service.ProfileManager().GetAllProfiles()
	if err != nil {
		return nil
	}
	return profiles
}

func (a *App) AddProfile(name, hashtag, firstComment string, proxyIP string, proxyPort string) error {
	err := service.ProfileManager().AddProfile(name, hashtag, firstComment, proxyIP, proxyPort)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) UpdateProfile(id int, name, hashtag, firstComment string, proxyIP string, proxyPort string) error {
	err := service.ProfileManager().UpdateProfile(id, name, hashtag, firstComment, proxyIP, proxyPort)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) DeleteProfile(id int) error {
	err := service.ProfileManager().DeleteProfile(id)
	if err != nil {
		return err
	}
	return nil
}

// Profile Douyin
func (a *App) GetAllDouyinProfiles() ([]service.ProfileDouyin, error) {
	profiles, err := service.ProfileDouyinManager().GetAllProfiles()
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

func (a *App) AddDouyinProfile(nickname, url string) error {
	err := service.ProfileDouyinManager().AddProfile(nickname, url)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) UpdateDouyinProfile(id int, nickname, url string) error {
	err := service.ProfileDouyinManager().UpdateProfile(id, nickname, url)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) DeleteDouyinProfile(id int) error {
	err := service.ProfileDouyinManager().DeleteProfile(id)
	if err != nil {
		return err
	}
	return nil
}

// Videos
func (a *App) GetAllVideos(page int, pageSize int) ([]service.Video, error) {
	videos, err := service.VideoManager().GetAllVideos(page, pageSize)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func (a *App) AddVideo(title string, videoURL string, thumbnailURL string, duration int, likeCount int, profileDouyinID int) (service.Video, error) {
	result, err := service.VideoManager().AddVideo(title, videoURL, thumbnailURL, duration, likeCount, profileDouyinID)
	if err != nil {
		return service.Video{}, err
	}
	return result, nil
}

// Settings
func (a *App) GetAllSettings() (map[string]string, error) {
	settings, err := service.SettingManager().GetAllSettings()
	if err != nil {
		return nil, err
	}
	return settings, nil
}

func (a *App) GetSetting(key string) (string, error) {
	value, err := service.SettingManager().GetSetting(key)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (a *App) SetSetting(key, value string) error {
	err := service.SettingManager().SetSetting(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) SeederSetting() error {
	err := service.SettingManager().SetSetting("path_chrome", "C:/Program Files/Google/Chrome/Application/chrome.exe")
	if err != nil {
		return err
	}
	return nil
}

// Connect profile with douyin profiles
func (a *App) ConnectWithProfileDouyin(profileId int, listProfileDouyinId []int) error {
	err := service.ProfileManager().ConnectWithProfileDouyin(profileId, listProfileDouyinId)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) GetAllDouyinProfilesFromProfile(profileId int) ([]service.ProfileDouyin, error) {
	profiles, err := service.ProfileDouyinManager().GetAllProfileDouyinFromProfile(profileId)
	if err != nil {
		return nil, err
	}
	return profiles, nil
}
