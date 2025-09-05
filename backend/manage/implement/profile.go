package implement

import (
	"database/sql"
	"fmt"
	"tiktok-wails/backend/manage/service"
)

type ProfileManager struct {
	db *sql.DB
}

func NewProfileManager(db *sql.DB) *ProfileManager {
	return &ProfileManager{db: db}
}

func (am *ProfileManager) UpdateProfile(id int, name, hashtag, firstComment, proxy_ip, proxy_port string) error {
	updateSQL := `UPDATE profiles SET name=?, hashtag=?, first_comment=? WHERE id=? AND proxy_ip=? AND proxy_port=?`
	_, err := am.db.Exec(updateSQL, name, hashtag, firstComment, id, proxy_ip, proxy_port)
	return err
}

// Hàm xóa tài khoản
func (am *ProfileManager) DeleteProfile(id int) error {
	deleteSQL := `DELETE FROM profiles WHERE id=?`
	_, err := am.db.Exec(deleteSQL, id)
	return err
}

func (am *ProfileManager) AddProfile(name, hashtag, firstComment, proxy_ip, proxy_port string) error {
	// Bắt đầu transaction
	tx, err := am.db.Begin()
	if err != nil {
		return fmt.Errorf("lỗi khi bắt đầu transaction: %w", err)
	}

	// Defer rollback để đảm bảo rollback nếu có lỗi
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	isAuthenticated := false
	// Thử đăng nhập TikTok
	err = service.VideoManager().LoginTiktok(name)
	if err == nil {
		isAuthenticated = true
	} else {
		fmt.Println("Log err:", err)
	}

	// Thực hiện insert trong transaction
	insertSQL := `INSERT INTO profiles (name, hashtag, first_comment, is_authenticated, proxy_ip, proxy_port) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(insertSQL, name, hashtag, firstComment, isAuthenticated, proxy_ip, proxy_port)
	if err != nil {
		return fmt.Errorf("lỗi khi thêm tài khoản: %w", err)
	}

	// Nếu đăng nhập thành công, commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("lỗi khi commit transaction: %w", err)
	}

	return nil
}

func (am *ProfileManager) GetAllProfiles() ([]service.Profiles, error) {
	rows, err := am.db.Query("SELECT id, name, hashtag, first_comment, is_authenticated, proxy_ip, proxy_port FROM profiles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []service.Profiles

	for rows.Next() {
		var account service.Profiles
		err := rows.Scan(&account.ID, &account.Name, &account.Hashtag, &account.FirstComment, &account.IsAuthenticated, &account.ProxyIP, &account.ProxyPort)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}

func (am *ProfileManager) GetAllProfileCheckAuthenticated() ([]service.Profiles, error) {
	rows, err := am.db.Query("SELECT id, name, is_authenticated FROM profiles WHERE is_authenticated = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []service.Profiles

	for rows.Next() {
		var account service.Profiles
		err := rows.Scan(&account.ID, &account.Name, &account.IsAuthenticated)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}

func (am *ProfileManager) ConnectWithProfileDouyin(profileId int, listProfileDouyinId []int) error {
	existingProfiles, err := service.ProfileDouyinManager().GetAllProfileDouyinFromProfile(profileId)
	if err != nil {
		return err
	}

	existingMap := make(map[int]struct{}, len(existingProfiles))
	for _, p := range existingProfiles {
		existingMap[p.ID] = struct{}{}
	}

	newMap := make(map[int]struct{}, len(listProfileDouyinId))
	for _, id := range listProfileDouyinId {
		newMap[id] = struct{}{}
	}

	tx, err := am.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	for id := range newMap {
		if _, exists := existingMap[id]; !exists {
			_, err := tx.Exec(`INSERT INTO profiles_profile_douyin (profile_id, profile_douyin_id) VALUES (?, ?)`, profileId, id)
			if err != nil {
				return err
			}
		}
	}

	for id := range existingMap {
		if _, exists := newMap[id]; !exists {
			_, err := tx.Exec(`DELETE FROM profiles_profile_douyin WHERE profile_id = ? AND profile_douyin_id = ?`, profileId, id)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (am *ProfileManager) UpdateAuthenticatedStatus(id int, isAuthenticated bool) error {
	updateSQL := `UPDATE profiles SET is_authenticated=? WHERE id=?`
	_, err := am.db.Exec(updateSQL, isAuthenticated, id)
	return err
}
