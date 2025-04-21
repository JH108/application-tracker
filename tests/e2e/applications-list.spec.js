const { test, expect } = require('@playwright/test');

test.describe('Applications List Page', () => {
  test('should load the applications list page', async ({ page }) => {
    // Navigate to the applications list page
    await page.goto('/applications');
    
    // Check that the page title is correct
    await expect(page).toHaveTitle(/Applications/);
    
    // Check that the main heading is present
    const heading = page.locator('h1');
    await expect(heading).toBeVisible();
    await expect(heading).toContainText('Applications');
    
    // Check that the "Add New Application" button is present
    const addButton = page.locator('a', { hasText: 'Add New Application' });
    await expect(addButton).toBeVisible();
    
    // Check that the search form is present
    const searchForm = page.locator('form[hx-get="/htmx/applications/search"]');
    await expect(searchForm).toBeVisible();
    
    // Check that the applications container is present
    const applicationsContainer = page.locator('#applications-container');
    await expect(applicationsContainer).toBeVisible();
  });
  
  test('should navigate to new application form when clicking Add New Application', async ({ page }) => {
    // Navigate to the applications list page
    await page.goto('/applications');
    
    // Click on the Add New Application button
    await page.click('a', { hasText: 'Add New Application' });
    
    // Check that we've navigated to the new application form
    await expect(page).toHaveURL('/applications/new');
    await expect(page).toHaveTitle(/Add New Application/);
  });
  
  test('should filter applications when using search', async ({ page }) => {
    // Navigate to the applications list page
    await page.goto('/applications');
    
    // Fill in the search input
    await page.fill('input[name="q"]', 'test');
    
    // Submit the search form
    await page.click('button[type="submit"]');
    
    // Wait for the HTMX request to complete
    await page.waitForResponse(response => 
      response.url().includes('/htmx/applications/search') && 
      response.status() === 200
    );
    
    // Check that the applications container has been updated
    const applicationsContainer = page.locator('#applications-container');
    await expect(applicationsContainer).toBeVisible();
  });
  
  test('should filter applications by status', async ({ page }) => {
    // Navigate to the applications list page
    await page.goto('/applications');
    
    // Select a status from the dropdown
    await page.selectOption('select[name="status"]', 'applied');
    
    // Wait for the HTMX request to complete
    await page.waitForResponse(response => 
      response.url().includes('/htmx/applications/search') && 
      response.status() === 200
    );
    
    // Check that the applications container has been updated
    const applicationsContainer = page.locator('#applications-container');
    await expect(applicationsContainer).toBeVisible();
  });
});