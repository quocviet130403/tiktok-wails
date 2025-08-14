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

func (a *App) AddProfile(name, hashtag, firstComment string) error {
	err := service.ProfileManager().AddProfile(name, hashtag, firstComment)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) UpdateProfile(id int, name, hashtag, firstComment string) error {
	err := service.ProfileManager().UpdateProfile(id, name, hashtag, firstComment)
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

func (a *App) AddDouyinProfile(accountID int, nickname, url, lastVideoReup string, retryCount int) error {
	err := service.ProfileDouyinManager().AddProfile(accountID, nickname, url, lastVideoReup, retryCount)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) UpdateDouyinProfile(id, accountID, retryCount int, nickname, url, lastVideoReup string) error {
	err := service.ProfileDouyinManager().UpdateProfile(id, accountID, retryCount, nickname, url, lastVideoReup)
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

func (a *App) AddVideo(title string, videoURL string, thumbnailURL string, duration int, likeCount int, profileDouyinID int) error {
	err := service.VideoManager().AddVideo(title, videoURL, thumbnailURL, duration, likeCount, profileDouyinID)
	if err != nil {
		return err
	}
	return nil
}
