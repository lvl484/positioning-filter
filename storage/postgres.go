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

func Connect(psqlHost, psqlPort, psqlUser, psqlPass, psqlDB string) (*sql.DB, error) {
	connStr := fmt.Sprintf(connFmt, psqlHost, psqlPort, psqlUser, psqlPass, psqlDB)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
