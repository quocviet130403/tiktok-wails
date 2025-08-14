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

// Accounts
func (a *App) GetAllAccounts() []service.Profiles {
	accounts, err := service.ProfileManager().GetAllProfiles()
	if err != nil {
		return nil
	}
	return accounts
}

func (a *App) AddAccount(name, hashtag, firstComment string) error {
	err := service.ProfileManager().AddProfile(name, hashtag, firstComment)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) UpdateAccount(id int, name, hashtag, firstComment string) error {
	err := service.ProfileManager().UpdateProfile(id, name, hashtag, firstComment)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) DeleteAccount(id int) error {
	err := service.ProfileManager().DeleteProfile(id)
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
