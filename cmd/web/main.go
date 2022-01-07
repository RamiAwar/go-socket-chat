package main

import (
	"chat/internal/handlers"
	"log"
	"net/http"
)

func main() {
	mux := routes()

	log.Println("Starting channel listener...")
	go handlers.ListenToWsChannel()

	log.Println("Starting web server...")
	_ = http.ListenAndServe(":8000", mux)
}
