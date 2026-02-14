export function openArea(areaName: string) {
  cy.wait('@listAreas').its('response.statusCode').should('eq', 200);
  cy.contains('[data-cy="area-item"]', areaName).click();
}

export function openItemGroup(itemGroupName: string) {
  cy.wait('@listItemGroups').its('response.statusCode').should('eq', 200);
  cy.contains('[data-cy="item-group-item"]', itemGroupName).click();
}

export function resetAndLogin() {
  cy.clearCookies();
  cy.clearLocalStorage();
  cy.login();
}

export function visitHomeAndWaitForAreas() {
  cy.visit('/');
  cy.waitForVuetify();
  cy.wait('@listAreas').its('response.statusCode').should('eq', 200);
}

export function openFirstAreaAndWaitItemGroups() {
  cy.get('[data-cy="area-item"]').first().click();
  cy.wait('@listItemGroups').its('response.statusCode').should('eq', 200);
  cy.location('pathname').should('match', /\/areas\/.*\/item-groups/);
}

export function openFirstItemGroupAndWaitItems() {
  cy.get('[data-cy="item-group-item"]').first().click();
  cy.wait('@listItems').its('response.statusCode').should('eq', 200);
  cy.location('pathname').should('match', /\/item-groups\/.*\/items/);
}

/** Sets up intercepts for areas, item groups, items, and availability API endpoints */
export function setupItemsPageIntercepts() {
  cy.intercept('GET', '/api/v1/areas').as('listAreas');
  cy.intercept('GET', /\/api\/v1\/areas\/[^/]+\/item-groups\/availability/).as('availability');
  cy.intercept('GET', /\/api\/v1\/areas\/[^/]+\/item-groups$/).as('listItemGroups');
  cy.intercept('GET', '/api/v1/item-groups/*/items*').as('listItems');
}

/** Creates a mock item object for testing */
export function createMockItem(
  id: string,
  name: string,
  availability: 'available' | 'occupied' = 'available',
  equipment: string[] = ['Monitor'],
  bookerName?: string
) {
  return {
    id,
    type: 'items',
    attributes: {
      name,
      equipment,
      availability,
      ...(bookerName ? { booker_name: bookerName } : {})
    }
  };
}

/** Creates a mock items API response body */
export function createMockItemsResponse(items: ReturnType<typeof createMockItem>[]) {
  return {
    statusCode: 200,
    body: {
      data: items
    }
  };
}
