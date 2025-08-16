package service

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
		panic("SettingManager is not initialized")
	}
	return localSetting
}

func InitSettingManager(i SettingInterface) {
	localSetting = i
}
