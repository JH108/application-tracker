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

    // Check that the applications list is present
    const applicationsList = page.locator('#applications-list');
    await expect(applicationsList).toBeVisible();

    // Check that the pagination controls are present
    const paginationControls = page.locator('#pagination-controls');
    await expect(paginationControls).toBeVisible();
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

    // Check that the applications list has been updated
    const applicationsList = page.locator('#applications-list');
    await expect(applicationsList).toBeVisible();
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

    // Check that the applications list has been updated
    const applicationsList = page.locator('#applications-list');
    await expect(applicationsList).toBeVisible();
  });

  test('should change page size when selecting from dropdown', async ({ page }) => {
    // Navigate to the applications list page
    await page.goto('/applications');

    // Select a different page size
    await page.selectOption('#page-size-selector', '25');

    // Wait for the HTMX request to complete
    await page.waitForResponse(response => 
      response.url().includes('/htmx/applications') && 
      response.status() === 200
    );

    // Check that the page size has been updated in the pagination data
    const paginationData = page.locator('#pagination-data');
    await expect(paginationData).toHaveAttribute('data-page-size', '25');
  });

  test('should navigate to next page when clicking Next button', async ({ page }) => {
    // Navigate to the applications list page
    await page.goto('/applications');

    // Wait for the initial load to complete
    await page.waitForResponse(response => 
      response.url().includes('/htmx/applications') && 
      response.status() === 200
    );

    // Check if there's a next page before trying to click the button
    const nextButton = page.locator('#next-page');
    const isDisabled = await nextButton.getAttribute('disabled');

    if (isDisabled === null) {
      // Click the Next button if it's not disabled
      await nextButton.click();

      // Wait for the HTMX request to complete
      await page.waitForResponse(response => 
        response.url().includes('/htmx/applications') && 
        response.status() === 200
      );

      // Check that the current page has been updated in the pagination data
      const paginationData = page.locator('#pagination-data');
      const currentPage = await paginationData.getAttribute('data-current-page');
      expect(parseInt(currentPage)).toBeGreaterThan(1);
    }
  });
});
