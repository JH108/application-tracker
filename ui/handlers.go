package ui

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ApplicationTracker/models"
	"ApplicationTracker/storage"
)

// TemplateData holds data to be passed to templates
type TemplateData struct {
	Title        string
	CurrentYear  int
	Application  *models.Application
	Applications []models.Application
	Error        string
	Query        string
	Tags         string
	Status       string
	Page         int
}

// renderTemplate renders a template with the given data
func renderTemplate(w http.ResponseWriter, tmpl string, data TemplateData) {
	// Add current year to all template data
	data.CurrentYear = time.Now().Year()

	// Parse templates
	templates := template.Must(template.ParseGlob("templates/layouts/*.html"))
	templates = template.Must(templates.ParseGlob("templates/pages/*.html"))
	templates = template.Must(templates.ParseGlob("templates/pages/*/*.html"))
	templates = template.Must(templates.ParseGlob("templates/partials/*.html"))

	// Execute template
	// First set the content template, then execute the base template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Use the tmpl parameter to determine which template to execute
	// For the form template, we need to execute the form.html template directly
	if tmpl == "form" {
		// Look for the form template in the pages/applications directory
		formTemplate := template.Must(template.ParseFiles(
			"templates/layouts/base.html",
			"templates/partials/header.html",
			"templates/partials/footer.html",
			"templates/pages/applications/form.html",
		))
		err := formTemplate.ExecuteTemplate(w, "base.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else if tmpl == "detail" {
		// Look for the detail template in the pages/applications directory
		detailTemplate := template.Must(template.ParseFiles(
			"templates/layouts/base.html",
			"templates/partials/header.html",
			"templates/partials/footer.html",
			"templates/pages/applications/detail.html",
		))
		err := detailTemplate.ExecuteTemplate(w, "base.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// For other templates, execute the base.html template
	err := templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HomeHandler handles the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "content", TemplateData{
		Title: "Home",
	})
}

// ApplicationsListHandler handles the applications list page
func ApplicationsListHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "content", TemplateData{
		Title: "Applications",
	})
}

// ApplicationDetailHandler handles the application detail page
func ApplicationDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	id := strings.TrimPrefix(r.URL.Path, "/applications/")
	id = strings.TrimSuffix(id, "/")

	// Check if this is an edit request
	if strings.HasSuffix(id, "/edit") {
		id = strings.TrimSuffix(id, "/edit")
		ApplicationEditHandler(w, r, id)
		return
	}

	// Get application
	application, err := storage.GetApplicationByID(id)
	if err != nil {
		if err == storage.ErrNotFound {
			http.Error(w, "Application not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve application", http.StatusInternalServerError)
		}
		return
	}

	renderTemplate(w, "content", TemplateData{
		Title:       application.Company + " - " + application.Position,
		Application: application,
	})
}

// NewApplicationHandler handles the new application page
func NewApplicationHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "form", TemplateData{
		Title:       "Add New Application",
		Application: &models.Application{}, // Pass an empty application object
	})
}

// ApplicationEditHandler handles the edit application page
func ApplicationEditHandler(w http.ResponseWriter, r *http.Request, id string) {
	// Get application
	application, err := storage.GetApplicationByID(id)
	if err != nil {
		if err == storage.ErrNotFound {
			http.Error(w, "Application not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve application", http.StatusInternalServerError)
		}
		return
	}

	renderTemplate(w, "content", TemplateData{
		Title:       "Edit Application - " + application.Company,
		Application: application,
	})
}

// HtmxApplicationsHandler handles HTMX requests for applications list
func HtmxApplicationsHandler(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	query := r.URL.Query().Get("q")
	tagsParam := r.URL.Query().Get("tags")
	status := r.URL.Query().Get("status")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	// Parse page number
	page := 1
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
	}

	// Parse page size
	pageSize := 10
	if pageSizeStr != "" {
		var err error
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			pageSize = 10
		}
		// Limit page size to valid options
		if pageSize != 10 && pageSize != 25 && pageSize != 50 {
			pageSize = 10
		}
	}

	// Parse tags
	var tags []string
	if tagsParam != "" {
		tags = strings.Split(tagsParam, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
	}

	// Get applications
	applications, err := storage.SearchApplications(query, tags)
	if err != nil {
		http.Error(w, "Failed to search applications", http.StatusInternalServerError)
		return
	}

	// Filter by status if provided
	if status != "" {
		var filtered []models.Application
		for _, app := range applications {
			if app.Status == status {
				filtered = append(filtered, app)
			}
		}
		applications = filtered
	}

	// Calculate total count for pagination
	totalCount := len(applications)

	// Pagination
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize

	if startIndex >= totalCount {
		applications = []models.Application{}
	} else if endIndex > totalCount {
		applications = applications[startIndex:]
	} else {
		applications = applications[startIndex:endIndex]
	}

	// Calculate pagination metadata
	totalPages := (totalCount + pageSize - 1) / pageSize
	hasNextPage := page < totalPages

	// Set HX-Has-More header if there are more results
	if hasNextPage {
		w.Header().Set("HX-Has-More", "true")
	} else {
		w.Header().Set("HX-Has-More", "false")
	}

	// Set custom headers for pagination metadata
	w.Header().Set("HX-Current-Page", strconv.Itoa(page))
	w.Header().Set("HX-Page-Size", strconv.Itoa(pageSize))
	w.Header().Set("HX-Total-Pages", strconv.Itoa(totalPages))
	w.Header().Set("HX-Total-Count", strconv.Itoa(totalCount))

	// Render template
	tmpl := template.Must(template.ParseFiles("templates/htmx/applications/list.html"))
	err = tmpl.Execute(w, map[string]interface{}{
		"Applications": applications,
		"Pagination": map[string]interface{}{
			"CurrentPage": page,
			"PageSize":    pageSize,
			"TotalPages":  totalPages,
			"TotalCount":  totalCount,
			"HasNextPage": hasNextPage,
			"HasPrevPage": page > 1,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HtmxApplicationsCountHandler handles HTMX requests for applications count
func HtmxApplicationsCountHandler(w http.ResponseWriter, r *http.Request) {
	// Get applications
	applications, err := storage.GetAllApplications()
	if err != nil {
		http.Error(w, "Failed to retrieve applications", http.StatusInternalServerError)
		return
	}

	// Write count
	w.Write([]byte(strconv.Itoa(len(applications))))
}

// HtmxStatsHandler handles HTMX requests for application statistics
func HtmxStatsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract stat type from URL
	statType := strings.TrimPrefix(r.URL.Path, "/htmx/stats/")

	// Get applications
	applications, err := storage.GetAllApplications()
	if err != nil {
		http.Error(w, "Failed to retrieve applications", http.StatusInternalServerError)
		return
	}

	// Calculate stats
	var count int
	switch statType {
	case "total":
		count = len(applications)
	case "in-progress":
		for _, app := range applications {
			if app.Status == models.ApplicationStatus.InProgress {
				count++
			}
		}
	case "accepted":
		for _, app := range applications {
			if app.Status == models.ApplicationStatus.Accepted {
				count++
			}
		}
	case "rejected":
		for _, app := range applications {
			if app.Status == models.ApplicationStatus.Rejected {
				count++
			}
		}
	default:
		http.Error(w, "Invalid stat type", http.StatusBadRequest)
		return
	}

	// Write count
	w.Write([]byte(strconv.Itoa(count)))
}
