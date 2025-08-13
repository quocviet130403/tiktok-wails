package initialize

import "database/sql"

func InitDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "accounts.db")
	if err != nil {
		return nil, err
	}

	// Tạo bảng nếu chưa tồn tại
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS accounts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        url_reup TEXT NOT NULL,
        hashtag TEXT NOT NULL,
        first_comment TEXT NOT NULL,
		last_video_reup TEXT,
		retry_count INTEGER DEFAULT 0,
		is_authenticated BOOLEAN DEFAULT FALSE
    );`

	_, err = db.Exec(createTableSQL)
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
		account_id INTEGER,
		FOREIGN KEY (account_id) REFERENCES accounts (id),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		duration INTEGER,
		like_count INTEGER DEFAULT 0
	);`
	_, err = db.Exec(createVideosTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}
