package main

import (
	"net/http"

	"chat/internal/handlers"

	"github.com/bmizerany/pat"
)

func routes() http.Handler {
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/connect", http.HandlerFunc(handlers.Connect))

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
