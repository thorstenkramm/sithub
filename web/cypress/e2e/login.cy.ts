describe('login', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
  });

  it('should show login form for unauthenticated users', () => {
    cy.visit('/');
    cy.location('pathname').should('eq', '/login');
    cy.get('[data-cy="login-form"]').should('exist');
    cy.get('[data-cy="login-email"]').should('exist');
    cy.get('[data-cy="login-password"]').should('exist');
    cy.get('[data-cy="login-submit"]').should('exist');
    cy.get('[data-cy="login-entraid"]').should('exist');
  });

  it('should login with valid credentials and redirect to home', () => {
    cy.intercept('POST', '/api/v1/auth/login').as('loginRequest');
    cy.intercept('GET', '/api/v1/me').as('me');
    cy.intercept('GET', '/api/v1/areas').as('listAreas');

    cy.visit('/login');

    cy.get('[data-cy="login-email"]').find('input').type(Cypress.env('testUserEmail'));
    cy.get('[data-cy="login-password"]').find('input').type(Cypress.env('testUserPassword'));
    cy.get('[data-cy="login-submit"]').click();

    cy.wait('@loginRequest').its('response.statusCode').should('eq', 200);
    cy.location('pathname').should('eq', '/');
    cy.wait('@me').its('response.statusCode').should('eq', 200);
    cy.get('[data-cy="user-menu-trigger"]').should('exist');
  });

  it('should show error for invalid credentials', () => {
    cy.intercept('POST', '/api/v1/auth/login').as('loginRequest');

    cy.visit('/login');

    cy.get('[data-cy="login-email"]').find('input').type('wrong@example.com');
    cy.get('[data-cy="login-password"]').find('input').type('wrong-password-here');
    cy.get('[data-cy="login-submit"]').click();

    cy.wait('@loginRequest').its('response.statusCode').should('eq', 401);
    cy.get('[data-cy="login-error"]').should('contain', 'Invalid email or password');
    cy.location('pathname').should('eq', '/login');
  });

  it('should allow authenticated users to access protected pages', () => {
    cy.intercept('GET', '/api/v1/areas').as('listAreas');

    cy.login();
    cy.visit('/');

    // Should stay on the home page (not redirect to login)
    cy.location('pathname').should('eq', '/');
    cy.wait('@listAreas').its('response.statusCode').should('eq', 200);
  });
});
