// Package web provides managing HTTP protocol.
// Implement CRUD operations for filters over HTTP API.
package web

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lvl484/positioning-filter/repository"
)

type WebFilters struct {
	filters repository.Filters
}

func NewWebFilters(filters repository.Filters) *WebFilters {
	return &WebFilters{
		filters: filters,
	}
}

func (wb *WebFilters) AddFilter(rw http.ResponseWriter, r *http.Request) {
	var filter repository.Filter

	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := wb.filters.Add(&filter); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}

}

func (wb *WebFilters) GetFiltersByUser(rw http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	userIDstring := m["USERID"]
	userID, err := uuid.Parse(userIDstring)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
	}

	filters, err := wb.filters.AllByUser(userID)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}

	if err := json.NewEncoder(rw).Encode(filters); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (wb *WebFilters) UpdateFilter(rw http.ResponseWriter, r *http.Request) {
	var filter repository.Filter

	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	m := mux.Vars(r)
	filterName := m["FILTERNAME"]
	userIDstring := m["USERID"]
	userID, err := uuid.Parse(userIDstring)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	if err := wb.filters.Update(userID, filterName, &filter); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent) //swagger
}

func (wb *WebFilters) DeleteFilter(rw http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	filterName := m["FILTERNAME"]
	n := mux.Vars(r)
	userIDstring := n["USERID"]
	userID, err := uuid.Parse(userIDstring)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
	}

	filters := wb.filters.Delete(userID, filterName) //!!!???
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}

	if err := json.NewEncoder(rw).Encode(filters); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
