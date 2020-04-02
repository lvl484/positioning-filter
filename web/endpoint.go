// Package web provides managing HTTP protocol.
// Implement CRUD operations for filters over HTTP API.
package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lvl484/positioning-filter/repository"
)

func NewRouter(filters repository.Filters) *mux.Router {
	handle := NewWebFilters(filters)

	router := mux.NewRouter()
	router.Use(Middleware)

	filtersRouter := router.PathPrefix("/router/{user_id}/filters").Subrouter()

	filtersRouter.HandleFunc("/", handle.AddFilter).Methods(http.MethodPost)
	filtersRouter.HandleFunc("/", handle.GetFiltersByUser).Methods(http.MethodGet)
	filtersRouter.HandleFunc("/{id}", handle.UpdateFilter).Methods(http.MethodPut)
	filtersRouter.HandleFunc("/{id}", handle.DeleteFilter).Methods(http.MethodDelete)
	return router
}
