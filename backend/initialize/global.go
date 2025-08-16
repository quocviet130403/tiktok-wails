package initialize

import (
	"database/sql"
	"fmt"
	"os"
	"tiktok-wails/backend/global"
)

const (
	KEY_PATH_CHROME           = "path_chrome"
	VALUE_DEFAULT_PATH_CHROME = "C:/Program Files/Google/Chrome/Application/chrome.exe"
)

func InitGlobal(db *sql.DB) error {
	global.DB = db

	pathChrome, err := db.Query("SELECT value FROM settings WHERE key = ?", KEY_PATH_CHROME)
	if err != nil {
		return err
	}
	defer pathChrome.Close()

	global.PathAppChrome = VALUE_DEFAULT_PATH_CHROME
	if pathChrome.Next() {
		_ = pathChrome.Scan(&global.PathAppChrome)
	}

	fmt.Println("PathChrome1:", global.PathAppChrome)

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
