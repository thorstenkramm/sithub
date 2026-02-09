export function openArea(areaName: string) {
  cy.wait('@listAreas').its('response.statusCode').should('eq', 200);
  cy.contains('[data-cy="area-item"]', areaName).click();
}

export function openItemGroup(itemGroupName: string) {
  cy.wait('@listItemGroups').its('response.statusCode').should('eq', 200);
  cy.contains('[data-cy="item-group-item"]', itemGroupName).click();
}

/** Sets up intercepts for areas, item groups, and items API endpoints */
export function setupItemsPageIntercepts() {
  cy.intercept('GET', '/api/v1/areas').as('listAreas');
  cy.intercept('GET', '/api/v1/areas/*/item-groups').as('listItemGroups');
  cy.intercept('GET', '/api/v1/item-groups/*/items*').as('listItems');
}

/** Creates a mock item object for testing */
export function createMockItem(
  id: string,
  name: string,
  availability: 'available' | 'occupied' = 'available',
  equipment: string[] = ['Monitor']
) {
  return {
    id,
    type: 'items',
    attributes: {
      name,
      equipment,
      availability
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
