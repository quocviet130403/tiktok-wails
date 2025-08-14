package initialize

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func InitDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "wails.db")
	if err != nil {
		return nil, err
	}

	// Tạo bảng nếu chưa tồn tại
	createProfileTableSQL := `
    CREATE TABLE IF NOT EXISTS profiles (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        hashtag TEXT NOT NULL,
        first_comment TEXT NOT NULL,
		is_authenticated BOOLEAN DEFAULT FALSE
    );`

	_, err = db.Exec(createProfileTableSQL)
	if err != nil {
		return nil, err
	}

	createProfileDouyinTableSQL := `
	CREATE TABLE IF NOT EXISTS profile_douyin (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		account_id INTEGER,
		nickname TEXT NOT NULL,
		url TEXT NOT NULL,
		last_video_reup TEXT,
		retry_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createProfileDouyinTableSQL)
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

	return db, nil
}
