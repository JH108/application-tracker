#!/bin/bash

# Test script for ApplicationTracker API
# This script tests the basic CRUD operations and search functionality

# Set the API base URL
API_URL="http://localhost:8080/api"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print success or failure
print_result() {
  if [ $1 -eq 0 ]; then
    echo -e "${GREEN}SUCCESS${NC}: $2"
  else
    echo -e "${RED}FAILED${NC}: $2"
  fi
}

echo "Starting API tests..."

# Test health endpoint
echo -e "\n--- Testing Health Endpoint ---"
response=$(curl -s -o /dev/null -w "%{http_code}" $API_URL/health)
print_result $? "Health check returned status code: $response"

# Create a new application
echo -e "\n--- Testing Create Application ---"
create_response=$(curl -s -X POST $API_URL/applications \
  -H "Content-Type: application/json" \
  -d '{
    "company": "Example Corp",
    "position": "Software Engineer",
    "description": "Building awesome software",
    "url": "https://example.com/jobs/123",
    "tags": ["remote", "golang", "backend"]
  }')

# Extract the ID from the response
app_id=$(echo $create_response | grep -o '"id":"[^"]*' | cut -d'"' -f4)
print_result $? "Created application with ID: $app_id"

# Get all applications
echo -e "\n--- Testing Get All Applications ---"
get_all_response=$(curl -s $API_URL/applications)
print_result $? "Retrieved all applications"

# Get application by ID
echo -e "\n--- Testing Get Application by ID ---"
get_response=$(curl -s $API_URL/applications/$app_id)
print_result $? "Retrieved application with ID: $app_id"

# Update application
echo -e "\n--- Testing Update Application ---"
update_response=$(curl -s -X PUT $API_URL/applications/$app_id \
  -H "Content-Type: application/json" \
  -d '{
    "status": "in_progress",
    "tags": ["remote", "golang", "backend", "senior"]
  }')
print_result $? "Updated application with ID: $app_id"

# Search applications by tag
echo -e "\n--- Testing Search Applications by Tag ---"
search_response=$(curl -s "$API_URL/applications/search?tags=golang,remote")
print_result $? "Searched applications with tags: golang, remote"

# Search applications by query
echo -e "\n--- Testing Search Applications by Query ---"
search_response=$(curl -s "$API_URL/applications/search?q=Software")
print_result $? "Searched applications with query: Software"

# Delete application
echo -e "\n--- Testing Delete Application ---"
delete_response=$(curl -s -X DELETE $API_URL/applications/$app_id)
print_result $? "Deleted application with ID: $app_id"

# Verify deletion
echo -e "\n--- Verifying Deletion ---"
get_deleted_response=$(curl -s -o /dev/null -w "%{http_code}" $API_URL/applications/$app_id)
if [ $get_deleted_response -eq 404 ]; then
  print_result 0 "Application was successfully deleted (404 Not Found)"
else
  print_result 1 "Application was not deleted, got status code: $get_deleted_response"
fi

echo -e "\nAPI tests completed."