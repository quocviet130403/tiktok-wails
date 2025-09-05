package initialize

func InitServer() error {

	dbInit, err := InitDatabase()
	if err != nil {
		panic("Error initializing database: " + err.Error())
	}
	err = InitGlobal(dbInit)
	if err != nil {
		panic("Error initializing global: " + err.Error())
	}

	InitManage(dbInit)
	InitSchedule()

	return nil
}
