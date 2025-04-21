# Application Tracker End-to-End Tests

This directory contains end-to-end tests for the Application Tracker project using Playwright.

## Prerequisites

- Node.js (v14 or later)
- npm (v6 or later)

## Test Structure

The tests are organized as follows:

- `e2e/home.spec.js` - Tests for the home page
- `e2e/applications-list.spec.js` - Tests for the applications list page
- `e2e/application-detail.spec.js` - Tests for the application detail page
- `e2e/application-form.spec.js` - Tests for the application form page
- `e2e/api.spec.js` - Tests for the API endpoints
- `e2e/user-flows.spec.js` - Tests for complete user flows

## Running the Tests

### Using the Test Runner Script

The easiest way to run the tests is to use the provided test runner script:

```bash
./run-tests.sh
```

This script will:
1. Check if Node.js and npm are installed
2. Install dependencies if they're not already installed
3. Install Playwright browsers if they're not already installed
4. Run the tests
5. Report the test results

### Manual Setup and Running

If you prefer to run the tests manually, follow these steps:

1. Install dependencies:
   ```bash
   npm install
   ```

2. Install Playwright browsers:
   ```bash
   npx playwright install
   ```

3. Run the tests:
   ```bash
   npm test
   ```

### Running Specific Tests

To run a specific test file:

```bash
npx playwright test e2e/home.spec.js
```

To run tests in a specific browser:

```bash
npx playwright test --project=chromium
```

To run tests in debug mode:

```bash
npm run test:debug
```

To run tests with the Playwright UI:

```bash
npm run test:ui
```

## Viewing Test Reports

After running the tests, you can view the HTML report:

```bash
npx playwright show-report
```

## Test Configuration

The test configuration is defined in `playwright.config.js`. Key settings include:

- Tests run against a local server started automatically at http://localhost:8080
- Tests run in Chromium, Firefox, and WebKit browsers
- Screenshots are taken on test failures
- Traces are collected on first retry

## Adding New Tests

To add new tests:

1. Create a new test file in the `e2e` directory
2. Import the Playwright test utilities:
   ```javascript
   const { test, expect } = require('@playwright/test');
   ```
3. Write your tests using the Playwright API
4. Run the tests to verify they work as expected