package models

import (
	"time"
)

// Application represents a job application
type Application struct {
	ID          string    `json:"id"`
	Company     string    `json:"company"`
	Position    string    `json:"position"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Status      string    `json:"status"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ApplicationStatus defines the possible statuses for a job application
var ApplicationStatus = struct {
	Applied     string
	InProgress  string
	Rejected    string
	Accepted    string
}{
	Applied:     "applied",
	InProgress:  "in_progress",
	Rejected:    "rejected",
	Accepted:    "accepted",
}

// NewApplication creates a new application with default values
func NewApplication(company, position, description, url string, tags []string) *Application {
	now := time.Now()
	return &Application{
		ID:          generateID(),
		Company:     company,
		Position:    position,
		Description: description,
		URL:         url,
		Status:      ApplicationStatus.Applied,
		Tags:        tags,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// AddTag adds a tag to the application if it doesn't already exist
func (a *Application) AddTag(tag string) {
	for _, t := range a.Tags {
		if t == tag {
			return
		}
	}
	a.Tags = append(a.Tags, tag)
	a.UpdatedAt = time.Now()
}

// RemoveTag removes a tag from the application
func (a *Application) RemoveTag(tag string) {
	for i, t := range a.Tags {
		if t == tag {
			a.Tags = append(a.Tags[:i], a.Tags[i+1:]...)
			a.UpdatedAt = time.Now()
			return
		}
	}
}

// UpdateStatus updates the application status
func (a *Application) UpdateStatus(status string) {
	a.Status = status
	a.UpdatedAt = time.Now()
}