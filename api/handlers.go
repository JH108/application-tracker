package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"ApplicationTracker/models"
	"ApplicationTracker/storage"
)

// Response is a generic API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ApplicationRequest is the structure for application creation/update requests
type ApplicationRequest struct {
	Company     string   `json:"company"`
	Position    string   `json:"position"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Status      string   `json:"status,omitempty"`
	Tags        []string `json:"tags"`
}

// GetAllApplicationsHandler returns all applications
func GetAllApplicationsHandler(w http.ResponseWriter, r *http.Request) {
	applications, err := storage.GetAllApplications()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve applications: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    applications,
	})
}

// GetApplicationHandler returns a specific application by ID
func GetApplicationHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/applications/")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Application ID is required")
		return
	}

	application, err := storage.GetApplicationByID(id)
	if err != nil {
		if err == storage.ErrNotFound {
			respondWithError(w, http.StatusNotFound, "Application not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve application: "+err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    application,
	})
}

// CreateApplicationHandler creates a new application
func CreateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	var req ApplicationRequest

	// Handle form submissions from HTMX
	if isHtmxRequest(r) && r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid form data: "+err.Error())
			return
		}

		// Extract form values
		req.Company = r.FormValue("company")
		req.Position = r.FormValue("position")
		req.Description = r.FormValue("description")
		req.URL = r.FormValue("url")

		// Parse tags
		if tagsStr := r.FormValue("tags"); tagsStr != "" {
			tags := strings.Split(tagsStr, ",")
			for i, tag := range tags {
				tags[i] = strings.TrimSpace(tag)
			}
			req.Tags = tags
		}
	} else {
		// Handle JSON API requests
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
			return
		}
	}

	// Validate required fields
	if req.Company == "" || req.Position == "" {
		respondWithError(w, http.StatusBadRequest, "Company and position are required fields")
		return
	}

	// Create new application
	application := models.NewApplication(
		req.Company,
		req.Position,
		req.Description,
		req.URL,
		req.Tags,
	)

	// Set status if provided
	if req.Status != "" {
		application.UpdateStatus(req.Status)
	}

	// Save to storage
	if err := storage.SaveApplication(application); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to save application: "+err.Error())
		return
	}

	// Handle HTMX response
	if isHtmxRequest(r) {
		w.Header().Set("HX-Redirect", "/applications")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Handle JSON API response
	respondWithJSON(w, http.StatusCreated, Response{
		Success: true,
		Message: "Application created successfully",
		Data:    application,
	})
}

// UpdateApplicationHandler updates an existing application
func UpdateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/applications/")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Application ID is required")
		return
	}

	// Get existing application
	application, err := storage.GetApplicationByID(id)
	if err != nil {
		if err == storage.ErrNotFound {
			respondWithError(w, http.StatusNotFound, "Application not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve application: "+err.Error())
		}
		return
	}

	var req ApplicationRequest

	// Handle form submissions from HTMX
	if isHtmxRequest(r) && r.Method == http.MethodPut {
		if err := r.ParseForm(); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid form data: "+err.Error())
			return
		}

		// Extract form values
		req.Company = r.FormValue("company")
		req.Position = r.FormValue("position")
		req.Description = r.FormValue("description")
		req.URL = r.FormValue("url")
		req.Status = r.FormValue("status")

		// Parse tags
		if tagsStr := r.FormValue("tags"); tagsStr != "" {
			tags := strings.Split(tagsStr, ",")
			for i, tag := range tags {
				tags[i] = strings.TrimSpace(tag)
			}
			req.Tags = tags
		}
	} else {
		// Parse JSON request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
			return
		}
	}

	// Update fields
	if req.Company != "" {
		application.Company = req.Company
	}
	if req.Position != "" {
		application.Position = req.Position
	}
	application.Description = req.Description // Allow empty description
	application.URL = req.URL // Allow empty URL
	if req.Status != "" {
		application.Status = req.Status
	}
	if req.Tags != nil {
		application.Tags = req.Tags
	}

	// Update timestamp
	application.UpdatedAt = time.Now()

	// Save to storage
	if err := storage.SaveApplication(application); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update application: "+err.Error())
		return
	}

	// Handle HTMX response
	if isHtmxRequest(r) {
		w.Header().Set("HX-Redirect", "/applications/"+id)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Handle JSON API response
	respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "Application updated successfully",
		Data:    application,
	})
}

// UpdateApplicationStatusHandler updates the status of an application
func UpdateApplicationStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	path := strings.TrimPrefix(r.URL.Path, "/api/applications/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[1] != "status" {
		respondWithError(w, http.StatusBadRequest, "Invalid URL format")
		return
	}
	id := parts[0]

	// Get existing application
	application, err := storage.GetApplicationByID(id)
	if err != nil {
		if err == storage.ErrNotFound {
			respondWithError(w, http.StatusNotFound, "Application not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve application: "+err.Error())
		}
		return
	}

	// Get status from request
	var status string
	if isHtmxRequest(r) {
		if err := r.ParseForm(); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid form data: "+err.Error())
			return
		}
		status = r.FormValue("status")
	} else {
		var req struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
			return
		}
		status = req.Status
	}

	// Validate status
	validStatuses := []string{
		models.ApplicationStatus.Applied,
		models.ApplicationStatus.InProgress,
		models.ApplicationStatus.Accepted,
		models.ApplicationStatus.Rejected,
	}
	valid := false
	for _, s := range validStatuses {
		if status == s {
			valid = true
			break
		}
	}
	if !valid {
		respondWithError(w, http.StatusBadRequest, "Invalid status value")
		return
	}

	// Update status
	application.UpdateStatus(status)

	// Save to storage
	if err := storage.SaveApplication(application); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update application: "+err.Error())
		return
	}

	// Handle HTMX response
	if isHtmxRequest(r) {
		w.Header().Set("HX-Refresh", "true")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Handle JSON API response
	respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "Application status updated successfully",
		Data:    application,
	})
}

// DeleteApplicationHandler deletes an application
func DeleteApplicationHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/applications/")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Application ID is required")
		return
	}

	if err := storage.DeleteApplication(id); err != nil {
		if err == storage.ErrNotFound {
			respondWithError(w, http.StatusNotFound, "Application not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Failed to delete application: "+err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "Application deleted successfully",
	})
}

// SearchApplicationsHandler searches applications by query and tags
func SearchApplicationsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	// Get tags from query parameters
	var tags []string
	if tagsParam := r.URL.Query().Get("tags"); tagsParam != "" {
		tags = strings.Split(tagsParam, ",")
	}

	applications, err := storage.SearchApplications(query, tags)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to search applications: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    applications,
	})
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, Response{
		Success: false,
		Message: message,
	})
}

// isHtmxRequest checks if the request is from HTMX
func isHtmxRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to marshal JSON response"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
