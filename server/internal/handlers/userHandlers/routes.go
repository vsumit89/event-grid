package userHandlers

import (
	"net/http"
	"server/internal/commons"
	"server/internal/handlers/middlewares"
	"server/internal/services"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	userSvc services.IUserSvc
	jwtSvc  *commons.JwtSvc
}

type AuthHandler struct {
	userSvc services.IUserSvc
	jwtSvc  *commons.JwtSvc
}

func NewHandler(userSvc services.IUserSvc, jwt *commons.JwtSvc) *Handler {
	return &Handler{
		userSvc: userSvc,
		jwtSvc:  jwt,
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

	userRouter.Use(middlewares.JWTAuth(h.jwtSvc))

	userRouter.Get("/profile", h.getUser)

	userRouter.Get("/search", h.searchUsers)

	return userRouter
}

func (h *AuthHandler) AuthRoutes() http.Handler {
	authRouter := chi.NewRouter()

	authRouter.Post("/register", h.registerUser)

	authRouter.Post("/login", h.loginUser)

	return authRouter
}
