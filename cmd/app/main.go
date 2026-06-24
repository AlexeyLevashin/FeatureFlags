package main

import (
	http2 "FeatureFlags/internal/transport/http"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/auth/login", http2.LoginHandler)

	fmt.Println("running on port 8080")
	http.ListenAndServe(":8080", nil)
}
