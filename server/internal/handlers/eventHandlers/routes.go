package eventHandlers

import (
	"net/http"
	"server/internal/commons"
	"server/internal/handlers/middlewares"
	"server/internal/services"
	"time"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	eventSvc services.IEventSvc
	jwtSvc   *commons.JwtSvc
	timezone *time.Location
}

func NewEventHandler(eventSvc services.IEventSvc, jwtSvc *commons.JwtSvc) *Handler {
	return &Handler{
		eventSvc: eventSvc,
		jwtSvc:   jwtSvc,
		timezone: time.FixedZone(commons.IST_TIMEZONE, commons.IST_OFFSET),
	}
}

func (h *Handler) GetRoutes() http.Handler {
	eventRouter := chi.NewRouter()

	eventRouter.Use(middlewares.JWTAuth(h.jwtSvc))

	eventRouter.Post("/", h.createEvent)

	eventRouter.Get("/{eventID}", h.getEvent)

	eventRouter.Delete("/{eventID}", h.deleteEvent)

	eventRouter.Get("/", h.getEvents)

	return eventRouter
}
