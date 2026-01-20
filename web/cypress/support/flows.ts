export function openArea(areaName: string) {
  cy.wait('@listAreas').its('response.statusCode').should('eq', 200);
  cy.contains('[data-cy="area-item"]', areaName).click();
}

export function openRoom(roomName: string) {
  cy.wait('@listRooms').its('response.statusCode').should('eq', 200);
  cy.contains('[data-cy="room-item"]', roomName).click();
}

/** Sets up intercepts for areas, rooms, and desks API endpoints */
export function setupDesksPageIntercepts() {
  cy.intercept('GET', '/api/v1/areas').as('listAreas');
  cy.intercept('GET', '/api/v1/areas/*/rooms').as('listRooms');
  cy.intercept('GET', '/api/v1/rooms/*/desks*').as('listDesks');
}

/** Creates a mock desk object for testing */
export function createMockDesk(
  id: string,
  name: string,
  availability: 'available' | 'occupied' = 'available',
  equipment: string[] = ['Monitor']
) {
  return {
    id,
    type: 'desks',
    attributes: {
      name,
      equipment,
      availability
    }
  };
}

/** Creates a mock desks API response body */
export function createMockDesksResponse(desks: ReturnType<typeof createMockDesk>[]) {
  return {
    statusCode: 200,
    body: {
      data: desks
    }
  };
}
