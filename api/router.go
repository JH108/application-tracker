package api

import (
	"ApplicationTracker/storage"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// Middleware for logging requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// Middleware for CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// applicationHandler handles all application-related requests
func applicationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request:", r.Method, r.URL.Path)
	// Extract the ID from the path if present
	path := strings.TrimPrefix(r.URL.Path, "/applications")

	// Route based on HTTP method and path
	switch {
	case r.Method == http.MethodGet && path == "":
		// GET /api/applications - Get all applications
		GetAllApplicationsHandler(w, r)

	case r.Method == http.MethodGet && path == "/search":
		// GET /api/applications/search - Search applications
		SearchApplicationsHandler(w, r)

	case r.Method == http.MethodGet && path != "":
		// GET /api/applications/{id} - Get application by ID
		GetApplicationHandler(w, r)

	case r.Method == http.MethodPost && path == "":
		// POST /api/applications - Create new application
		CreateApplicationHandler(w, r)

	case r.Method == http.MethodPut && strings.Contains(path, "/status"):
		// PUT /api/applications/{id}/status - Update application status
		UpdateApplicationStatusHandler(w, r)

	case r.Method == http.MethodPut && path != "":
		// PUT /api/applications/{id} - Update application
		UpdateApplicationHandler(w, r)

	case r.Method == http.MethodDelete && path != "":
		// DELETE /api/applications/{id} - Delete application
		DeleteApplicationHandler(w, r)

	default:
		// Method not allowed or route not found
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed or route not found"))
	}
}

// healthCheckHandler handles health check requests
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// SetupRouter initializes and returns the HTTP router
func SetupRouter() http.Handler {
	// Initialize storage
	if err := storage.Initialize(); err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/applications", applicationHandler)
	mux.HandleFunc("/applications/", applicationHandler)
	mux.HandleFunc("/health", healthCheckHandler)

	// Add middleware
	handler := loggingMiddleware(corsMiddleware(mux))

	return handler
}
