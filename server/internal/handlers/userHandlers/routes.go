package userHandlers

import (
	"fmt"
	"net/http"
	"server/internal/commons"
	"server/internal/services"
	"server/pkg/logger"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	userSvc services.IUserSvc
}

type AuthHandler struct {
	userSvc services.IUserSvc
	jwtSvc  *commons.JwtSvc
}

func NewHandler(userSvc services.IUserSvc) *Handler {
	return &Handler{
		userSvc: userSvc,
	}
}

func NewAuthHandler(userSvc services.IUserSvc, jwtSvc *commons.JwtSvc) *AuthHandler {
	return &AuthHandler{
		userSvc: userSvc,
		jwtSvc:  jwtSvc,
	}
}

func (h *Handler) GetRoutes() http.Handler {
	userRouter := chi.NewRouter()

	userRouter.Get("/", h.getUsers)

	return userRouter
}

func (h *AuthHandler) AuthRoutes() http.Handler {
	authRouter := chi.NewRouter()

	authRouter.Post("/register", h.registerUser)

	authRouter.Post("/login", h.loginUser)

	return authRouter
}

func (h *Handler) getUsers(w http.ResponseWriter, r *http.Request) {

	logger.Info("retrieving users")
	fmt.Println(h.userSvc)
	// TODO: get users from the user service

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully retrieved users"))
}
