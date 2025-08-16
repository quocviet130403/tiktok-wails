package implement

import "database/sql"

type SettingManager struct {
	db *sql.DB
}

func NewSettingManager(db *sql.DB) *SettingManager {
	return &SettingManager{
		db: db,
	}
}

func (sm *SettingManager) GetAllSettings() (map[string]string, error) {
	rows, err := sm.db.Query("SELECT key, value FROM settings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		settings[key] = value
	}
	return settings, nil
}

func (sm *SettingManager) GetSetting(key string) (string, error) {
	var value string
	err := sm.db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (sm *SettingManager) SetSetting(key, value string) error {
	_, err := sm.db.Exec("INSERT INTO settings (key, value) VALUES (?, ?) ON DUPLICATE KEY UPDATE value = ?", key, value, value)
	return err
}
