describe('auth', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
  });

  it('should redirect unauthenticated users to login page', () => {
    cy.visit('/');
    cy.location('pathname').should('eq', '/login');
  });

  it('should show the signed-in user name after login', () => {
    cy.intercept('GET', '/api/v1/me').as('me');
    cy.intercept('GET', '/api/v1/areas').as('listAreas');

    cy.login();
    cy.visit('/');

    cy.wait('@me').its('response.statusCode').should('eq', 200);
    cy.get('[data-cy="user-menu-trigger"]').should('exist');
    cy.wait('@listAreas').its('response.statusCode').should('eq', 200);
    cy.get('[data-cy="areas-list"]').should('exist');
  });

  it('should log out and redirect to login page', () => {
    cy.login();
    cy.visit('/');

    cy.get('[data-cy="user-menu-trigger"]').click();
    cy.get('[data-cy="logout-btn"]').click();

    cy.location('pathname').should('eq', '/login');
  });

  it('should return 401 for unauthenticated API calls', () => {
    cy.request({
      url: '/api/v1/me',
      failOnStatusCode: false
    }).then((response) => {
      expect(response.status).to.eq(401);
      expect(response.headers['content-type']).to.contain('application/vnd.api+json');
      expect(response.body.errors[0].code).to.eq('auth_required');
    });
  });
});
