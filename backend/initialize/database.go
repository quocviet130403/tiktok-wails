package initialize

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const (
	KEY_PATH_CHROME     = "path_chrome"
	VALUE_PATH_CHROME   = "C:/Program Files/Google/Chrome/Application/chrome.exe"
	KEY_SCHEDULE_TIME   = "schedule_time"
	VALUE_SCHEDULE_TIME = "daily"
	KEY_RUN_AT_TIME     = "run_at_time"
	VALUE_RUN_AT_TIME   = "24"
)

func InitDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "wails.db")
	if err != nil {
		return nil, err
	}

	// createProxyTableSQL := `
	// CREATE TABLE IF NOT EXISTS proxies (
	// 	id INTEGER PRIMARY KEY AUTOINCREMENT,
	// 	ip TEXT NOT NULL,
	// 	port INTEGER NOT NULL
	// );`

	// Tạo bảng nếu chưa tồn tại
	createProfileTableSQL := `
    CREATE TABLE IF NOT EXISTS profiles (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL UNIQUE,
        hashtag TEXT NOT NULL,
        first_comment TEXT NOT NULL,
		is_authenticated BOOLEAN DEFAULT FALSE,
		proxy_ip TEXT DEFAULT NULL,
		proxy_port INTEGER DEFAULT NULL
    );`

	_, err = db.Exec(createProfileTableSQL)
	if err != nil {
		return nil, err
	}

	createProfileDouyinTableSQL := `
	CREATE TABLE IF NOT EXISTS profile_douyin (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nickname TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL UNIQUE,
		last_video_reup TEXT DEFAULT NULL,
		retry_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createProfileDouyinTableSQL)
	if err != nil {
		return nil, err
	}

	createProfilesProfileDouyin := `
	CREATE TABLE IF NOT EXISTS profiles_profile_douyin (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_id INTEGER NOT NULL,
		profile_douyin_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createProfilesProfileDouyin)
	if err != nil {
		return nil, err
	}

	// tạo thêm table videos
	createVideosTableSQL := `
	CREATE TABLE IF NOT EXISTS videos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		video_url TEXT NOT NULL,
		thumbnail_url TEXT,
		profile_douyin_id INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		duration INTEGER,
		status TEXT DEFAULT 'pending',
		like_count INTEGER DEFAULT 0
	);`
	_, err = db.Exec(createVideosTableSQL)
	if err != nil {
		return nil, err
	}

	// tao them table settings
	createSettingsTableSQL := `
	CREATE TABLE IF NOT EXISTS settings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key TEXT NOT NULL UNIQUE,
		value TEXT NOT NULL
	);`

	_, err = db.Exec(createSettingsTableSQL)
	if err != nil {
		return nil, err
	}

	// insert data setting
	_, err = db.Exec("INSERT OR IGNORE INTO settings (key, value) VALUES (?, ?)", KEY_PATH_CHROME, VALUE_PATH_CHROME)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("INSERT OR IGNORE INTO settings (key, value) VALUES (?, ?)", KEY_SCHEDULE_TIME, VALUE_SCHEDULE_TIME)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("INSERT OR IGNORE INTO settings (key, value) VALUES (?, ?)", KEY_RUN_AT_TIME, VALUE_RUN_AT_TIME)
	if err != nil {
		return nil, err
	}

	return db, nil
}
