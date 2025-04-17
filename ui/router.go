package ui

import (
	"net/http"
)

// SetupUIRouter sets up the UI routes
func SetupUIRouter(mux *http.ServeMux) {
	// Serve static files
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	
	// UI routes
	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/applications", ApplicationsListHandler)
	mux.HandleFunc("/applications/new", NewApplicationHandler)
	mux.HandleFunc("/applications/", ApplicationDetailHandler)
	
	// HTMX routes
	mux.HandleFunc("/htmx/applications", HtmxApplicationsHandler)
	mux.HandleFunc("/htmx/applications/search", HtmxApplicationsHandler)
	mux.HandleFunc("/htmx/applications/count", HtmxApplicationsCountHandler)
	mux.HandleFunc("/htmx/stats/", HtmxStatsHandler)
}