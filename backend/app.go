package backend

import (
	"context"
	"database/sql"
)

// App struct
type App struct {
	db  *sql.DB
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp(db *sql.DB) *App {
	return &App{
		db: db,
	}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// GetAllAccounts returns a list of all accounts
func (a *App) GetAllAccounts() []string {
	return []string{"Account1", "Account2", "Account3"}
}
