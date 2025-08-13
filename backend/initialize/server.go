package initialize

import (
	"database/sql"
	"tiktok-wails/backend/global"
)

func InitServer(db *sql.DB) error {
	global.DB = db
	return nil
}
