const testAuthEnabled = ['true', true, '1', 'yes'].includes(Cypress.env('testAuthEnabled'));
const itIfAuth = testAuthEnabled ? it : it.skip;

describe('areas', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
  });

  itIfAuth('should show an empty state when no areas exist', () => {
    cy.intercept('GET', '/api/v1/areas').as('listAreas');

    cy.visit('/oauth/callback');
    cy.wait('@listAreas').its('response.statusCode').should('eq', 200);
    cy.get('[data-cy="areas-empty"]').should('contain', 'No areas available.');
  });
});
