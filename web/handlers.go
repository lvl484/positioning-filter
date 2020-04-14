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

func (h *handler) AddFilter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDstring := vars[inputUserID]

	userid, err := uuid.Parse(userIDstring)

	if err != nil {
		h.log.Errorf("Can't parse UserID from URL query: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var filter repository.Filter

	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		h.log.Errorf("Can't decode request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filter.UserID = userid

	if err := h.filters.Add(&filter); err != nil {
		h.log.Errorf("Can't add filter to database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.log.Infof("Added filter %v for user %v", filter.Name, filter.UserID)
	w.WriteHeader(http.StatusCreated)
}

func (h *handler) GetOneFilter(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	filterName := m[inputName]
	userIDstring := m[inputUserID]
	userID, err := uuid.Parse(userIDstring)

	if err != nil {
		h.log.Errorf("Can't parse UserID from URL query: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	f, err := h.filters.OneByUser(userID, filterName)
	if err != nil {
		h.log.Errorf("Can't get filter from database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(f); err != nil {
		h.log.Errorf("Can't encode body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	h.log.Infof("Got one filter %v for user %v", filterName, userID)
	w.WriteHeader(http.StatusAccepted)
}

func (h *handler) GetOffset(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	userIDstring := m[inputUserID]
	userID, err := uuid.Parse(userIDstring)
	if err != nil {
		h.log.Errorf("Can't parse UserID from URL query: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	offsetString := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetString)
	if err != nil {
		h.log.Errorf("Can't convert offset to integer: %v", err)
	}

	filters, err := h.filters.OffsetByUser(userID, offset)
	if err != nil {
		h.log.Errorf("Can't get filters from database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(filters); err != nil {
		h.log.Errorf("Can't encode body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.log.Infof("Got filters for user %v", userID)
	w.WriteHeader(http.StatusOK)
}

func (h *handler) UpdateFilter(w http.ResponseWriter, r *http.Request) {
	var filter repository.Filter

	vars := mux.Vars(r)

	filter.Name = vars[inputName]
	userIDstring := vars[inputUserID]
	userUUID, err := uuid.Parse(userIDstring)

	if err != nil {
		h.log.Errorf("Can't parse UserID from URL query: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	filter.UserID = userUUID

	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		h.log.Errorf("Can't decode request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.filters.Update(&filter); err != nil {
		h.log.Errorf("Can't update filter in database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.log.Infof("Updated filter %v for user %v", filter.Name, filter.UserID)
	w.WriteHeader(http.StatusOK)
}

func (h *handler) DeleteFilter(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	filterName := m[inputName]
	userIDstring := m[inputUserID]
	userID, err := uuid.Parse(userIDstring)
	if err != nil {
		h.log.Errorf("Can't parse UserID from URL query: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	filters := h.filters.Delete(userID, filterName)
	if err != nil {
		h.log.Errorf("Can't delete filter in database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(filters); err != nil {
		h.log.Errorf("Can't encode body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.log.Infof("Deleted filter %v for user %v", filterName, userID)
	w.WriteHeader(http.StatusNoContent)
}
