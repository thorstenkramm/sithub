import {
  openFirstAreaAndWaitItemGroups,
  resetAndLogin,
  setupItemsPageIntercepts,
  visitHomeAndWaitForAreas
} from '../support/flows';

describe('item groups', () => {
  beforeEach(() => {
    resetAndLogin();
    setupItemsPageIntercepts();
  });

  it('should show item groups for selected area', () => {
    visitHomeAndWaitForAreas();
    openFirstAreaAndWaitItemGroups();
    cy.get('[data-cy="item-groups-list"]').should('exist');
    // Check that at least one item group exists (name depends on config)
    cy.get('[data-cy="item-group-item"]').should('have.length.at.least', 1);
  });
});
