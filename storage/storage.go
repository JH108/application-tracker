package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"ApplicationTracker/models"
)

const (
	dataDir          = "./data"
	applicationsFile = "applications.json"
)

var (
	// ErrNotFound is returned when an application is not found
	ErrNotFound = errors.New("application not found")

	// mutex to prevent concurrent file access
	mutex = &sync.RWMutex{}
)

// Initialize creates the data directory if it doesn't exist
func Initialize() error {
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return fmt.Errorf("failed to create data directory: %w", err)
		}
	}

	// Create applications file if it doesn't exist
	filePath := filepath.Join(dataDir, applicationsFile)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create an empty applications array
		emptyData := []models.Application{}
		jsonData, err := json.MarshalIndent(emptyData, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal empty applications: %w", err)
		}

		if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
			return fmt.Errorf("failed to create applications file: %w", err)
		}
	}

	return nil
}

// GetAllApplications returns all applications
func GetAllApplications() ([]models.Application, error) {
	fmt.Println("Getting all applications")
	mutex.RLock()
	defer mutex.RUnlock()

	filePath := filepath.Join(dataDir, applicationsFile)
	data, err := os.ReadFile(filePath)

	fmt.Println("Reading applications file:", filePath)

	if err != nil {
		return nil, fmt.Errorf("failed to read applications file: %w", err)
	}

	var applications []models.Application
	if err := json.Unmarshal(data, &applications); err != nil {
		return nil, fmt.Errorf("failed to unmarshal applications: %w", err)
	}

	fmt.Println("Loaded applications from file:", filePath)
	fmt.Println("Loaded applications from JSON:", applications)

	return applications, nil
}

// GetApplicationByID returns an application by ID
func GetApplicationByID(id string) (*models.Application, error) {
	applications, err := GetAllApplications()
	if err != nil {
		return nil, err
	}

	for _, app := range applications {
		if app.ID == id {
			return &app, nil
		}
	}

	return nil, ErrNotFound
}

// SaveApplication saves an application (creates or updates)
func SaveApplication(app *models.Application) error {
	fmt.Println("Saving application:", app)

	applications, err := GetAllApplications()
	if err != nil {
		return err
	}

	fmt.Println("Saving applications:", applications)

	// Check if application already exists
	found := false
	for i, a := range applications {
		if a.ID == app.ID {
			// Update existing application
			applications[i] = *app
			found = true
			break
		}
	}

	// Add new application if not found
	if !found {
		applications = append(applications, *app)
	}

	// Save to file
	return saveApplicationsToFile(applications)
}

// DeleteApplication deletes an application by ID
func DeleteApplication(id string) error {
	applications, err := GetAllApplications()
	if err != nil {
		return err
	}

	found := false
	var updatedApps []models.Application

	for _, app := range applications {
		if app.ID != id {
			updatedApps = append(updatedApps, app)
		} else {
			found = true
		}
	}

	if !found {
		return ErrNotFound
	}

	return saveApplicationsToFile(updatedApps)
}

// SearchApplications searches applications by tags and text
func SearchApplications(query string, tags []string) ([]models.Application, error) {
	applications, err := GetAllApplications()
	if err != nil {
		return nil, err
	}

	var results []models.Application

	for _, app := range applications {
		// Check if all specified tags are present
		if len(tags) > 0 {
			allTagsPresent := true
			for _, searchTag := range tags {
				tagFound := false
				for _, appTag := range app.Tags {
					if strings.EqualFold(appTag, searchTag) {
						tagFound = true
						break
					}
				}
				if !tagFound {
					allTagsPresent = false
					break
				}
			}
			if !allTagsPresent {
				continue
			}
		}

		// Check if query matches any field
		if query != "" {
			query = strings.ToLower(query)
			if !strings.Contains(strings.ToLower(app.Company), query) &&
				!strings.Contains(strings.ToLower(app.Position), query) &&
				!strings.Contains(strings.ToLower(app.Description), query) {
				continue
			}
		}

		results = append(results, app)
	}

	return results, nil
}

// saveApplicationsToFile saves applications to the JSON file
func saveApplicationsToFile(applications []models.Application) error {
	filePath := filepath.Join(dataDir, applicationsFile)

	jsonData, err := json.MarshalIndent(applications, "", "  ")

	fmt.Println("Saving applications to file:", filePath)
	fmt.Println("Saving applications to JSON:", string(jsonData))

	if err != nil {
		return fmt.Errorf("failed to marshal applications: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write applications file: %w", err)
	}

	return nil
}
