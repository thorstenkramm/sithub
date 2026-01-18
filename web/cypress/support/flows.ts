export function openArea(areaName: string) {
  cy.wait('@listAreas').its('response.statusCode').should('eq', 200);
  cy.contains('[data-cy="area-item"]', areaName).click();
}

export function openRoom(roomName: string) {
  cy.wait('@listRooms').its('response.statusCode').should('eq', 200);
  cy.contains('[data-cy="room-item"]', roomName).click();
}
