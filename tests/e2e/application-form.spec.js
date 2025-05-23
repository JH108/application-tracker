const { test, expect } = require('@playwright/test');

test.describe('Application Form', () => {
  test('should load the new application form', async ({ page }) => {
    // Navigate to the new application form
    await page.goto('/applications/new');

    // Check that the page title is correct
    await expect(page).toHaveTitle(/Add New Application/);

    // Check that the form is present
    const form = page.locator('form[hx-post="/api/applications"]');
    await expect(form).toBeVisible();

    // Check that all form fields are present
    await expect(page.locator('input[name="company"]')).toBeVisible();
    await expect(page.locator('input[name="position"]')).toBeVisible();
    await expect(page.locator('textarea[name="description"]')).toBeVisible();
    await expect(page.locator('input[name="url"]')).toBeVisible();
    await expect(page.locator('input[name="tags"]')).toBeVisible();

    // Check that the submit button is present
    const submitButton = page.locator('button[type="submit"]');
    await expect(submitButton).toBeVisible();
    await expect(submitButton).toContainText('Save Application');

    // Check that the cancel button is present
    const cancelButton = page.locator('a', { hasText: 'Cancel' });
    await expect(cancelButton).toBeVisible();
  });

  test('should create a new application when submitting the form', async ({ page }) => {
    // Navigate to the new application form
    await page.goto('/applications/new');

    // Fill in the form with a unique company name
    const timestamp = Date.now();
    const companyName = `Form Test Company ${timestamp}`;
    await page.fill('input[name="company"]', companyName);
    await page.fill('input[name="position"]', 'Form Test Position');
    await page.fill('textarea[name="description"]', 'Form Test Description');
    await page.fill('input[name="url"]', 'https://example.com/form-test');
    await page.fill('input[name="tags"]', 'form,test');

    // Submit the form
    await page.click('button[type="submit"]');

    // Wait for redirect to applications page
    await page.waitForURL('/applications');

    // Wait for the applications to load
    await page.waitForResponse(response => 
      response.url().includes('/htmx/applications') && 
      response.status() === 200
    );

    // Check that we're back on the applications page
    await expect(page).toHaveTitle(/Applications/);

    // Try to find the new application in the list
    let applicationLink = page.locator(`a:has-text("${companyName}")`).first();
    let isVisible = await applicationLink.isVisible();

    // If not found, try searching for it
    if (!isVisible) {
      // Fill in the search input with the company name
      await page.fill('input[name="q"]', companyName);

      // Submit the search form
      await page.click('button[type="submit"]');

      // Wait for the search results
      await page.waitForResponse(response => 
        response.url().includes('/htmx/applications/search') && 
        response.status() === 200
      );

      // Try to find the application again
      applicationLink = page.locator(`a:has-text("${companyName}")`).first();
      await expect(applicationLink).toBeVisible();
    }
  });

  test('should validate required fields', async ({ page }) => {
    // Navigate to the new application form
    await page.goto('/applications/new');

    // Submit the form without filling in required fields
    await page.click('button[type="submit"]');

    // Check that we're still on the form page
    await expect(page).toHaveURL('/applications/new');

    // Check for validation messages
    // Note: The exact validation behavior depends on the application implementation
    // This test assumes HTML5 validation is being used
    const companyInput = page.locator('input[name="company"]');
    const positionInput = page.locator('input[name="position"]');

    // Check if the inputs have the :invalid pseudo-class
    const isCompanyInvalid = await companyInput.evaluate(el => el.validity.valid === false);
    const isPositionInvalid = await positionInput.evaluate(el => el.validity.valid === false);

    expect(isCompanyInvalid || isPositionInvalid).toBeTruthy();
  });

  test('should navigate back to applications page when clicking Cancel', async ({ page }) => {
    // Navigate to the new application form
    await page.goto('/applications/new');

    // Click on the Cancel button
    await page.click('a', { hasText: 'Cancel' });

    // Check that we're back on the applications page
    await expect(page).toHaveURL('/applications');
    await expect(page).toHaveTitle(/Applications/);
  });

  test('should load the edit form with application data', async ({ page }) => {
    // First create an application with a unique company name
    await page.goto('/applications/new');
    const timestamp = Date.now();
    const companyName = `Edit Test Company ${timestamp}`;
    await page.fill('input[name="company"]', companyName);
    await page.fill('input[name="position"]', 'Edit Test Position');
    await page.fill('textarea[name="description"]', 'Edit Test Description');
    await page.fill('input[name="url"]', 'https://example.com/edit-test');
    await page.fill('input[name="tags"]', 'edit,test');
    await page.click('button[type="submit"]');

    // Wait for redirect to applications page
    await page.waitForURL('/applications');

    // Wait for the applications to load
    await page.waitForResponse(response => 
      response.url().includes('/htmx/applications') && 
      response.status() === 200
    );

    // Try to find the new application in the list
    let applicationLink = page.locator(`a:has-text("${companyName}")`).first();
    let isVisible = await applicationLink.isVisible();

    // If not found, try searching for it
    if (!isVisible) {
      // Fill in the search input with the company name
      await page.fill('input[name="q"]', companyName);

      // Submit the search form
      await page.click('button[type="submit"]');

      // Wait for the search results
      await page.waitForResponse(response => 
        response.url().includes('/htmx/applications/search') && 
        response.status() === 200
      );

      // Try to find the application again
      applicationLink = page.locator(`a:has-text("${companyName}")`).first();
      isVisible = await applicationLink.isVisible();

      if (!isVisible) {
        throw new Error(`Could not find application with company name: ${companyName}`);
      }
    }

    // Navigate to the detail page
    await applicationLink.click();

    // Wait for the detail page to load
    await page.waitForSelector(`h1:has-text("${companyName}")`);

    // Click on the Edit button
    await page.click('a', { hasText: 'Edit' });

    // Check that the form is pre-filled with the application data
    await expect(page.locator('input[name="company"]')).toHaveValue(companyName);
    await expect(page.locator('input[name="position"]')).toHaveValue('Edit Test Position');
    await expect(page.locator('textarea[name="description"]')).toHaveValue('Edit Test Description');
    await expect(page.locator('input[name="url"]')).toHaveValue('https://example.com/edit-test');
    await expect(page.locator('input[name="tags"]')).toHaveValue('edit,test');
  });
});
