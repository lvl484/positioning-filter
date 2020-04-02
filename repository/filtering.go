// Package repository provides filter model
package repository

import (
	"database/sql"

	"github.com/google/uuid"
)

const (
	addQuery    = "INSERT INTO FILTERS(id,type,configutation,reversed,user_id) VALUES ($1,$2,$3,$4,$5)"
	getAllQuery = "SELECT * FROM FILTERS WHERE user_id=$1"
	getOneQuery = "SELECT * FROM FILTERS WHERE user_id=$1 AND id=$2"
	updateQuery = "UPDATE FILTERS SET (type,configutation,reversed) = ($1,$2,$3) WHERE user_id=$4 AND id=$5"
	deleteQuery = "DELETE FROM FILTERS WHERE user_id=$1 AND id=$2"
)

type Filters interface {
	OneByUser(userID uuid.UUID, ID uuid.UUID) (*Filter, error)
	AllByUser(userID uuid.UUID) ([]*Filter, error)
	Add(filter *Filter) error
	Update(filter *Filter) error
	Delete(userID uuid.UUID, ID uuid.UUID) error
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
	filterID, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	filter.ID = filterID

	_, err = p.db.Exec(addQuery, filter.ID, filter.Type, filter.Configuration, filter.Reversed, filter.UserID)
	return err
}

// OneByUser returns filter for relevant user
func (p *filtersRepo) OneByUser(userID uuid.UUID, ID uuid.UUID) (*Filter, error) {
	var filter *Filter

	row := p.db.QueryRow(getOneQuery, userID)

	err := row.Scan(filter.ID, filter.Type, filter.Configuration, filter.Reversed, filter.UserID)
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
		err = rows.Scan(&f.ID, &f.Type, &f.Configuration, &f.Reversed, &f.UserID)
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
	_, err := p.db.Exec(updateQuery, filter.Type, filter.Configuration, filter.Reversed, filter.UserID, filter.ID)
	return err
}

// Delete deletes filter by id for relevant user
func (p *filtersRepo) Delete(userID uuid.UUID, ID uuid.UUID) error {
	_, err := p.db.Exec(deleteQuery, userID, ID)
	return err
}
