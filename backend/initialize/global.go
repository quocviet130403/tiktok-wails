package initialize

import (
	"database/sql"
	"os"
	"tiktok-wails/backend/global"
)

func InitGlobal(db *sql.DB) error {
	global.DB = db

	settings, err := db.Query("SELECT key, value FROM settings")
	if err != nil {
		return err
	}
	defer settings.Close()

	for settings.Next() {
		var key, value string
		if err := settings.Scan(&key, &value); err != nil {
			return err
		}

		switch key {
		case KEY_PATH_CHROME:
			global.PathAppChrome = value
		case KEY_SCHEDULE_TIME:
			global.ScheduleSetting.Time = value
		case KEY_RUN_AT_TIME:
			global.ScheduleSetting.RunAtTime = value
		}
	}

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
