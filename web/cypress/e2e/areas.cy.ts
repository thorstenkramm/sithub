describe('areas', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
    cy.login();
  });

  it('should show configured areas', () => {
    cy.intercept('GET', '/api/v1/areas').as('listAreas');

    cy.visit('/');

    // Verify Vuetify is properly loaded (catches component import issues)
    cy.waitForVuetify();

    cy.wait('@listAreas').its('response.statusCode').should('eq', 200);
    cy.get('[data-cy="areas-list"]').should('exist');
    // Check that at least one area is displayed (name depends on config)
    cy.get('[data-cy="area-item"]').should('have.length.at.least', 1);

    // Verify area cards render as proper Vuetify cards
    // Note: area-item IS the v-card in the redesigned UI
    cy.get('[data-cy="area-item"]').first().should('have.class', 'v-card');
  });
});
