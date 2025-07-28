package main

import (
	"log"
	"net/http"

	"github.com/gemini/go-blog-api/internal/database"
	"github.com/gemini/go-blog-api/internal/handler"
)

func main() {
	// Initialize the in-memory database
	db := database.NewMemoryStore()

	// Initialize handlers
	postHandler := handler.NewPostHandler(db)

	// Setup the router
	mux := http.NewServeMux()
	mux.Handle("/posts/", postHandler)
	mux.HandleFunc("/health", handler.HealthCheckHandler)

	// Configure the server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server starting on port 8080...")
	log.Fatal(server.ListenAndServe())
}
