package initialize

import (
	"database/sql"
	"tiktok-wails/backend/manage/implement"
	"tiktok-wails/backend/manage/service"
)

func InitManage(db *sql.DB) {
	service.InitVideoManager(implement.NewVideoManager(db))
	service.InitAccountManager(implement.NewAccountManager(db))
}
