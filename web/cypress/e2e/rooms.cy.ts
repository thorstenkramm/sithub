import { openArea } from '../support/flows';

const testAuthEnabled = ['true', true, '1', 'yes'].includes(Cypress.env('testAuthEnabled'));
const itIfAuth = testAuthEnabled ? it : it.skip;

describe('rooms', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
  });

  itIfAuth('should show rooms for selected area', () => {
    cy.intercept('GET', '/api/v1/areas').as('listAreas');
    cy.intercept('GET', '/api/v1/areas/office_1st_floor/rooms').as('listRooms');

    cy.visit('/oauth/callback');
    openArea('Office 1st Floor');
    cy.wait('@listRooms').its('response.statusCode').should('eq', 200);
    cy.location('pathname').should('eq', '/areas/office_1st_floor/rooms');
    cy.get('[data-cy="rooms-list"]').should('exist');
    cy.get('[data-cy="room-item"]').first().should('contain', 'Room 101');
  });
});
