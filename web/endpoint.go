// Package web provides managing HTTP protocol.
// Implement CRUD operations for filters over HTTP API.
package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lvl484/positioning-filter/repository"
	"github.com/sirupsen/logrus"
)

func newRouter(filters repository.Filters, log *logrus.Logger) *mux.Router {
	handle := newHandler(filters, log)

	router := mux.NewRouter()
	filtersRouter := router.PathPrefix("/users/{user_id}/filters").Subrouter()

	router.HandleFunc("/health", handle.Health).Methods(http.MethodGet)
	filtersRouter.HandleFunc("/", handle.AddFilter).Methods(http.MethodPost)
	filtersRouter.HandleFunc("/", handle.GetOffset).Methods(http.MethodGet)
	filtersRouter.HandleFunc("/{name}", handle.UpdateFilter).Methods(http.MethodPatch)
	filtersRouter.HandleFunc("/{name}", handle.GetOneFilter).Methods(http.MethodGet)
	filtersRouter.HandleFunc("/{name}", handle.DeleteFilter).Methods(http.MethodDelete)

	return router
}
