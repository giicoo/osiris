package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/giicoo/osiris/notification-service/internal/config"
)

type Server struct {
	cfg        *config.Config
	httpServer *http.Server
}

func NewServer(cfg *config.Config, r http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
			Handler: r,
		},
	}
}

func (srv *Server) StartServer() error {
	return srv.httpServer.ListenAndServe()
}

func (srv *Server) ShutdownServer(ctx context.Context) error {
	return srv.httpServer.Shutdown(ctx)
}
