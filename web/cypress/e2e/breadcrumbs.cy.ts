import {
  openFirstAreaAndWaitItemGroups,
  openFirstItemGroupAndWaitItems,
  resetAndLogin,
  setupItemsPageIntercepts,
  visitHomeAndWaitForAreas
} from '../support/flows';

describe('breadcrumbs', () => {
  beforeEach(() => {
    resetAndLogin();
    setupItemsPageIntercepts();
  });

  it('should navigate from items back to item groups via area breadcrumb', () => {
    visitHomeAndWaitForAreas();
    openFirstAreaAndWaitItemGroups();

    // Store the area URL for comparison
    cy.location('pathname').then((areaPath) => {
      openFirstItemGroupAndWaitItems();

      // areaId should be in query params
      cy.location('search').should('include', 'areaId=');

      // Click the area breadcrumb (second breadcrumb link)
      cy.get('[data-cy="breadcrumb-item-1"]').find('a').click();

      // Should navigate back to item groups for that area
      cy.location('pathname').should('eq', areaPath);
    });
  });

  it('should navigate from items back to areas via Home breadcrumb', () => {
    visitHomeAndWaitForAreas();
    openFirstAreaAndWaitItemGroups();
    openFirstItemGroupAndWaitItems();

    // Click Home breadcrumb
    cy.get('[data-cy="breadcrumb-item-0"]').find('a').click();

    // Should navigate back to areas list
    cy.location('pathname').should('eq', '/');
  });

  it('should navigate from item groups back to areas via Home breadcrumb', () => {
    visitHomeAndWaitForAreas();
    openFirstAreaAndWaitItemGroups();

    // Click Home breadcrumb
    cy.get('[data-cy="breadcrumb-item-0"]').find('a').click();

    // Should navigate back to areas list
    cy.location('pathname').should('eq', '/');
  });

  it('should preserve correct areaId across different areas', () => {
    visitHomeAndWaitForAreas();

    // Get the list of areas and check if we have at least 2
    cy.get('[data-cy="area-item"]').then(($areas) => {
      if ($areas.length < 2) {
        // Only one area available - skip cross-area test
        return;
      }

      // Navigate to first area
      openFirstAreaAndWaitItemGroups();

      cy.location('pathname').then((firstAreaPath) => {
        openFirstItemGroupAndWaitItems();

        // Click area breadcrumb - should go back to first area
        cy.get('[data-cy="breadcrumb-item-1"]').find('a').click();
        cy.location('pathname').should('eq', firstAreaPath);

        // Go back to areas
        cy.get('[data-cy="breadcrumb-item-0"]').find('a').click();
        cy.location('pathname').should('eq', '/');

        // Navigate to second area
        cy.wait('@listAreas').its('response.statusCode').should('eq', 200);
        cy.get('[data-cy="area-item"]').eq(1).click();
        cy.wait('@listItemGroups').its('response.statusCode').should('eq', 200);

        cy.location('pathname').then((secondAreaPath) => {
          // Verify it's a different area
          expect(secondAreaPath).not.to.eq(firstAreaPath);

          openFirstItemGroupAndWaitItems();

          // Click area breadcrumb - should go to SECOND area, not first
          cy.get('[data-cy="breadcrumb-item-1"]').find('a').click();
          cy.location('pathname').should('eq', secondAreaPath);
        });
      });
    });
  });
});
