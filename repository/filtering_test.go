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
	repo := NewFiltersRepository(db)
	filter1 := newTestFilter("Name1", "round")
	filter2 := newTestFilter("Name2", "round")
	filter3 := newTestFilter("Name2", "round")
	err := repo.Add(filter1)
	assert.NoError(t, err)
	err = repo.Add(filter2)
	assert.NoError(t, err)
	err = repo.Add(filter3)
	assert.NotNil(t, err)
}

func TestOneByUser(t *testing.T) {
	db := newTestDB(t)
	repo := NewFiltersRepository(db)
	userID, err := uuid.Parse("d5cadefb-4d4d-4105-8244-1c354f936e69")
	assert.NoError(t, err)
	filter, err := repo.OneByUser(userID, "Name2")
	assert.NoError(t, err)
	assert.NotNil(t, filter)
}

func TestAllByUser(t *testing.T) {
	db := newTestDB(t)
	repo := NewFiltersRepository(db)
	userID, err := uuid.Parse("d5cadefb-4d4d-4105-8244-1c354f936e69")
	assert.NoError(t, err)
	gotFilters, err := repo.AllByUser(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, gotFilters)
}

func TestUpdate(t *testing.T) {
	db := newTestDB(t)
	repo := NewFiltersRepository(db)
	filter := newTestFilter("Name1", "rectangular")
	err := repo.Update(filter)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	db := newTestDB(t)
	repo := NewFiltersRepository(db)
	userID, err := uuid.Parse("d5cadefb-4d4d-4105-8244-1c354f936e69")
	assert.NoError(t, err)
	err = repo.Delete(userID, "Name1")
	assert.NoError(t, err)
	err = repo.Delete(userID, "Name2")
	assert.NoError(t, err)
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

func newTestFilter(filterName, filterType string) *Filter {
	conf, _ := json.Marshal("someString")
	userID, _ := uuid.Parse("d5cadefb-4d4d-4105-8244-1c354f936e69")

	return &Filter{
		Name:          filterName,
		Type:          filterType,
		Configuration: conf,
		Reversed:      false,
		UserID:        userID,
	}
}
