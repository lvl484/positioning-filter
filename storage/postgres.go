// Package storage provides postgres database configuration for storing messages.
package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	connFmt = "host=%v port=%v user=%v password=%v dbname=%v sslmode=disable"
)

// StorageConfig contains fileds used in Connect for DSN
type DBConfig struct {
	Host string
	Port string
	User string
	Pass string
	DB   string
}

// Connect connects to DB with args and returns pointer to DB instance
func Connect(p *DBConfig) (*sql.DB, error) {
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
