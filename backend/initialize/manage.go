package initialize

import (
	"database/sql"
	"tiktok-wails/backend/manage"
)

var (
	localVideoManager *manage.VideoManager
)

func VideoManager() *manage.VideoManager {
	return localVideoManager
}

func InitManage(db *sql.DB) {
	localVideoManager = manage.NewVideoManager(db)
}
