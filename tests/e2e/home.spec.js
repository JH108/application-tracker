const { test, expect } = require('@playwright/test');

test.describe('Home Page', () => {
  test('should load the home page', async ({ page }) => {
    // Navigate to the home page
    await page.goto('/');

    // Check that the page title is correct
    await expect(page).toHaveTitle(/Home/);

    // Check that the main heading is present
    const heading = page.locator('h1');
    await expect(heading).toBeVisible();

    // Check that the navigation links are present
    const navLinks = page.locator('nav a');
    await expect(navLinks).toHaveCount(2); // Home and Applications links

    // Check that the Applications link is present
    const applicationsLink = page.locator('nav a', { hasText: 'Applications' });
    await expect(applicationsLink).toBeVisible();

    // Check that the stats section is present
    const statsSection = page.locator('div.grid.grid-cols-1.md\\:grid-cols-4');
    await expect(statsSection).toBeVisible();

    // Check that the recent applications section is present
    const recentApplications = page.locator('#recent-applications');
    await expect(recentApplications).toBeVisible();
  });

  test('should navigate to applications page when clicking Applications link', async ({ page }) => {
    // Navigate to the home page
    await page.goto('/');

    // Click on the Applications link
    await page.click('nav a', { hasText: 'Applications' });

    // Check that we've navigated to the applications page
    await expect(page).toHaveURL('/applications');
    await expect(page).toHaveTitle(/Applications/);
  });

  test('should filter recent applications by status', async ({ page }) => {
    // Navigate to the home page
    await page.goto('/');

    // Wait for the initial load to complete
    await page.waitForResponse(response => 
      response.url().includes('/htmx/applications') && 
      response.status() === 200
    );

    // Select a status from the dropdown
    await page.selectOption('select', 'applied');

    // Wait for the HTMX request to complete
    await page.waitForResponse(response => 
      response.url().includes('/htmx/applications') && 
      response.status() === 200
    );

    // Check that the recent applications section has been updated
    const recentApplications = page.locator('#recent-applications');
    await expect(recentApplications).toBeVisible();
  });
});
