const testAuthEnabled = ['true', true, '1', 'yes'].includes(Cypress.env('testAuthEnabled'));
const itIfAuth = testAuthEnabled ? it : it.skip;
const itIfNoAuth = testAuthEnabled ? it.skip : it;

describe('auth', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
  });

  itIfNoAuth('should redirect unauthenticated users to Entra ID', () => {
    cy.intercept('GET', '/api/v1/me').as('me');

    cy.request({
      url: '/oauth/login',
      followRedirect: false
    }).then((response) => {
      expect(response.status).to.eq(302);
      expect(response.headers.location).to.contain('/entra/login');
    });

    cy.visit('/');
    cy.wait('@me').its('response.statusCode').should('eq', 401);
    cy.location('pathname').should('eq', '/entra/login');
  });

  itIfAuth('should show the signed-in user name after callback', () => {
    cy.intercept('GET', '/api/v1/me').as('me');

    cy.visit('/oauth/callback');
    cy.location('pathname').should('eq', '/');
    cy.wait('@me').its('response.statusCode').should('eq', 200);
    cy.get('[data-cy="areas-title"]').should('contain', 'Signed in as');
  });
});
