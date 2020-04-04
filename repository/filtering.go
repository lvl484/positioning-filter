// Package repository provides filter model
package repository

import (
	"database/sql"

	"github.com/google/uuid"
)

const (
	addQuery    = "INSERT INTO FILTERS(name,type,configutation,reversed,user_id) VALUES ($1,$2,$3,$4,$5)"
	getOneQuery = "SELECT name, type, configuration, reversed, user_id FROM FILTERS WHERE user_id=$1 AND name=$2"
	getAllQuery = "SELECT name, type, configuration, reversed, user_id  FROM FILTERS WHERE user_id=$1"
	updateQuery = "UPDATE FILTERS SET (type,configutation,reversed) = ($1,$2,$3) WHERE user_id=$4 AND name=$5"
	deleteQuery = "DELETE FROM FILTERS WHERE user_id=$1 AND name=$2"
)

type Filters interface {
	Add(filter *Filter) error
	OneByUser(userID uuid.UUID, filterName string) (*Filter, error)
	AllByUser(userID uuid.UUID) ([]*Filter, error)
	Update(filter *Filter) error
	Delete(userID uuid.UUID, filterName string) error
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

// OneByUser returns filter for relevant user
func (p *filtersRepo) OneByUser(userID uuid.UUID, filterName string) (*Filter, error) {
	filter := new(Filter)

	row := p.db.QueryRow(getOneQuery, userID, filterName)

	err := row.Scan(&filter.Name, &filter.Type, &filter.Configuration, &filter.Reversed, &filter.UserID)
	if err != nil {
		return nil, err
	}

	return filter, nil
}

// AllByUser returns set of filters for relevant user
func (p *filtersRepo) AllByUser(userID uuid.UUID) ([]*Filter, error) {
	filters := []*Filter{}

	rows, err := p.db.Query(getAllQuery, userID)
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

// Update updates filter fields by filter id for relevant user
func (p *filtersRepo) Update(filter *Filter) error {
	_, err := p.db.Exec(updateQuery, filter.Type, filter.Configuration, filter.Reversed, filter.UserID, filter.Name)
	return err
}

// Delete deletes filter by id for relevant user
func (p *filtersRepo) Delete(userID uuid.UUID, filterName string) error {
	_, err := p.db.Exec(deleteQuery, userID, filterName)
	return err
}
