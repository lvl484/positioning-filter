// Package web provides managing HTTP protocol.
// Implement CRUD operations for filters over HTTP API.
package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lvl484/positioning-filter/repository"
)

func InitServer(filters repository.Filters, port string) http.Server {
	handle := NewWebFilters(filters)

	router := mux.NewRouter()
	router.Use(Middleware)

	router.HandleFunc("/", handle.AddFilter).Methods(http.MethodPost)
	router.HandleFunc("/{USERID}/{FILTERNAME}", handle.UpdateFilter).Methods(http.MethodPut)
	router.HandleFunc("/{USERID}/{FILTERNAME}", handle.DeleteFilter).Methods(http.MethodDelete)
	router.HandleFunc("/{USERID}", handle.GetFiltersByUser).Methods(http.MethodGet)

	return http.Server{
		Addr:    port,
		Handler: router,
	}
}

/*
	filtersRouter := router.PathPrefix("/router/{user_id}/filters").Subrouter()

	filtersRouter.HandleFunc("/", handle.AddFilter).Methods(http.MethodPost)
	filtersRouter.HandleFunc("/", handle.UpdateFilter).Methods(http.MethodPut)
	filtersRouter.HandleFunc("/", handle.DeleteFilter).Methods(http.MethodDelete)
	filtersRouter.HandleFunc("/", handle.GetFiltersByUser).Methods(http.MethodGet)
*/
