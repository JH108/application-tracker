const { test, expect } = require('@playwright/test');

// Helper function to create a test application and return its data
async function createTestApplication(request) {
  const data = {
    company: 'API Test Company',
    position: 'API Test Position',
    description: 'API Test Description',
    url: 'https://example.com/api-test',
    tags: ['api', 'test']
  };
  
  const response = await request.post('/api/applications', {
    data: data,
    headers: {
      'Content-Type': 'application/json'
    }
  });
  
  expect(response.ok()).toBeTruthy();
  const responseData = await response.json();
  return responseData.data;
}

test.describe('API Endpoints', () => {
  test.describe('GET Endpoints', () => {
    test('GET /api/health should return status ok', async ({ request }) => {
      const response = await request.get('/api/health');
      
      expect(response.ok()).toBeTruthy();
      const data = await response.json();
      expect(data.status).toBe('ok');
    });
    
    test('GET /api/applications should return all applications', async ({ request }) => {
      const response = await request.get('/api/applications');
      
      expect(response.ok()).toBeTruthy();
      const data = await response.json();
      expect(data.success).toBeTruthy();
      expect(Array.isArray(data.data)).toBeTruthy();
    });
    
    test('GET /api/applications/:id should return a specific application', async ({ request }) => {
      // First create a test application
      const application = await createTestApplication(request);
      
      // Then get it by ID
      const response = await request.get(`/api/applications/${application.id}`);
      
      expect(response.ok()).toBeTruthy();
      const data = await response.json();
      expect(data.success).toBeTruthy();
      expect(data.data.id).toBe(application.id);
      expect(data.data.company).toBe('API Test Company');
    });
    
    test('GET /api/applications/search should filter applications', async ({ request }) => {
      // First create a test application
      await createTestApplication(request);
      
      // Then search for it
      const response = await request.get('/api/applications/search?q=API Test');
      
      expect(response.ok()).toBeTruthy();
      const data = await response.json();
      expect(data.success).toBeTruthy();
      expect(Array.isArray(data.data)).toBeTruthy();
      
      // Check if at least one result contains our test company
      const hasTestCompany = data.data.some(app => app.company === 'API Test Company');
      expect(hasTestCompany).toBeTruthy();
    });
  });
  
  test.describe('POST Endpoints', () => {
    test('POST /api/applications should create a new application', async ({ request }) => {
      const data = {
        company: 'POST Test Company',
        position: 'POST Test Position',
        description: 'POST Test Description',
        url: 'https://example.com/post-test',
        tags: ['post', 'test']
      };
      
      const response = await request.post('/api/applications', {
        data: data,
        headers: {
          'Content-Type': 'application/json'
        }
      });
      
      expect(response.ok()).toBeTruthy();
      const responseData = await response.json();
      expect(responseData.success).toBeTruthy();
      expect(responseData.data.company).toBe('POST Test Company');
      expect(responseData.data.position).toBe('POST Test Position');
      expect(responseData.data.id).toBeTruthy();
    });
    
    test('POST /api/applications should validate required fields', async ({ request }) => {
      const data = {
        // Missing required fields: company and position
        description: 'Invalid Test Description',
        url: 'https://example.com/invalid-test',
        tags: ['invalid', 'test']
      };
      
      const response = await request.post('/api/applications', {
        data: data,
        headers: {
          'Content-Type': 'application/json'
        }
      });
      
      expect(response.ok()).toBeFalsy();
      const responseData = await response.json();
      expect(responseData.success).toBeFalsy();
      expect(responseData.message).toContain('required');
    });
  });
  
  test.describe('PUT Endpoints', () => {
    test('PUT /api/applications/:id should update an application', async ({ request }) => {
      // First create a test application
      const application = await createTestApplication(request);
      
      // Then update it
      const updateData = {
        company: 'Updated Company',
        position: 'Updated Position'
      };
      
      const response = await request.put(`/api/applications/${application.id}`, {
        data: updateData,
        headers: {
          'Content-Type': 'application/json'
        }
      });
      
      expect(response.ok()).toBeTruthy();
      const responseData = await response.json();
      expect(responseData.success).toBeTruthy();
      expect(responseData.data.company).toBe('Updated Company');
      expect(responseData.data.position).toBe('Updated Position');
      expect(responseData.data.id).toBe(application.id);
    });
    
    test('PUT /api/applications/:id/status should update application status', async ({ request }) => {
      // First create a test application
      const application = await createTestApplication(request);
      
      // Then update its status
      const updateData = {
        status: 'in_progress'
      };
      
      const response = await request.put(`/api/applications/${application.id}/status`, {
        data: updateData,
        headers: {
          'Content-Type': 'application/json'
        }
      });
      
      expect(response.ok()).toBeTruthy();
      const responseData = await response.json();
      expect(responseData.success).toBeTruthy();
      expect(responseData.data.status).toBe('in_progress');
      expect(responseData.data.id).toBe(application.id);
    });
  });
  
  test.describe('DELETE Endpoints', () => {
    test('DELETE /api/applications/:id should delete an application', async ({ request }) => {
      // First create a test application
      const application = await createTestApplication(request);
      
      // Then delete it
      const response = await request.delete(`/api/applications/${application.id}`);
      
      expect(response.ok()).toBeTruthy();
      const responseData = await response.json();
      expect(responseData.success).toBeTruthy();
      
      // Verify it's deleted by trying to get it
      const getResponse = await request.get(`/api/applications/${application.id}`);
      expect(getResponse.status()).toBe(404);
    });
  });
});