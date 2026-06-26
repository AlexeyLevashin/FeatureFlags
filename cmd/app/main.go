package main

import (
	"FeatureFlags/internal/config"
	"FeatureFlags/internal/repository"
	"FeatureFlags/internal/service"
	"FeatureFlags/internal/transport/http"
	pkgPostgres "FeatureFlags/pkg/postgres"
	"log"
	"net/http"
)

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
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)
	flagHandler := handlers.NewFlagHandler(flagRepo)
	//userRepo := repository.NewFlagRepository(&db)

	http.HandleFunc("/auth/login", authHandler.Login)
	http.HandleFunc("/flags", flagHandler.GetAllFlags)
	log.Println("Запуск Feature Flags API на порту 8080...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Критическая ошибка сервера: %v", err)
	}
}
