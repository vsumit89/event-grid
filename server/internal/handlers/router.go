package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetRoutes() http.Handler {
	router := chi.NewRouter()

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return router
}
