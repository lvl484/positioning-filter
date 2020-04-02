// Package web provides managing HTTP protocol.
// Implement CRUD operations for filters over HTTP API.
package web

import (
	"context"
	"net/http"

	"github.com/lvl484/positioning-filter/repository"
)

type WebServer struct {
	server *http.Server
}

func NewWebServer(filters repository.Filters, port string) *WebServer {
	router := NewRouter(filters)
	return &WebServer{
		&http.Server{
			Addr:    port,
			Handler: router,
		},
	}
}

func (ws *WebServer) Close() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return ws.server.Shutdown(ctx)
}

func (ws *WebServer) Run() error {
	return ws.server.ListenAndServe()
}
