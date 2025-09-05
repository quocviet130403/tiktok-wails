package service

type PythonInterface interface {
	TranslateVideo(id int) error
}

var (
	localPython PythonInterface
)

func PythonManager() PythonInterface {
	if localPython == nil {
		panic("PythonManager is not initialized")
	}
	return localPython
}

func InitPythonManager(i PythonInterface) {
	localPython = i
}
