package initialize

import (
	"database/sql"
	"os"
	"tiktok-wails/backend/global"
)

func InitGlobal(db *sql.DB) error {
	global.DB = db
	global.PathAppChrome = "C:/Program Files/Google/Chrome/Application/chrome.exe"
	home, _ := os.UserHomeDir()
	global.PathTempProfile = home + "/TiktokReupVM/TempProfile/"
	global.PathVideoReup = home + "/TiktokReupVM/VideoReup/"
	global.PathHandleCaptcha = home + "/TiktokReupVM/c/"

	if _, err := os.Stat(global.PathTempProfile); os.IsNotExist(err) {
		err = os.MkdirAll(global.PathTempProfile, 0755)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(global.PathVideoReup); os.IsNotExist(err) {
		err = os.MkdirAll(global.PathVideoReup, 0755)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(global.PathHandleCaptcha); os.IsNotExist(err) {
		err = os.MkdirAll(global.PathHandleCaptcha, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
