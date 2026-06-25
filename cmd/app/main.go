package main

import (
	"FeatureFlags/internal/repository"
	"FeatureFlags/internal/service"
	"FeatureFlags/internal/transport/http"
	"log"
	"net/http"
)

func main() {
	// userRepo := repository.NewUserRepository(db)
	userRepo := repository.UserRepository{}
	authService := service.NewAuthService(&userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	http.HandleFunc("/auth/login", authHandler.Login)
	log.Println("Запуск Feature Flags API на порту 8080...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Критическая ошибка сервера: %v", err)
	}
}
