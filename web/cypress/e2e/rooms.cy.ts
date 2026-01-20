const testAuthEnabled = ['true', true, '1', 'yes'].includes(Cypress.env('testAuthEnabled'));
const itIfAuth = testAuthEnabled ? it : it.skip;

describe('rooms', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
  });

  itIfAuth('should show rooms for selected area', () => {
    cy.intercept('GET', '/api/v1/areas').as('listAreas');
    cy.intercept('GET', '/api/v1/areas/*/rooms').as('listRooms');

    cy.visit('/oauth/callback');

    // Wait for areas to load and click the first one
    cy.wait('@listAreas').its('response.statusCode').should('eq', 200);
    cy.get('[data-cy="area-item"]').first().click();

    cy.wait('@listRooms').its('response.statusCode').should('eq', 200);
    cy.location('pathname').should('match', /\/areas\/.*\/rooms/);
    cy.get('[data-cy="rooms-list"]').should('exist');
    // Check that at least one room exists (name depends on config)
    cy.get('[data-cy="room-item"]').should('have.length.at.least', 1);
  });
});
