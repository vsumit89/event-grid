package userHandlers

import (
	"net/http"
	"server/internal/services"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	userSvc services.IUserSvc
}

func NewHandler(userSvc services.IUserSvc) *Handler {
	return &Handler{
		userSvc: userSvc,
	}
}

func (h *Handler) GetRoutes() http.Handler {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	return router
}
