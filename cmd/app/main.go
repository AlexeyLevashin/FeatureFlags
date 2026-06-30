package main

import (
	"FeatureFlags/internal/config"
	"FeatureFlags/internal/repository"
	"FeatureFlags/internal/service"
	"FeatureFlags/internal/transport/handlers"
	pkgPostgres "FeatureFlags/pkg/postgres"
	"log"
	"net/http"

	_ "FeatureFlags/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Feature Flags API
// @version 1.0
// @description API для управления фича-флагами
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg := config.LoadConfig()

	db, err := pkgPostgres.New(cfg.DatabaseDSN)

	if err != nil {
		log.Fatalf("Не удалось запустить БД: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка при закрытии соединения с БД: %v", err)
		}
	}()

	userRepo := repository.NewUserRepository(db)
	flagRepo := repository.NewFlagRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	flagService := service.NewFlagService(flagRepo, userRepo, teamRepo)
	authHandler := handlers.NewAuthHandler(authService)
	flagHandler := handlers.NewFlagHandler(flagService)

	r := chi.NewRouter()

	r.Post("/auth/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(handlers.AuthMiddleware(cfg.JWTSecret))

		r.Get("/me", authHandler.GetMe)
		r.Get("/flags", flagHandler.GetAllFlags)
		r.Get("/flags/{id}", flagHandler.GetFlagById)
		r.Post("/flags", flagHandler.CreateFlag)
		r.Put("/flags/{id}", flagHandler.UpdateFlagById)
		r.Patch("/flags/{id}/status", flagHandler.UpdateFlagStatusById)
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Println("Запуск Feature Flags API на порту 8080...")
	log.Println("Swagger доступен по адресу: http://localhost:8080/swagger/index.html")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Критическая ошибка сервера: %v", err)
	}
}
