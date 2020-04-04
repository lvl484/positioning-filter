package repository

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/lvl484/positioning-filter/config"
	"github.com/lvl484/positioning-filter/storage"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	db := newTestDB(t)
	tx, err := db.Begin()
	assert.NoError(t, err)

	filter := addTestFilterToDB(tx, t)
	assert.NotNil(t, filter)

	defer func(*testing.T) {
		err := tx.Rollback()
		assert.NoError(t, err)
	}(t)
}

func TestOneByUser(t *testing.T) {
	db := newTestDB(t)
	tx, err := db.Begin()
	assert.NoError(t, err)

	filter := addTestFilterToDB(tx, t)

	row := tx.QueryRow(getOneQuery, filter.UserID, filter.Name)

	var filterFromDB Filter
	err = row.Scan(
		&filterFromDB.Name,
		&filterFromDB.Type,
		&filterFromDB.Configuration,
		&filterFromDB.Reversed,
		&filterFromDB.UserID,
	)
	assert.NoError(t, err)
	assert.Equal(t, *filter, filterFromDB)

	defer func(*testing.T) {
		err := tx.Rollback()
		assert.NoError(t, err)
	}(t)
}

func TestAllByUser(t *testing.T) {
	db := newTestDB(t)
	tx, err := db.Begin()
	assert.NoError(t, err)

	filter := addTestFilterToDB(tx, t)

	rows, err := tx.Query(getOneQuery, filter.UserID, filter.Name)
	assert.NoError(t, err)

	var filtersFromDB []Filter

	for rows.Next() {
		var filterFromDB Filter
		err = rows.Scan(
			&filterFromDB.Name,
			&filterFromDB.Type,
			&filterFromDB.Configuration,
			&filterFromDB.Reversed,
			&filterFromDB.UserID,
		)
		assert.NoError(t, err)

		filtersFromDB = append(filtersFromDB, filterFromDB)

		err = rows.Err()
		assert.NoError(t, err)
	}

	assert.NotNil(t, filtersFromDB)

	defer func(*testing.T) {
		err := tx.Rollback()
		assert.NoError(t, err)
	}(t)
}

func TestUpdate(t *testing.T) {
	db := newTestDB(t)
	tx, err := db.Begin()
	assert.NoError(t, err)

	filter := addTestFilterToDB(tx, t)

	confUpdate, err := json.Marshal("someOtherString")
	assert.NoError(t, err)

	filter.Type = "rectangular"
	filter.Configuration = confUpdate
	filter.Reversed = true

	_, err = tx.Exec(updateQuery, filter.Type, filter.Configuration, filter.Reversed, filter.UserID, filter.Name)
	assert.NoError(t, err)

	defer func(*testing.T) {
		err := tx.Rollback()
		assert.NoError(t, err)
	}(t)
}

func TestDelete(t *testing.T) {
	db := newTestDB(t)
	tx, err := db.Begin()
	assert.NoError(t, err)

	filter := addTestFilterToDB(tx, t)

	_, err = tx.Exec(deleteQuery, filter.UserID, filter.Name)
	assert.NoError(t, err)

	defer func(*testing.T) {
		err := tx.Rollback()
		assert.NoError(t, err)
	}(t)
}

func newTestDB(t *testing.T) *sql.DB {
	viper, err := config.NewConfig("viper.config", "../config/")
	assert.NoError(t, err)

	dbConfig := viper.NewDBConfig()
	dbConfig.Host = "localhost"
	db, err := storage.Connect(dbConfig)
	assert.NoError(t, err)

	return db
}

func addTestFilterToDB(tx *sql.Tx, t *testing.T) *Filter {
	conf, err := json.Marshal("someString")
	assert.NoError(t, err)

	filter := Filter{
		Name:          "Filter",
		Type:          "round",
		Configuration: conf,
		Reversed:      false,
		UserID:        uuid.New(),
	}

	_, err = tx.Exec(addQuery, filter.Name, filter.Type, filter.Configuration, filter.Reversed, filter.UserID)
	assert.NoError(t, err)

	return &filter
}
