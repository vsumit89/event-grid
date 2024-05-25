package handlers

import (
	"net/http"
	"server/internal/commons"
	"server/internal/handlers/eventHandlers"
	"server/internal/handlers/middlewares"
	"server/internal/handlers/userHandlers"
	"server/internal/services"

	"github.com/go-chi/chi/v5"
)

// Container is a struct that contains all the services dependencies for the router
type Container struct {
	UserSvc  services.IUserSvc
	JWTSvc   *commons.JwtSvc
	EventSvc services.IEventSvc
}

func GetRoutes(c *Container) http.Handler {
	router := chi.NewRouter()

	router.Use(middlewares.Logging)
	router.Use(middlewares.Recoverer)

	v1Router := chi.NewRouter()

	router.Mount("/api/v1", v1Router)

	userRouter := userHandlers.NewHandler(c.UserSvc).GetRoutes()

	authRouter := userHandlers.NewAuthHandler(c.UserSvc, c.JWTSvc).AuthRoutes()

	eventRouter := eventHandlers.NewEventHandler(c.EventSvc, c.JWTSvc).GetRoutes()

	v1Router.Mount("/users", userRouter)

	v1Router.Mount("/auth", authRouter)

	v1Router.Mount("/events", eventRouter)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return router
}
