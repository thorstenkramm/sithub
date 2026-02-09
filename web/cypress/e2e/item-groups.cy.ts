describe('item groups', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
    cy.login();
  });

  it('should show item groups for selected area', () => {
    cy.intercept('GET', '/api/v1/areas').as('listAreas');
    cy.intercept('GET', '/api/v1/areas/*/item-groups').as('listItemGroups');

    cy.visit('/');

    // Wait for areas to load and click the first one
    cy.wait('@listAreas').its('response.statusCode').should('eq', 200);
    cy.get('[data-cy="area-item"]').first().click();

    cy.wait('@listItemGroups').its('response.statusCode').should('eq', 200);
    cy.location('pathname').should('match', /\/areas\/.*\/item-groups/);
    cy.get('[data-cy="item-groups-list"]').should('exist');
    // Check that at least one item group exists (name depends on config)
    cy.get('[data-cy="item-group-item"]').should('have.length.at.least', 1);
  });
});
