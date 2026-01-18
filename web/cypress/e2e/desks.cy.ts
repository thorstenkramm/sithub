import { openArea, openRoom } from '../support/flows';

const testAuthEnabled = ['true', true, '1', 'yes'].includes(Cypress.env('testAuthEnabled'));
const itIfAuth = testAuthEnabled ? it : it.skip;

describe('desks', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
  });

  itIfAuth('should show desks with equipment for a room', () => {
    cy.intercept('GET', '/api/v1/areas').as('listAreas');
    cy.intercept('GET', '/api/v1/areas/office_1st_floor/rooms').as('listRooms');
    cy.intercept('GET', '/api/v1/rooms/room_101/desks*').as('listDesks');

    cy.visit('/oauth/callback');
    openArea('Office 1st Floor');
    openRoom('Room 101');

    cy.wait('@listDesks').then((interception) => {
      expect(interception.response?.statusCode).to.eq(200);
      expect(interception.request.url).to.include('date=');
    });
    cy.location('pathname').should('eq', '/rooms/room_101/desks');
    cy.get('[data-cy="desk-item"]').first().should('contain', 'Desk 1');
    cy.get('[data-cy="desk-equipment"]').first().should('contain', '32 inch curved display, 2K');
    cy.get('[data-cy="desk-status"]').first().should('contain', 'Available');

    cy.get('[data-cy="desks-date"]')
      .invoke('val', '2026-01-20')
      .trigger('input')
      .trigger('change');
    cy.wait('@listDesks').its('request.url').should('include', 'date=2026-01-20');
  });
});
