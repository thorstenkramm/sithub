const testAuthEnabled = ['true', true, '1', 'yes'].includes(Cypress.env('testAuthEnabled'));
const testAuthPermitted = ['true', true, '1', 'yes'].includes(Cypress.env('testAuthPermitted'));
const itIfAuth = testAuthEnabled ? it : it.skip;
const itIfPermitted = testAuthEnabled && testAuthPermitted ? it : it.skip;
const itIfNoAuth = testAuthEnabled ? it.skip : it;
const itIfForbidden = testAuthEnabled && !testAuthPermitted ? it : it.skip;

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

    cy.request({
      url: '/api/v1/me',
      failOnStatusCode: false
    }).then((response) => {
      expect(response.status).to.eq(401);
      expect(response.headers['content-type']).to.contain('application/vnd.api+json');
      expect(response.body.errors[0].code).to.eq('auth_required');
    });
  });

  itIfPermitted('should show the signed-in user name after callback', () => {
    cy.intercept('GET', '/api/v1/me').as('me');

    cy.visit('/oauth/callback');
    cy.location('pathname').should('eq', '/');
    cy.wait('@me').its('response.statusCode').should('eq', 200);
    cy.get('[data-cy="areas-title"]').should('contain', 'Signed in as');
  });

  itIfForbidden('should show access denied for forbidden users', () => {
    cy.visit('/oauth/callback');
    cy.location('pathname').should('eq', '/access-denied');
    cy.get('[data-cy="access-denied-title"]').should('contain', 'Access denied');

    cy.request({
      url: '/api/v1/me',
      failOnStatusCode: false
    }).then((response) => {
      expect(response.status).to.eq(403);
      expect(response.headers['content-type']).to.contain('application/vnd.api+json');
      expect(response.body.errors[0].code).to.eq('forbidden');
    });
  });
});
