// Package web provides managing HTTP protocol.
// Implement CRUD operations for filters over HTTP API.
package web

import (
	"context"
	"net/http"

	"github.com/lvl484/positioning-filter/repository"
	"github.com/sirupsen/logrus"
)

type WebServer struct {
	server *http.Server
}

func NewWebServer(filters repository.Filters, addr string, log *logrus.Logger) *WebServer {
	router := newRouter(filters, log)
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
