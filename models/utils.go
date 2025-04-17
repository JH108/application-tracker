package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// generateID creates a unique ID for an application
func generateID() string {
	// Create a random component
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		// Fallback to timestamp if random generation fails
		return fmt.Sprintf("app_%d", time.Now().UnixNano())
	}
	
	// Combine timestamp and random component for uniqueness
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("app_%d_%s", timestamp, hex.EncodeToString(randomBytes))
}