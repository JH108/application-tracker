package main

import (
	"fmt"
	"log"
	"net/http"

	"ApplicationTracker/api"
	"ApplicationTracker/storage"
	"ApplicationTracker/ui"
)

func main() {
	port := 8080

	// Initialize storage
	if err := storage.Initialize(); err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Set up API routes
	apiRouter := api.SetupRouter()
	mux.Handle("/api/", http.StripPrefix("/api", apiRouter))

	// Set up UI routes
	ui.SetupUIRouter(mux)

	// Start the server
	fmt.Printf("Server running on port %d...\n", port)
	fmt.Printf("UI available at http://localhost:%d\n", port)
	fmt.Printf("API available at http://localhost:%d/api\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
