// Real-backend requests only. Expired keys cannot be created through the API;
// expiry, including the exact expiry instant, is covered by Go security tests.
describe('API key bearer authentication', () => {
  let id = '';
  let key = '';
  const query = 'filter[from]=2026-10-01&filter[to]=2026-10-08';

  beforeEach(() => {
    cy.login();
    cy.request({
      method: 'POST', url: '/api/v1/api-keys', log: false,
      body: { data: { type: 'api-keys', attributes: {
        name: `e2e-bearer-${Date.now()}`, scope: 'read-bookings', expires_in_months: 3
      } } }
    }).then(response => {
      expect(response.status).to.eq(201);
      id = response.body.data.id;
      key = response.body.data.attributes.key;
    });
  });

  afterEach(() => {
    if (!id) return;
    cy.login();
    cy.request({ method: 'DELETE', url: `/api/v1/api-keys/${id}`, failOnStatusCode: false })
      .its('status').should('be.oneOf', [204, 404]);
    cy.then(() => { id = ''; key = ''; });
  });

  function bearerRequest(method: string, url: string, expected: number, token?: string) {
    return cy.then(() => cy.request({
      method, url, headers: { Authorization: `Bearer ${token ?? key}` },
      failOnStatusCode: false, log: false
    })).then(response => {
      expect(response.status).to.eq(expected);
      return response;
    });
  }

  it('should export JSON and CSV without cookies using the same filters', () => {
    cy.clearCookies();
    bearerRequest('GET', `/api/v1/bookings/export?${query}&page[limit]=1&page[offset]=0`, 200)
      .then(response => {
        expect(response.body.data).to.be.an('array');
        expect(response.body.meta.limit).to.eq(1);
        expect(response.body.meta.offset).to.eq(0);
        expect(response.body.links).to.have.property('self');
      });
    bearerRequest('GET', `/api/v1/bookings/export.csv?${query}&filter[area]=test_area`, 200)
      .then(response => {
        expect(response.headers['content-type']).to.contain('text/csv');
        expect(response.body).to.contain('Name,Email,Booked item,Booking date');
      });
    bearerRequest('GET', '/api/v1/bookings/export?filter[from]=bad&filter[to]=2026-10-08', 400);
  });

  it('should deny user and admin routes and reject unknown credentials', () => {
    cy.clearCookies();
    for (const url of ['/api/v1/users', '/api/v1/api-keys', '/api/v1/me']) {
      bearerRequest('GET', url, 403);
    }
    bearerRequest('POST', '/api/v1/bookings', 403);
    bearerRequest('GET', `/api/v1/bookings/export?${query}`, 401, crypto.randomUUID())
      .then(response => {
        expect(response.body.errors).to.be.an('array');
        expect(response.headers['www-authenticate']).to.eq('Bearer error="invalid_token"');
      });
  });

  it('should prefer a bearer credential over an administrator cookie', () => {
    cy.getCookie('sithub_user').should('exist');
    cy.request('/api/v1/users').its('status').should('eq', 200);
    bearerRequest('GET', '/api/v1/users', 403);
    bearerRequest('GET', '/api/v1/users', 401, 'invalid');
  });

  it('should reject a key immediately after revocation', () => {
    cy.then(() => cy.request('DELETE', `/api/v1/api-keys/${id}`)).its('status').should('eq', 204);
    cy.clearCookies();
    bearerRequest('GET', `/api/v1/bookings/export?${query}`, 401);
  });
});
