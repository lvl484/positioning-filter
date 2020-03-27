// Package filter provides filter model
package filter

import (
	"database/sql"

	"github.com/google/uuid"
)

type Postgres struct {
	db *sql.DB
}

const (
	addQuery    = "INSERT INTO FILTERS(name,type,configutation,reversed,user_id) VALUES ($1,$2,$3,$4,$5)"
	getQuery    = "SELECT * FROM FILTERS WHERE user_id=$1"
	updateQuery = "UPDATE FILTERS SET (type,configutation,reversed) = ($1,$2,$3) WHERE name = $4"
	deleteQuery = "DELETE FROM FILTERS WHERE name=$1"
)

// NewPostgres returns Postgres with db
func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

// AddFilter adds new filter to database
func (p *Postgres) AddFilter(filter *Filter) error {
	_, err := p.db.Exec(addQuery, filter.Name, filter.Type, filter.Configuration, filter.Reversed, filter.UserID)
	return err
}

// GetFilters returns slice of filters relevant to user
func (p *Postgres) GetFilters(userID uuid.UUID) ([]Filter, error) {
	filterSlice := []Filter{}
	r, err := p.db.Query(getQuery, userID)

	if err != nil {
		return nil, err
	}

	for r.Next() {
		f := Filter{}
		if err := r.Scan(&f); err != nil {
			continue
		}

		if err := r.Err(); err != nil {
			continue
		}

		filterSlice = append(filterSlice, f)
	}

	return filterSlice, nil
}

// UpdateFilter updates filter fields by filter name
func (p *Postgres) UpdateFilter(filter *Filter) error {
	_, err := p.db.Exec(updateQuery, filter.Type, filter.Configuration, filter.Reversed, filter.Name)

	return err
}

// DeleteFilter deletes filter by name
func (p *Postgres) DeleteFilter(filterName string) error {
	_, err := p.db.Exec(deleteQuery, filterName)
	return err
}