package main

import (
	http2 "FeatureFlags/internal/transport/http"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/auth/login", http2.LoginHandler)
	log.Println("Запуск Feature Flags API на порту 8080...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Критическая ошибка сервера: %v", err)
	}
}
