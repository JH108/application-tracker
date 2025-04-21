const { test, expect } = require('@playwright/test');

// Helper function to create a test application and return its ID
async function createTestApplication(page) {
  // Navigate to the new application form
  await page.goto('/applications/new');
  
  // Fill in the form
  await page.fill('input[name="company"]', 'Test Company');
  await page.fill('input[name="position"]', 'Test Position');
  await page.fill('textarea[name="description"]', 'Test Description');
  await page.fill('input[name="url"]', 'https://example.com');
  await page.fill('input[name="tags"]', 'test,playwright');
  
  // Submit the form
  await page.click('button[type="submit"]');
  
  // Wait for redirect to applications page
  await page.waitForURL('/applications');
  
  // Get the ID of the newly created application
  // We'll need to find the application in the list and extract its ID
  const applicationLink = page.locator('a', { hasText: 'Test Company' }).first();
  const href = await applicationLink.getAttribute('href');
  const id = href.split('/').pop();
  
  return id;
}

test.describe('Application Detail Page', () => {
  let applicationId;
  
  test.beforeEach(async ({ page }) => {
    // Create a test application before each test
    applicationId = await createTestApplication(page);
  });
  
  test('should load the application detail page', async ({ page }) => {
    // Navigate to the application detail page
    await page.goto(`/applications/${applicationId}`);
    
    // Check that the page title contains the company and position
    await expect(page).toHaveTitle(/Test Company.*Test Position/);
    
    // Check that the company name is displayed
    const companyHeading = page.locator('h1');
    await expect(companyHeading).toBeVisible();
    await expect(companyHeading).toContainText('Test Company');
    
    // Check that the position is displayed
    const positionText = page.locator('p.text-xl');
    await expect(positionText).toBeVisible();
    await expect(positionText).toContainText('Test Position');
    
    // Check that the description is displayed
    const description = page.locator('p.whitespace-pre-line');
    await expect(description).toBeVisible();
    await expect(description).toContainText('Test Description');
    
    // Check that the URL is displayed
    const url = page.locator('a[href="https://example.com"]');
    await expect(url).toBeVisible();
    await expect(url).toContainText('https://example.com');
    
    // Check that the tags are displayed
    const tags = page.locator('.flex.flex-wrap.gap-2 span');
    await expect(tags).toHaveCount(2);
    await expect(tags.nth(0)).toContainText('test');
    await expect(tags.nth(1)).toContainText('playwright');
    
    // Check that the status buttons are present
    const inProgressButton = page.locator('button', { hasText: 'Mark as In Progress' });
    const acceptedButton = page.locator('button', { hasText: 'Mark as Accepted' });
    const rejectedButton = page.locator('button', { hasText: 'Mark as Rejected' });
    
    await expect(inProgressButton).toBeVisible();
    await expect(acceptedButton).toBeVisible();
    await expect(rejectedButton).toBeVisible();
  });
  
  test('should navigate to edit page when clicking Edit button', async ({ page }) => {
    // Navigate to the application detail page
    await page.goto(`/applications/${applicationId}`);
    
    // Click on the Edit button
    await page.click('a', { hasText: 'Edit' });
    
    // Check that we've navigated to the edit page
    await expect(page).toHaveURL(`/applications/${applicationId}/edit`);
  });
  
  test('should update status when clicking status buttons', async ({ page }) => {
    // Navigate to the application detail page
    await page.goto(`/applications/${applicationId}`);
    
    // Click on the "Mark as In Progress" button
    await page.click('button', { hasText: 'Mark as In Progress' });
    
    // Wait for the page to reload
    await page.waitForLoadState('networkidle');
    
    // Check that the status has been updated
    const status = page.locator('.px-3.py-1.bg-yellow-100.text-yellow-800');
    await expect(status).toBeVisible();
    await expect(status).toContainText('In Progress');
  });
  
  test('should delete application when clicking Delete button', async ({ page }) => {
    // Navigate to the application detail page
    await page.goto(`/applications/${applicationId}`);
    
    // Set up a dialog handler to accept the confirmation dialog
    page.on('dialog', dialog => dialog.accept());
    
    // Click on the Delete button
    await page.click('button', { hasText: 'Delete' });
    
    // Wait for redirect to applications page
    await page.waitForURL('/applications');
    
    // Check that we're back on the applications page
    await expect(page).toHaveTitle(/Applications/);
  });
});