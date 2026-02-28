package initialize

import (
	"database/sql"
	"tiktok-wails/backend/manage/implement"
	"tiktok-wails/backend/manage/service"
)

func InitManage(db *sql.DB) {
	service.InitSettingManager(implement.NewSettingManager(db))
	service.InitVideoManager(implement.NewVideoManager(db))
	service.InitProfileManager(implement.NewProfileManager(db))
	service.InitProfileDouyinManager(implement.NewProfileDouyinManager(db))
	service.InitPythonManager(implement.NewPythonManager(db))
}
