# Application Tracker

A backend service for tracking job applications with tagging and search capabilities.

## Overview

This application provides a complete solution for tracking job applications. It includes:

- A RESTful API for managing job applications
- A user-friendly web interface built with HTMX and Tailwind CSS
- Local storage using JSON files (JSON schema compliant)

### Features

- Create, read, update, and delete job applications
- Add and remove tags from applications
- Search applications by text and tags
- Filter applications by status
- Track application status changes
- Responsive design that works on mobile and desktop

## Getting Started

### Prerequisites

- Go 1.22 or higher

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/ApplicationTracker.git
   cd ApplicationTracker
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Build the application:
   ```
   go build -o app
   ```

4. Run the application:
   ```
   ./app
   ```

The server will start on port 8080 by default.

## API Endpoints

### Health Check

- `GET /api/health` - Check if the API is running

### Applications

- `GET /api/applications` - Get all applications
- `GET /api/applications/{id}` - Get application by ID
- `POST /api/applications` - Create a new application
- `PUT /api/applications/{id}` - Update an application
- `DELETE /api/applications/{id}` - Delete an application
- `GET /api/applications/search?q={query}&tags={tag1,tag2}` - Search applications by text and/or tags

## Data Model

### Application

```json
{
  "id": "string",
  "company": "string",
  "position": "string",
  "description": "string",
  "url": "string",
  "status": "string",
  "tags": ["string"],
  "createdAt": "string (ISO date)",
  "updatedAt": "string (ISO date)"
}
```

### Application Status Values

- `applied` - Initial application submitted
- `in_progress` - In the interview process
- `rejected` - Application was rejected
- `accepted` - Received an offer

## Example Requests

### Create Application

```bash
curl -X POST http://localhost:8080/api/applications \
  -H "Content-Type: application/json" \
  -d '{
    "company": "Example Corp",
    "position": "Software Engineer",
    "description": "Building awesome software",
    "url": "https://example.com/jobs/123",
    "tags": ["remote", "golang", "backend"]
  }'
```

### Search Applications

```bash
# Search by text
curl "http://localhost:8080/api/applications/search?q=Software"

# Search by tags
curl "http://localhost:8080/api/applications/search?tags=remote,golang"

# Search by both
curl "http://localhost:8080/api/applications/search?q=Engineer&tags=remote"
```

## Testing

### API Tests

Run the included test script to verify the API functionality:

```bash
chmod +x test_api.sh
./test_api.sh
```

### End-to-End Tests

The project includes a comprehensive suite of end-to-end tests using Playwright. These tests verify that all routes and user flows work as expected.

To run the end-to-end tests:

```bash
cd tests
./run-tests.sh
```

For more details on the end-to-end tests, see the [tests/README.md](tests/README.md) file.

## Project Structure

- `main.go` - Application entry point
- `models/` - Data models
- `storage/` - JSON file storage implementation
- `api/` - API handlers and routing
- `ui/` - UI handlers and routing
- `templates/` - HTML templates for the UI
  - `layouts/` - Base layout templates
  - `pages/` - Full page templates
  - `partials/` - Reusable template components
  - `htmx/` - HTMX partial templates for dynamic updates
- `static/` - Static assets (CSS, JS, images)
- `schemas/` - JSON schemas for data validation
- `tests/` - End-to-end tests using Playwright
  - `e2e/` - Test files organized by feature
  - `playwright.config.js` - Playwright configuration
  - `run-tests.sh` - Test runner script

## Technologies Used

- **Backend**: Go standard library
- **Frontend**: 
  - HTMX for interactive UI without writing JavaScript
  - Tailwind CSS for styling
- **Storage**: Local JSON files
- **Testing**:
  - Playwright for end-to-end testing
  - Bash scripts for API testing

## Future Enhancements

- User authentication
- Multiple user support
- Advanced statistics and reporting
- Email notifications for application status changes
- Calendar integration for interviews
- Document attachments (resume, cover letter)
