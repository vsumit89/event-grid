package eventHandlers

import (
	"net/http"
	"server/internal/commons"
	"server/internal/handlers/middlewares"
	"server/internal/services"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	eventSvc services.IEventSvc
	jwtSvc   *commons.JwtSvc
}

func NewEventHandler(eventSvc services.IEventSvc, jwtSvc *commons.JwtSvc) *Handler {
	return &Handler{
		eventSvc: eventSvc,
		jwtSvc:   jwtSvc,
	}
}

func (h *Handler) GetRoutes() http.Handler {
	eventRouter := chi.NewRouter()

	eventRouter.Use(middlewares.JWTAuth(h.jwtSvc))

	eventRouter.Post("/", h.createEvent)

	return eventRouter
}
