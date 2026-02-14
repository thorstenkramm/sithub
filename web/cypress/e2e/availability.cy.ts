import {
  openFirstAreaAndWaitItemGroups,
  resetAndLogin,
  setupItemsPageIntercepts,
  visitHomeAndWaitForAreas
} from '../support/flows';

describe('weekly availability preview', () => {
  beforeEach(() => {
    resetAndLogin();
    setupItemsPageIntercepts();
  });

  it('should show week selector on item groups page', () => {
    visitHomeAndWaitForAreas();
    openFirstAreaAndWaitItemGroups();
    cy.wait('@availability');

    cy.get('[data-cy="week-selector-card"]').should('exist');
    cy.get('[data-cy="week-selector"]').should('exist');
  });

  it('should show availability indicators with weekday labels', () => {
    visitHomeAndWaitForAreas();
    openFirstAreaAndWaitItemGroups();
    cy.wait('@availability');

    cy.get('[data-cy="availability-indicators"]').should('have.length.at.least', 1);
    cy.get('[data-cy="availability-indicators"]').first().within(() => {
      cy.get('.indicator-dot').should('have.length', 5);
      cy.get('[data-cy-weekday="MO"]').should('exist');
      cy.get('[data-cy-weekday="TU"]').should('exist');
      cy.get('[data-cy-weekday="WE"]').should('exist');
      cy.get('[data-cy-weekday="TH"]').should('exist');
      cy.get('[data-cy-weekday="FR"]').should('exist');
    });
  });

  it('should reload availability when week changes', () => {
    visitHomeAndWaitForAreas();
    openFirstAreaAndWaitItemGroups();
    cy.wait('@availability');

    // Open the week selector dropdown
    cy.get('[data-cy="week-selector"]').click();
    // Pick the second option in the dropdown overlay
    cy.get('.v-overlay__content .v-list-item').eq(1).click();
    cy.wait('@availability');
  });
});
