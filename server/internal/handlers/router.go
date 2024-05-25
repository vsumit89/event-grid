package handlers

import (
	"net/http"
	"server/internal/handlers/userHandlers"
	"server/internal/services"

	"github.com/go-chi/chi/v5"
)

// Container is a struct that contains all the services dependencies for the router
type Container struct {
	UserSvc services.IUserSvc
}

func GetRoutes(c *Container) http.Handler {
	router := chi.NewRouter()

	v1Router := chi.NewRouter()

	router.Mount("/api/v1", v1Router)

	userHandlers.NewHandler(c.UserSvc).GetRoutes(v1Router)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return router
}
