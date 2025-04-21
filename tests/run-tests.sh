#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== Application Tracker E2E Tests ===${NC}"
echo -e "This script will run the end-to-end tests for the Application Tracker."

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo -e "${RED}Error: Node.js is not installed. Please install Node.js to run the tests.${NC}"
    exit 1
fi

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    echo -e "${RED}Error: npm is not installed. Please install npm to run the tests.${NC}"
    exit 1
fi

# Navigate to the tests directory
cd "$(dirname "$0")"

# Install dependencies if node_modules doesn't exist
if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}Installing dependencies...${NC}"
    npm install
    if [ $? -ne 0 ]; then
        echo -e "${RED}Failed to install dependencies.${NC}"
        exit 1
    fi
    echo -e "${GREEN}Dependencies installed successfully.${NC}"
fi

# Install Playwright browsers if not already installed
if [ ! -d "node_modules/.cache/ms-playwright" ]; then
    echo -e "${YELLOW}Installing Playwright browsers...${NC}"
    npx playwright install
    if [ $? -ne 0 ]; then
        echo -e "${RED}Failed to install Playwright browsers.${NC}"
        exit 1
    fi
    echo -e "${GREEN}Playwright browsers installed successfully.${NC}"
fi

# Run the tests
echo -e "${YELLOW}Running tests...${NC}"
npm test

# Check the exit code
if [ $? -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
else
    echo -e "${RED}Some tests failed. Check the test report for details.${NC}"
    echo -e "You can view the HTML report by running: npx playwright show-report"
fi