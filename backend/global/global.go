package global

import (
	"database/sql"
)

var (
	DB                *sql.DB
	PathTempProfile   string
	PathVideoReup     string
	PathHandleCaptcha string
	PathAppChrome     string
)
