package main

import (
	"chat/internal/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := routes()

	log.Println("Starting channel listener...")
	go handlers.ListenToWsChannel()

	log.Println("Starting web server...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	_ = http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}
