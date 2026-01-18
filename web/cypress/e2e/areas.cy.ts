const testAuthEnabled = ['true', true, '1', 'yes'].includes(Cypress.env('testAuthEnabled'));
const itIfAuth = testAuthEnabled ? it : it.skip;

describe('areas', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
  });

  itIfAuth('should show configured areas', () => {
    cy.intercept('GET', '/api/v1/areas').as('listAreas');

    cy.visit('/oauth/callback');
    cy.wait('@listAreas').its('response.statusCode').should('eq', 200);
    cy.get('[data-cy="areas-list"]').should('exist');
    cy.get('[data-cy="area-item"]').first().should('contain', 'Office 1st Floor');
  });
});
