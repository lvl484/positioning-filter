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
	"github.com/sirupsen/logrus"
)

const (
	inputUserID = "user_id"
	inputName   = "name"
)

type handler struct {
	filters repository.Filters
	log     *logrus.Logger
}

func newHandler(filters repository.Filters, log *logrus.Logger) *handler {
	return &handler{
		filters: filters,
		log:     log,
	}
}

func (handler *handler) AddFilter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDstring := vars[inputUserID]

	userid, err := uuid.Parse(userIDstring)
	if err != nil {
		handler.log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var filter repository.Filter

	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		handler.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filter.UserID = userid

	if err := handler.filters.Add(&filter); err != nil {
		handler.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *handler) GetOneFilter(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	filterName := m[inputName]
	userIDstring := m[inputUserID]
	userID, err := uuid.Parse(userIDstring)

	if err != nil {
		handler.log.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	f, err := handler.filters.OneByUser(userID, filterName)
	if err != nil {
		handler.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(f); err != nil {
		handler.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusAccepted)
}

func (handler *handler) GetOffset(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	userIDstring := m[inputUserID]
	userID, err := uuid.Parse(userIDstring)
	if err != nil {
		handler.log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	offsetString := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetString)
	if err != nil {
		handler.log.Error(err)
	}

	filters, err := handler.filters.OffsetByUser(userID, offset)
	if err != nil {
		handler.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(filters); err != nil {
		handler.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *handler) UpdateFilter(w http.ResponseWriter, r *http.Request) {
	var filter repository.Filter

	vars := mux.Vars(r)

	filter.Name = vars[inputName]
	userIDstring := vars[inputUserID]
	userUUID, err := uuid.Parse(userIDstring)

	if err != nil {
		handler.log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	filter.UserID = userUUID

	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		handler.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := handler.filters.Update(&filter); err != nil {
		handler.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *handler) DeleteFilter(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	filterName := m[inputName]
	userIDstring := m[inputUserID]
	userID, err := uuid.Parse(userIDstring)
	if err != nil {
		handler.log.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	filters := handler.filters.Delete(userID, filterName)
	if err != nil {
		handler.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(filters); err != nil {
		handler.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
