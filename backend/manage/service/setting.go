package service

import "log"

type SettingInterface interface {
	GetAllSettings() (map[string]string, error)
	GetSetting(key string) (string, error)
	SetSetting(key, value string) error
}

var (
	localSetting SettingInterface
)

func SettingManager() SettingInterface {
	if localSetting == nil {
		log.Println("[Service] SettingManager chưa khởi tạo")
		return nil
	}
	return localSetting
}

func InitSettingManager(i SettingInterface) {
	localSetting = i
}
