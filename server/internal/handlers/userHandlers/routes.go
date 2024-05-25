package userHandlers

import (
	"fmt"
	"net/http"
	"server/internal/services"
	"server/pkg/logger"

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

func (h *Handler) GetRoutes(router *chi.Mux) {
	router.Get("/users", h.GetUsers)
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {

	logger.Info("retrieving users")
	fmt.Println(h.userSvc)
	// TODO: get users from the user service

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully retrieved users"))
}
