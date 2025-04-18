package models

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

// generateID creates a unique ID for an application using UUID v4
func generateID() string {
	// Generate a UUID v4
	id, err := uuid.NewRandom()
	if err != nil {
		// Fallback to timestamp if UUID generation fails
		return fmt.Sprintf("app_%d", time.Now().UnixNano())
	}

	return id.String()
}
