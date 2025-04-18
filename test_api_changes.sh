#!/bin/bash

# Test script for ApplicationTracker API changes
# This script tests the UUID v4 format for IDs and empty URL field

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

echo "Starting API tests for changes..."

# Test 1: Create an application and verify UUID v4 format
echo -e "\n--- Testing UUID v4 Format for Application ID ---"
create_response=$(curl -s -X POST $API_URL/applications \
  -H "Content-Type: application/json" \
  -d '{
    "company": "Test UUID Corp",
    "position": "Software Engineer",
    "description": "Testing UUID format",
    "url": "https://example.com/jobs/uuid",
    "tags": ["test", "uuid"]
  }')

# Extract the ID from the response
app_id=$(echo $create_response | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "Created application with ID: $app_id"

# Verify UUID v4 format (8-4-4-4-12 hexadecimal characters)
if [[ $app_id =~ ^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$ ]]; then
  print_result 0 "Application ID is in UUID v4 format"
else
  print_result 1 "Application ID is not in UUID v4 format: $app_id"
fi

# Test 2: Create an application with empty URL
echo -e "\n--- Testing Create Application with Empty URL ---"
create_empty_url_response=$(curl -s -X POST $API_URL/applications \
  -H "Content-Type: application/json" \
  -d '{
    "company": "Empty URL Corp",
    "position": "Software Engineer",
    "description": "Testing empty URL",
    "url": "",
    "tags": ["test", "empty-url"]
  }')

# Extract the ID from the response
empty_url_app_id=$(echo $create_empty_url_response | grep -o '"id":"[^"]*' | cut -d'"' -f4)
print_result $? "Created application with empty URL, ID: $empty_url_app_id"

# Get the application to verify empty URL
get_empty_url_response=$(curl -s $API_URL/applications/$empty_url_app_id)
url_value=$(echo $get_empty_url_response | grep -o '"url":"[^"]*' | cut -d'"' -f4)

if [ "$url_value" = "" ]; then
  print_result 0 "Application has empty URL as expected"
else
  print_result 1 "Application URL is not empty: $url_value"
fi

# Test 3: Update an application to have empty URL
echo -e "\n--- Testing Update Application to Empty URL ---"
update_response=$(curl -s -X PUT $API_URL/applications/$app_id \
  -H "Content-Type: application/json" \
  -d '{
    "url": ""
  }')
print_result $? "Updated application with ID: $app_id to have empty URL"

# Get the application to verify empty URL after update
get_updated_response=$(curl -s $API_URL/applications/$app_id)
updated_url_value=$(echo $get_updated_response | grep -o '"url":"[^"]*' | cut -d'"' -f4)

if [ "$updated_url_value" = "" ]; then
  print_result 0 "Application URL was updated to empty as expected"
else
  print_result 1 "Application URL was not updated to empty: $updated_url_value"
fi

# Clean up - Delete the applications
echo -e "\n--- Cleaning Up ---"
curl -s -X DELETE $API_URL/applications/$app_id > /dev/null
print_result $? "Deleted application with ID: $app_id"

curl -s -X DELETE $API_URL/applications/$empty_url_app_id > /dev/null
print_result $? "Deleted application with ID: $empty_url_app_id"

echo -e "\nAPI tests for changes completed."