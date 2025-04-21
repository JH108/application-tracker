const { test, expect } = require('@playwright/test');

test.describe('User Flows', () => {
  test('Complete application lifecycle: create, view, update, delete', async ({ page }) => {
    // 1. Navigate to the applications list page
    await page.goto('/applications');
    
    // 2. Click on the "Add New Application" button
    await page.click('a', { hasText: 'Add New Application' });
    await expect(page).toHaveURL('/applications/new');
    
    // 3. Fill in the new application form
    const companyName = `Flow Test Company ${Date.now()}`;
    await page.fill('input[name="company"]', companyName);
    await page.fill('input[name="position"]', 'Flow Test Position');
    await page.fill('textarea[name="description"]', 'Flow Test Description');
    await page.fill('input[name="url"]', 'https://example.com/flow-test');
    await page.fill('input[name="tags"]', 'flow,test,e2e');
    
    // 4. Submit the form
    await page.click('button[type="submit"]');
    
    // 5. Wait for redirect to applications page
    await page.waitForURL('/applications');
    
    // 6. Find and click on the newly created application
    const applicationLink = page.locator('a', { hasText: companyName }).first();
    await expect(applicationLink).toBeVisible();
    await applicationLink.click();
    
    // 7. Verify the application details
    await expect(page.locator('h1')).toContainText(companyName);
    await expect(page.locator('p.text-xl')).toContainText('Flow Test Position');
    await expect(page.locator('p.whitespace-pre-line')).toContainText('Flow Test Description');
    
    // 8. Update the application status
    await page.click('button', { hasText: 'Mark as In Progress' });
    await page.waitForLoadState('networkidle');
    
    // 9. Verify the status has been updated
    const status = page.locator('.px-3.py-1.bg-yellow-100.text-yellow-800');
    await expect(status).toBeVisible();
    await expect(status).toContainText('In Progress');
    
    // 10. Navigate to edit page
    await page.click('a', { hasText: 'Edit' });
    
    // 11. Update the application
    const updatedCompanyName = `${companyName} - Updated`;
    await page.fill('input[name="company"]', updatedCompanyName);
    await page.click('button[type="submit"]');
    
    // 12. Verify the update
    await expect(page.locator('h1')).toContainText(updatedCompanyName);
    
    // 13. Delete the application
    page.on('dialog', dialog => dialog.accept());
    await page.click('button', { hasText: 'Delete' });
    
    // 14. Verify we're back on the applications list page
    await page.waitForURL('/applications');
    
    // 15. Verify the application is no longer in the list
    await expect(page.locator('body')).toContainText('Applications');
    const deletedApplicationLink = page.locator('a', { hasText: updatedCompanyName });
    await expect(deletedApplicationLink).toHaveCount(0);
  });
  
  test('Search and filter applications', async ({ page }) => {
    // 1. Create a few test applications with different statuses
    await createTestApplication(page, 'Search Test Applied', 'applied');
    await createTestApplication(page, 'Search Test In Progress', 'in_progress');
    await createTestApplication(page, 'Search Test Accepted', 'accepted');
    await createTestApplication(page, 'Search Test Rejected', 'rejected');
    
    // 2. Navigate to the applications list page
    await page.goto('/applications');
    
    // 3. Search for "Search Test"
    await page.fill('input[name="q"]', 'Search Test');
    await page.click('button[type="submit"]');
    
    // 4. Wait for the search results
    await page.waitForResponse(response => 
      response.url().includes('/htmx/applications/search') && 
      response.status() === 200
    );
    
    // 5. Verify all test applications are shown
    await expect(page.locator('a', { hasText: 'Search Test' })).toHaveCount(4);
    
    // 6. Filter by "In Progress" status
    await page.selectOption('select[name="status"]', 'in_progress');
    
    // 7. Wait for the filtered results
    await page.waitForResponse(response => 
      response.url().includes('/htmx/applications/search') && 
      response.status() === 200
    );
    
    // 8. Verify only the "In Progress" application is shown
    await expect(page.locator('a', { hasText: 'Search Test In Progress' })).toHaveCount(1);
    await expect(page.locator('a', { hasText: 'Search Test Applied' })).toHaveCount(0);
    
    // 9. Clear the filters
    await page.fill('input[name="q"]', '');
    await page.selectOption('select[name="status"]', '');
    await page.click('button[type="submit"]');
    
    // 10. Wait for the unfiltered results
    await page.waitForResponse(response => 
      response.url().includes('/htmx/applications/search') && 
      response.status() === 200
    );
  });
});

// Helper function to create a test application with a specific status
async function createTestApplication(page, company, status) {
  // Navigate to the new application form
  await page.goto('/applications/new');
  
  // Fill in the form
  await page.fill('input[name="company"]', company);
  await page.fill('input[name="position"]', 'Test Position');
  await page.fill('textarea[name="description"]', 'Test Description');
  await page.fill('input[name="url"]', 'https://example.com');
  await page.fill('input[name="tags"]', 'test,search');
  
  // Submit the form
  await page.click('button[type="submit"]');
  
  // Wait for redirect to applications page
  await page.waitForURL('/applications');
  
  // Find the application and navigate to its detail page
  const applicationLink = page.locator('a', { hasText: company }).first();
  await applicationLink.click();
  
  // Wait for the detail page to load
  await page.waitForSelector('h1:has-text("' + company + '")');
  
  // Update the status if needed
  if (status !== 'applied') {
    let buttonText;
    switch (status) {
      case 'in_progress':
        buttonText = 'Mark as In Progress';
        break;
      case 'accepted':
        buttonText = 'Mark as Accepted';
        break;
      case 'rejected':
        buttonText = 'Mark as Rejected';
        break;
    }
    
    await page.click('button', { hasText: buttonText });
    await page.waitForLoadState('networkidle');
  }
  
  // Go back to the applications list
  await page.click('a', { hasText: 'Back to Applications' });
}