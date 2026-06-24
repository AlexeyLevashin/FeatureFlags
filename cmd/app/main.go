package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Запуск Feature Flags API на порту 8080...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Критическая ошибка сервера: %v", err)
	}
}
