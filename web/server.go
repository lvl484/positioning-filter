// Package web provides managing HTTP protocol.
// Implement CRUD operations for filters over HTTP API.
package web

import (
	"net/http"

	"github.com/lvl484/positioning-filter/repository"
)

func InitServer(filters repository.Filters, port string) http.Server {
	router := NewRouter(filters)
	return http.Server{
		Addr:    port,
		Handler: router,
	}
}
