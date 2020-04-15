// Package web provides managing HTTP protocol.
// Implement CRUD operations for filters over HTTP API.
package web

import (
	"context"
	"net/http"

	"github.com/lvl484/positioning-filter/repository"
	"github.com/sirupsen/logrus"
)

type Server struct {
	server *http.Server
}

func NewServer(filters repository.Filters, addr string, log *logrus.Logger) *Server {
	router := newRouter(filters, log)

	return &Server{
		&http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}

func (ws *Server) Close() error {
	ctx := context.Background()
	return ws.server.Shutdown(ctx)
}

func (ws *Server) Run() error {
	return ws.server.ListenAndServe()
}
