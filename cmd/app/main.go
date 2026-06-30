package main

import (
	"FeatureFlags/internal/config"
	"FeatureFlags/internal/repository"
	"FeatureFlags/internal/service"
	"FeatureFlags/internal/transport"
	"FeatureFlags/internal/transport/handlers"
	"FeatureFlags/pkg/logger"
	pkgPostgres "FeatureFlags/pkg/postgres"
	"net/http"
	"os"

	_ "FeatureFlags/docs"
)

// @title Feature Flags API
// @version 1.0
// @description API для управления фича-флагами
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	log := logger.Setup()
	log.Info("Инициализация приложения...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Критическая ошибка при загрузке конфигурации", "error", err)
		os.Exit(1)
	}

	db, err := pkgPostgres.New(cfg.DatabaseDSN)
	if err != nil {
		log.Error("Не удалось запустить БД", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error("Ошибка при закрытии соединения с БД", "error", err)
		}
	}()

	userRepo := repository.NewUserRepository(db)
	flagRepo := repository.NewFlagRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	flagService := service.NewFlagService(flagRepo, userRepo, teamRepo)
	authHandler := handlers.NewAuthHandler(authService)
	flagHandler := handlers.NewFlagHandler(flagService)

	r := transport.NewRouter(authHandler, flagHandler, cfg.JWTSecret)

	log.Info("Запуск Feature Flags API", "port", 8080)
	log.Info("Swagger доступен", "url", "http://localhost:8080/swagger/index.html")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Error("Критическая ошибка сервера:", "error", err)
		os.Exit(1)
	}
}
