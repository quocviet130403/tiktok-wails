package implement

import (
	"database/sql"
	"tiktok-wails/backend/manage/service"
)

type ProfileDouyinManager struct {
	db *sql.DB
}

func NewProfileDouyinManager(db *sql.DB) *ProfileDouyinManager {
	return &ProfileDouyinManager{db: db}
}

func (pm *ProfileDouyinManager) UpdateProfile(id, account_id, retry_count int, nickname, url, last_video_reup string) error {
	updateSQL := `UPDATE profile_douyin SET nickname = ?, url = ?, last_video_reup = ? WHERE id = ? AND account_id = ? AND retry_count = ?`
	_, err := pm.db.Exec(updateSQL, nickname, url, last_video_reup, id, account_id, retry_count)
	return err
}

func (pm *ProfileDouyinManager) DeleteProfile(id int) error {
	deleteSQL := `DELETE FROM profile_douyin WHERE id = ?`
	_, err := pm.db.Exec(deleteSQL, id)
	return err
}

func (pm *ProfileDouyinManager) AddProfile(account_id int, nickname, url, last_video_reup string, retry_count int) error {
	insertSQL := `INSERT INTO profile_douyin (account_id, nickname, url, last_video_reup, retry_count) VALUES (?, ?, ?, ?, ?)`
	_, err := pm.db.Exec(insertSQL, account_id, nickname, url, last_video_reup, retry_count)
	return err
}

func (pm *ProfileDouyinManager) GetAllProfiles() ([]service.ProfileDouyin, error) {
	rows, err := pm.db.Query(`SELECT id, nickname, url, last_video_reup, retry_count FROM profile_douyin`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []service.ProfileDouyin
	for rows.Next() {
		var p service.ProfileDouyin
		if err := rows.Scan(&p.ID, &p.Nickname, &p.URL, &p.LastVideoReup, &p.RetryCount); err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}
