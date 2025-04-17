# UI Implementation Plan for Application Tracker

## Overview

This document outlines the approach for implementing a UI layer for the Application Tracker using HTMX and Tailwind CSS. The UI will provide an interactive interface for managing job applications while leveraging the existing backend API.

## Architecture

The UI implementation will follow these architectural principles:

1. **Server-Side Rendering** - Go templates will render HTML on the server
2. **Progressive Enhancement** - Core functionality works without JavaScript
3. **HTMX for Interactivity** - HTMX will handle dynamic updates without full page reloads
4. **Tailwind CSS for Styling** - Utility-first CSS framework for responsive design

## Folder Structure

We'll extend the current project structure with the following additions:

```
ApplicationTracker/
├── static/              # Static assets
│   ├── css/             # CSS files (including Tailwind)
│   ├── js/              # JavaScript files (including HTMX)
│   └── images/          # Image assets
├── templates/           # HTML templates
│   ├── layouts/         # Base layout templates
│   │   └── base.html    # Main layout template
│   ├── partials/        # Reusable template components
│   │   ├── header.html
│   │   ├── footer.html
│   │   └── application-card.html
│   ├── pages/           # Full page templates
│   │   ├── index.html
│   │   ├── applications/
│   │   │   ├── list.html
│   │   │   ├── detail.html
│   │   │   ├── create.html
│   │   │   └── edit.html
│   └── htmx/            # HTMX partial templates for dynamic updates
│       └── applications/
│           ├── list-item.html
│           ├── form.html
│           └── search-results.html
├── ui/                  # UI-related Go code
│   ├── handlers.go      # UI route handlers
│   ├── middleware.go    # UI-specific middleware
│   └── router.go        # UI routes setup
```

## Implementation Steps

### 1. Set Up Template Engine

Add Go's template package to render HTML templates:

```go
// ui/handlers.go
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
    templates := template.Must(template.ParseGlob("templates/**/*.html"))
    err := templates.ExecuteTemplate(w, tmpl, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
```

### 2. Add Static File Serving

Configure the router to serve static files:

```go
// ui/router.go
func SetupUIRouter(mux *http.ServeMux) {
    // Serve static files
    fileServer := http.FileServer(http.Dir("static"))
    mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

    // UI routes
    mux.HandleFunc("/", homeHandler)
    mux.HandleFunc("/applications", applicationsListHandler)
    mux.HandleFunc("/applications/new", newApplicationHandler)
    mux.HandleFunc("/applications/", applicationDetailHandler)
}
```

### 3. Integrate HTMX and Tailwind CSS

Add the necessary libraries to the base template:

```html
<!-- templates/layouts/base.html -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Application Tracker</title>
    <link href="/static/css/tailwind.min.css" rel="stylesheet">
    <script src="/static/js/htmx.min.js"></script>
</head>
<body class="bg-gray-100">
    <div class="container mx-auto px-4 py-8">
        {{ template "content" . }}
    </div>
</body>
</html>
```

### 4. Create UI Pages

#### Home/Dashboard Page

```html
<!-- templates/pages/index.html -->
{{ define "content" }}
<div class="mb-8">
    <h1 class="text-3xl font-bold mb-4">Application Tracker</h1>
    <p class="text-gray-600">Track and manage your job applications</p>
    <a href="/applications/new" class="mt-4 inline-block bg-blue-500 text-white px-4 py-2 rounded">
        Add New Application
    </a>
</div>

<div class="mb-8">
    <h2 class="text-2xl font-bold mb-4">Recent Applications</h2>
    <div id="applications-list" hx-get="/htmx/applications" hx-trigger="load" hx-swap="innerHTML">
        <p>Loading applications...</p>
    </div>
</div>
{{ end }}
```

#### Application List

```html
<!-- templates/htmx/applications/list.html -->
{{ range .Applications }}
<div class="bg-white p-4 rounded shadow mb-4">
    <h3 class="text-xl font-bold">{{ .Company }} - {{ .Position }}</h3>
    <div class="flex mt-2">
        <span class="px-2 py-1 bg-gray-200 rounded text-sm mr-2">{{ .Status }}</span>
        {{ range .Tags }}
        <span class="px-2 py-1 bg-blue-100 text-blue-800 rounded text-sm mr-2">{{ . }}</span>
        {{ end }}
    </div>
    <div class="mt-4">
        <a href="/applications/{{ .ID }}" class="text-blue-500 hover:underline">View Details</a>
    </div>
</div>
{{ else }}
<p>No applications found.</p>
{{ end }}
```

### 5. Implement HTMX Interactions

#### Search with HTMX

