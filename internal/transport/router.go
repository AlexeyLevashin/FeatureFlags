package transport

import (
	"FeatureFlags/internal/transport/handlers"

	"FeatureFlags/internal/transport/middleware"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(
	authHandler *handlers.AuthHandler,
	flagHandler *handlers.FlagHandler,
	jwtSecret string,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.CorsMiddleware)

	r.Post("/auth/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(handlers.AuthMiddleware(jwtSecret))

		r.Get("/me", authHandler.GetMe)

		r.Get("/flags", flagHandler.GetAllFlags)
		r.Get("/flags/{id}", flagHandler.GetFlagDetailsById)
		r.Post("/flags", flagHandler.CreateFlag)
		r.Put("/flags/{id}", flagHandler.UpdateFlagById)
		r.Patch("/flags/{id}/status", flagHandler.UpdateFlagStatusById)
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r
}
