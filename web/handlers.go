// Package web provides managing HTTP protocol.
// Implement CRUD operations for filters over HTTP API.
package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lvl484/positioning-filter/repository"
)

const (
	inputUserID = "user_id"
	inputName   = "name"
)

type repo struct {
	filters repository.Filters
}

func newRepo(filters repository.Filters) *repo {
	return &repo{
		filters: filters,
	}
}

func (repo *repo) AddFilter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDstring := vars[inputUserID]

	userid, err := uuid.Parse(userIDstring)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var filter repository.Filter

	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filter.UserID = userid

	if err := repo.filters.Add(&filter); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (repo *repo) GetOneFilter(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	filterName := m[inputName]
	userIDstring := m[inputUserID]
	userID, err := uuid.Parse(userIDstring)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	f, err := repo.filters.OneByUser(userID, filterName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(f); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusAccepted)
}

func (repo *repo) GetOffset(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	userIDstring := m[inputUserID]
	userID, err := uuid.Parse(userIDstring)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	offsetString := r.URL.Query().Get("offset")
	offset, _ := strconv.Atoi(offsetString)

	filters, err := repo.filters.OffsetByUser(userID, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(filters); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (repo *repo) UpdateFilter(w http.ResponseWriter, r *http.Request) {
	var filter repository.Filter

	vars := mux.Vars(r)

	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	filter.Name = vars[inputName]
	userIDstring := vars[inputUserID]
	userUUID, err := uuid.Parse(userIDstring)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	filter.UserID = userUUID

	if err := repo.filters.Update(&filter); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (repo *repo) DeleteFilter(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	filterName := m[inputName]
	userIDstring := m[inputUserID]
	userID, err := uuid.Parse(userIDstring)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	filters := repo.filters.Delete(userID, filterName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(filters); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