```html
<!-- templates/partials/search-form.html -->
<form hx-get="/htmx/applications/search" hx-target="#applications-list" hx-trigger="submit, input[name='query'] changed delay:500ms">
    <div class="flex">
        <input type="text" name="query" placeholder="Search applications..." 
               class="px-4 py-2 border rounded-l w-full">
        <button type="submit" class="bg-blue-500 text-white px-4 py-2 rounded-r">
            Search
        </button>
    </div>
    <div class="mt-2">
        <input type="text" name="tags" placeholder="Tags (comma separated)" 
               class="px-4 py-2 border rounded w-full">
    </div>
</form>
```

#### HTMX Handler in Go

```go
// ui/handlers.go
func searchApplicationsHtmxHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("query")
    tagsParam := r.URL.Query().Get("tags")

    var tags []string
    if tagsParam != "" {
        tags = strings.Split(tagsParam, ",")
    }

    applications, err := storage.SearchApplications(query, tags)
    if err != nil {
        http.Error(w, "Failed to search applications", http.StatusInternalServerError)
        return
    }

    renderTemplate(w, "htmx/applications/list.html", map[string]interface{}{
        "Applications": applications,
    })
}
```

### 6. Form Handling with HTMX

```html
<!-- templates/pages/applications/create.html -->
{{ define "content" }}
<h1 class="text-3xl font-bold mb-6">Add New Application</h1>

<form hx-post="/api/applications" hx-swap="outerHTML" class="bg-white p-6 rounded shadow">
    <div class="mb-4">
        <label class="block text-gray-700 mb-2" for="company">Company</label>
        <input class="w-full px-3 py-2 border rounded" type="text" id="company" name="company" required>
    </div>

    <div class="mb-4">
        <label class="block text-gray-700 mb-2" for="position">Position</label>
        <input class="w-full px-3 py-2 border rounded" type="text" id="position" name="position" required>
    </div>

    <div class="mb-4">
        <label class="block text-gray-700 mb-2" for="description">Description</label>
        <textarea class="w-full px-3 py-2 border rounded" id="description" name="description" rows="3"></textarea>
    </div>

    <div class="mb-4">
        <label class="block text-gray-700 mb-2" for="url">URL</label>
        <input class="w-full px-3 py-2 border rounded" type="url" id="url" name="url">
    </div>

    <div class="mb-4">
        <label class="block text-gray-700 mb-2" for="tags">Tags (comma separated)</label>
        <input class="w-full px-3 py-2 border rounded" type="text" id="tags" name="tags">
    </div>

    <div class="mt-6">
        <button type="submit" class="bg-blue-500 text-white px-4 py-2 rounded">
            Save Application
        </button>
        <a href="/applications" class="ml-2 text-gray-500 hover:underline">Cancel</a>
    </div>
</form>
{{ end }}
```

## API Modifications for HTMX

To support HTMX, we'll need to add the following to our existing API:

1. Content negotiation to return HTML or JSON based on request headers
2. New endpoints for HTMX-specific partial responses
3. Support for form submissions with proper redirects

```go
// api/handlers.go
func isHtmxRequest(r *http.Request) bool {
    return r.Header.Get("HX-Request") == "true"
}

func respondWithHtml(w http.ResponseWriter, tmpl string, data interface{}) {
    w.Header().Set("Content-Type", "text/html")
    renderTemplate(w, tmpl, data)
}

// Modified handler example
func CreateApplicationHandler(w http.ResponseWriter, r *http.Request) {
    // ... existing code ...

    if isHtmxRequest(r) {
        // Return HTML response for HTMX
        w.Header().Set("HX-Redirect", "/applications")
        w.WriteHeader(http.StatusOK)
        return
    }

    // Return JSON response for API clients
    respondWithJSON(w, http.StatusCreated, Response{
        Success: true,
        Message: "Application created successfully",
        Data:    application,
    })
}
```

## Conclusion

This implementation plan provides a comprehensive approach to adding a UI layer to the Application Tracker using HTMX and Tailwind CSS. The plan leverages the existing backend API while adding server-side rendering capabilities and interactive UI components.

Key benefits of this approach:

1. **Minimal JavaScript** - HTMX handles most dynamic interactions
2. **Fast initial page load** - Server-side rendering provides quick first paint
3. **Progressive enhancement** - Core functionality works without JavaScript
4. **Responsive design** - Tailwind CSS makes responsive layouts easy
5. **Maintainable structure** - Clear separation of concerns between UI and API

Next steps would be to implement this plan, starting with the basic infrastructure (templates, static file serving) and then building out the UI components and HTMX interactions.
