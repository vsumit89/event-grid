package handlers

import (
	"net/http"
	"server/pkg/logger"

	"github.com/go-chi/chi/v5"
)

func GetRoutes() http.Handler {
	router := chi.NewRouter()

	logger.Info("setting up routes")
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return router
}
