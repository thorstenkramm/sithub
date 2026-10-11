describe('API key management', () => {
  const createdIds = new Set<string>();

  beforeEach(() => {
    cy.clearLocalStorage();
    cy.login();
    cy.intercept('GET', '/api/v1/api-keys').as('listApiKeys');
    cy.intercept('POST', '/api/v1/api-keys').as('createApiKey');
    cy.intercept('DELETE', '/api/v1/api-keys/*').as('deleteApiKey');
  });

  afterEach(() => {
    if (!createdIds.size) return;
    cy.login();
    for (const id of createdIds) {
      cy.request({ method: 'DELETE', url: `/api/v1/api-keys/${id}`, failOnStatusCode: false })
        .its('status').should('be.oneOf', [204, 404]);
    }
    cy.then(() => createdIds.clear());
  });

  it('should create a key, show it once, list its prefix and revoke after confirmation', () => {
    const name = `e2e-key-${Date.now()}`;
    cy.visit('/');
    cy.get('[data-cy="user-menu-trigger"]').click();
    cy.get('[data-cy="api-keys-btn"]').click();
    cy.location('pathname').should('eq', '/admin/api-keys');
    cy.wait('@listApiKeys').its('response.statusCode').should('eq', 200);
    cy.get('[data-cy="api-keys-add"]').click();
    cy.get('[data-cy="api-key-submit"]').should('be.disabled');
    cy.get('[data-cy="api-key-name-input"] input').type(name);
    cy.get('[data-cy="api-key-expiry-select"]').should('contain.text', '12 months');
    cy.get('[data-cy="api-key-scope-select"]').should('contain.text', 'Read bookings');
    cy.get('[data-cy="api-key-submit"]').click();
    cy.wait('@createApiKey').then(({ request, response }) => {
      expect(response?.statusCode).to.eq(201);
      const resource = response!.body.data;
      createdIds.add(resource.id);
      const key: string = resource.attributes.key;
      expect(key).to.match(/^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/);
      expect(resource.attributes).not.to.have.property('key_hash');
      expect(request.body.data.attributes).to.deep.eq({ name, scope: 'read-bookings', expires_in_months: 12 });
      cy.get('[data-cy="api-key-value"] input').should('have.value', key);
      cy.get('[data-cy="api-key-once-warning"]').should('be.visible').and('contain.text', 'cannot be shown again');
      cy.get('[data-cy="api-key-copy"]').should('be.visible');
      cy.get('[data-cy="api-key-done"]').click();
      cy.wait('@listApiKeys').then(({ response: list }) => {
        expect(JSON.stringify(list!.body)).not.to.contain(key);
        expect(JSON.stringify(list!.body)).not.to.contain('key_hash');
        for (const entry of list!.body.data) expect(entry.attributes).not.to.have.property('key');
      });
      cy.get('[data-cy="api-key-value"]').should('not.exist');
      cy.get('body').should('not.contain.text', key);
      cy.contains('[data-cy="api-key-row"]', name).within(() => {
        cy.get('[data-cy="api-key-prefix"]').should('have.text', `${key.slice(0, 10)}…`);
        cy.get('[data-cy="api-key-scope"]').should('have.text', 'Read bookings');
        cy.get('[data-cy="api-key-expiry"]').should('not.be.empty');
        cy.get('[data-cy="api-key-revoke"]').click();
      });
      cy.get('[data-cy="confirm-dialog"]').should('contain.text', name);
      cy.get('[data-cy="confirm-dialog-cancel"]').click();
      cy.contains('[data-cy="api-key-row"]', name).find('[data-cy="api-key-revoke"]').click();
      cy.get('[data-cy="confirm-dialog-confirm"]').click();
      cy.wait('@deleteApiKey').its('response.statusCode').should('eq', 204);
      cy.wait('@listApiKeys');
      cy.contains('[data-cy="api-key-row"]', name).should('not.exist');
      cy.get('[data-cy="api-keys-snackbar"]').should('contain.text', 'API key revoked');
    });
  });

  it('should open API Keys from the mobile drawer for admins', () => {
    cy.viewport('iphone-x');
    cy.intercept('GET', '/api/v1/areas').as('mobileAreas');
    cy.visit('/');
    cy.wait('@mobileAreas');
    cy.get('[data-cy="mobile-menu-btn"]').click();
    cy.get('[data-cy="mobile-api-keys-btn"]').scrollIntoView().should('be.visible').click();
    cy.location('pathname').should('eq', '/admin/api-keys');
    cy.wait('@listApiKeys');
  });

  it('should hide entries, deny the page and reject all endpoints for non-admins', () => {
    cy.login('alex@sithub.local', 'SitHubDemo2026!!');
    cy.visit('/');
    cy.get('[data-cy="user-menu-trigger"]').click();
    cy.get('[data-cy="api-keys-btn"]').should('not.exist');
    cy.get('[data-cy="mobile-api-keys-btn"]').should('not.exist');
    cy.visit('/admin/api-keys');
    cy.location('pathname').should('eq', '/access-denied');
    for (const method of ['GET', 'POST', 'DELETE']) {
      cy.request({ method, url: `/api/v1/api-keys${method === 'DELETE' ? '/missing' : ''}`, failOnStatusCode: false })
        .its('status').should('eq', 403);
    }
  });

  it('should require a session on all endpoints', () => {
    cy.clearCookies();
    for (const method of ['GET', 'POST', 'DELETE']) {
      cy.request({ method, url: `/api/v1/api-keys${method === 'DELETE' ? '/missing' : ''}`, failOnStatusCode: false })
        .its('status').should('eq', 401);
    }
  });
});
