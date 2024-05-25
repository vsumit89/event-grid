package transport

import (
	"context"
	"fmt"
	"net/http"
	"server/internal/config"
)

type httpServer struct {
	server *http.Server
}

func NewHTTPServer(cfg *config.ServerConfig, handler http.Handler) *httpServer {
	return &httpServer{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%s", cfg.Port),
			Handler: handler,
		},
	}
}

func (h *httpServer) Run() error {
	err := h.server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (h *httpServer) Stop() error {
	return h.server.Shutdown(context.Background())
}
