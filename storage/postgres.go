// Package storage provides postgres database configuration for storing messages.
package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/lvl484/positioning-filter/config"
)

const (
	connFmt = "host=%v port=%v user=%v password=%v dbname=%v sslmode=disable"
)

// Connect connects to DB with args and returns pointer to DB instance
func Connect(p *config.PostgresConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf(connFmt, p.Host, p.Port, p.User, p.Pass, p.DB)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
