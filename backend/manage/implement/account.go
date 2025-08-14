package implement

import (
	"database/sql"
	"fmt"
	"tiktok-wails/backend/manage/service"
)

type AccountManager struct {
	db *sql.DB
}

func NewAccountManager(db *sql.DB) *AccountManager {
	return &AccountManager{db: db}
}

func (am *AccountManager) UpdateAccount(id int, name, urlReup, hashtag, firstComment string) error {
	updateSQL := `UPDATE accounts SET name=?, url_reup=?, hashtag=?, first_comment=? WHERE id=?`
	_, err := am.db.Exec(updateSQL, name, urlReup, hashtag, firstComment, id)
	return err
}

// Hàm xóa tài khoản
func (am *AccountManager) DeleteAccount(id int) error {
	deleteSQL := `DELETE FROM accounts WHERE id=?`
	_, err := am.db.Exec(deleteSQL, id)
	return err
}

func (am *AccountManager) AddAccount(name, urlReup, hashtag, firstComment string) error {
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
	}

	// Thực hiện insert trong transaction
	insertSQL := `INSERT INTO accounts (name, url_reup, hashtag, first_comment, last_video_reup, retry_count, is_authenticated) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(insertSQL, name, urlReup, hashtag, firstComment, "", 0, isAuthenticated)
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

func (am *AccountManager) GetAllAccounts() ([]service.Accounts, error) {
	rows, err := am.db.Query("SELECT id, name, url_reup, hashtag, first_comment, last_video_reup, retry_count, is_authenticated FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []service.Accounts

	for rows.Next() {
		var account service.Accounts
		err := rows.Scan(&account.ID, &account.Name, &account.UrlReup, &account.Hashtag, &account.FirstComment, &account.LastVideoReup, &account.RetryCount, &account.IsAuthenticated)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}
