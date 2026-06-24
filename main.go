package main

import (
	"feature-flags/handlers"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/auth/login", handlers.LoginHandler)

	fmt.Println("running on port 8080")
	http.ListenAndServe(":8080", nil)
}
