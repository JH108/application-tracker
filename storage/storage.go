package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	log.Printf("Initializing storage...")

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		log.Printf("Data directory does not exist, creating: %s", dataDir)
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			log.Printf("ERROR: Failed to create data directory: %v", err)
			return fmt.Errorf("failed to create data directory: %w", err)
		}
		log.Printf("Created data directory: %s", dataDir)
	} else {
		log.Printf("Data directory exists: %s", dataDir)
	}

	// Create applications file if it doesn't exist
	filePath := filepath.Join(dataDir, applicationsFile)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("Applications file does not exist, creating: %s", filePath)

		// Create an empty applications array
		emptyData := []models.Application{}
		jsonData, err := json.MarshalIndent(emptyData, "", "  ")
		if err != nil {
			log.Printf("ERROR: Failed to marshal empty applications: %v", err)
			return fmt.Errorf("failed to marshal empty applications: %w", err)
		}

		if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
			log.Printf("ERROR: Failed to create applications file: %v", err)
			return fmt.Errorf("failed to create applications file: %w", err)
		}

		log.Printf("Created applications file: %s", filePath)
	} else {
		log.Printf("Applications file exists: %s", filePath)

		// Verify the file contains valid JSON
		if err := validateApplicationsFile(filePath); err != nil {
			log.Printf("WARNING: Applications file contains invalid JSON: %v", err)
			log.Printf("You may want to fix or recreate the file to avoid runtime errors")
			// We don't return an error here to allow the application to start,
			// but we log a warning so the user knows there might be issues
		}
	}

	log.Printf("Storage initialization complete")
	return nil
}

// validateApplicationsFile checks if the applications file contains valid JSON
func validateApplicationsFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read applications file: %w", err)
	}

	var applications []models.Application
	if err := json.Unmarshal(data, &applications); err != nil {
		return fmt.Errorf("invalid JSON in applications file: %w", err)
	}

	return nil
}

// GetAllApplications returns all applications
func GetAllApplications() ([]models.Application, error) {
	log.Printf("Getting all applications")
	mutex.RLock()
	defer mutex.RUnlock()

	filePath := filepath.Join(dataDir, applicationsFile)
	data, err := os.ReadFile(filePath)

	log.Printf("Reading applications file: %s", filePath)

	if err != nil {
		log.Printf("ERROR: Failed to read applications file: %v", err)
		return nil, fmt.Errorf("failed to read applications file: %w", err)
	}

	var applications []models.Application
	if err := json.Unmarshal(data, &applications); err != nil {
		// This is a critical error - log it with details
		log.Printf("ERROR: Failed to unmarshal applications JSON: %v", err)
		// Include a snippet of the problematic JSON in the log
		if len(data) > 100 {
			log.Printf("JSON snippet (first 100 bytes): %s", string(data[:100]))
		} else {
			log.Printf("JSON content: %s", string(data))
		}
		return nil, fmt.Errorf("failed to unmarshal applications: %w", err)
	}

	log.Printf("Loaded applications from file: %s", filePath)
	log.Printf("Loaded %d applications from JSON", len(applications))

	return applications, nil
}

// GetApplicationByID returns an application by ID
func GetApplicationByID(id string) (*models.Application, error) {
	log.Printf("Getting application by ID: %s", id)

	applications, err := GetAllApplications()
	if err != nil {
		log.Printf("ERROR: Failed to get applications when looking up by ID: %v", err)
		return nil, err
	}

	for _, app := range applications {
		if app.ID == id {
			log.Printf("Found application: %s - %s (ID: %s)", app.Company, app.Position, app.ID)
			return &app, nil
		}
	}

	log.Printf("ERROR: Application with ID %s not found", id)
	return nil, ErrNotFound
}

// SaveApplication saves an application (creates or updates)
func SaveApplication(app *models.Application) error {
	log.Printf("Saving application: %s - %s", app.Company, app.ID)

	applications, err := GetAllApplications()
	if err != nil {
		log.Printf("ERROR: Failed to get applications for saving: %v", err)
		return err
	}

	log.Printf("Retrieved %d existing applications for update", len(applications))

	// Check if application already exists
	found := false
	for i, a := range applications {
		if a.ID == app.ID {
			// Update existing application
			applications[i] = *app
			found = true
			log.Printf("Updated existing application: %s (ID: %s)", app.Company, app.ID)
			break
		}
	}

	// Add new application if not found
	if !found {
		applications = append(applications, *app)
		log.Printf("Added new application: %s (ID: %s)", app.Company, app.ID)
	}

	// Save to file
	return saveApplicationsToFile(applications)
}

// DeleteApplication deletes an application by ID
func DeleteApplication(id string) error {
	log.Printf("Deleting application with ID: %s", id)

	applications, err := GetAllApplications()
	if err != nil {
		log.Printf("ERROR: Failed to get applications for deletion: %v", err)
		return err
	}

	found := false
	var updatedApps []models.Application

	for _, app := range applications {
		if app.ID != id {
			updatedApps = append(updatedApps, app)
		} else {
			found = true
			log.Printf("Found application to delete: %s - %s", app.Company, app.ID)
		}
	}

	if !found {
		log.Printf("ERROR: Application with ID %s not found for deletion", id)
		return ErrNotFound
	}

	log.Printf("Deleting application with ID %s (removed from %d total applications)", 
		id, len(applications))
	return saveApplicationsToFile(updatedApps)
}

// SearchApplications searches applications by tags and text
func SearchApplications(query string, tags []string) ([]models.Application, error) {
	log.Printf("Searching applications with query: '%s', tags: %v", query, tags)

	applications, err := GetAllApplications()
	if err != nil {
		log.Printf("ERROR: Failed to get applications for search: %v", err)
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

	log.Printf("Search returned %d results", len(results))
	return results, nil
}

// saveApplicationsToFile saves applications to the JSON file
func saveApplicationsToFile(applications []models.Application) error {
	filePath := filepath.Join(dataDir, applicationsFile)
	log.Printf("Saving %d applications to file: %s", len(applications), filePath)

	jsonData, err := json.MarshalIndent(applications, "", "  ")
	if err != nil {
		log.Printf("ERROR: Failed to marshal applications to JSON: %v", err)
		return fmt.Errorf("failed to marshal applications: %w", err)
	}

	// Log a snippet of the JSON data for debugging
	if len(jsonData) > 100 {
		log.Printf("JSON data snippet (first 100 bytes): %s", string(jsonData[:100]))
	} else {
		log.Printf("JSON data: %s", string(jsonData))
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		log.Printf("ERROR: Failed to write applications file: %v", err)
		return fmt.Errorf("failed to write applications file: %w", err)
	}

	log.Printf("Successfully saved applications to file")
	return nil
}
