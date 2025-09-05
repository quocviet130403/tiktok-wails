package implement

import "database/sql"

type PythonManager struct {
	db *sql.DB
}

func NewPythonManager(db *sql.DB) *PythonManager {
	return &PythonManager{
		db: db,
	}
}

func (pm *PythonManager) TranslateVideo(id int) error {
	return nil
}
