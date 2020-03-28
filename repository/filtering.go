// Package repository provides filter model
package repository

import (
	"database/sql"

	"github.com/google/uuid"
)

const (
	addQuery    = "INSERT INTO FILTERS(name,type,configutation,reversed,user_id) VALUES ($1,$2,$3,$4,$5)"
	getQuery    = "SELECT * FROM FILTERS WHERE user_id=$1"
	updateQuery = "UPDATE FILTERS SET (type,configutation,reversed) = ($1,$2,$3) WHERE name = $4"
	deleteQuery = "DELETE FROM FILTERS WHERE user_id=$1 AND name=$2"
)

type Filters interface {
	AllByUser(userID uuid.UUID) ([]*Filter, error)
	Add(filter *Filter) error
	Update(filter *Filter) error
	Delete(userID uuid.UUID, name string) error
}

type filtersRepo struct {
	db *sql.DB
}

func NewFiltersRepository(db *sql.DB) Filters {
	return &filtersRepo{
		db: db,
	}
}

// Add adds new filter to database
func (p *filtersRepo) Add(filter *Filter) error {
	_, err := p.db.Exec(addQuery, filter.Name, filter.Type, filter.Configuration, filter.Reversed, filter.UserID)
	return err
}

// AllByUser returns set of filters for relevant user
func (p *filtersRepo) AllByUser(userID uuid.UUID) ([]*Filter, error) {
	filters := []*Filter{}

	rows, err := p.db.Query(getQuery, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		f := new(Filter)
		err = rows.Scan(&f.Name, &f.Type, &f.Configuration, &f.Reversed, &f.UserID)
		if err != nil {
			return nil, err
		}

		filters = append(filters, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return filters, nil
}

// Update updates filter fields by filter name
func (p *filtersRepo) Update(filter *Filter) error {
	_, err := p.db.Exec(updateQuery, filter.Type, filter.Configuration, filter.Reversed, filter.Name)
	return err
}

// Delete deletes filter by name for relevant user
func (p *filtersRepo) Delete(userID uuid.UUID, filterName string) error {
	_, err := p.db.Exec(deleteQuery, userID, filterName)
	return err
}
