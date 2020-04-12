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

func NewWebServer(filters repository.Filters, addr string) *WebServer {
	router := newRouter(filters)
	return &WebServer{
		&http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}

func (ws *WebServer) Close() error {
	ctx := context.Background()
	return ws.server.Shutdown(ctx)
}

func (ws *WebServer) Run() error {
	return ws.server.ListenAndServe()
}
