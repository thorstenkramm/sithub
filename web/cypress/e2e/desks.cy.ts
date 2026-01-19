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

  itIfAuth('should book an available desk and show success message', () => {
    cy.intercept('GET', '/api/v1/areas').as('listAreas');
    cy.intercept('GET', '/api/v1/areas/office_1st_floor/rooms').as('listRooms');
    cy.intercept('GET', '/api/v1/rooms/room_101/desks*').as('listDesks');
    cy.intercept('POST', '/api/v1/bookings').as('createBooking');

    cy.visit('/oauth/callback');
    openArea('Office 1st Floor');
    openRoom('Room 101');

    cy.wait('@listDesks');
    cy.location('pathname').should('eq', '/rooms/room_101/desks');

    // Issue 7: Use improved selector with data-cy-availability attribute
    cy.get('[data-cy="desk-item"][data-cy-availability="available"]')
      .first()
      .find('[data-cy="book-desk-btn"]')
      .click();

    cy.wait('@createBooking').then((interception) => {
      expect(interception.response?.statusCode).to.eq(201);
      expect(interception.request.body.data.type).to.eq('bookings');
    });

    // Success message should appear
    cy.get('[data-cy="booking-success"]').should('contain', 'Desk booked successfully');

    // After reload, the desk should show as occupied
    cy.wait('@listDesks');
  });

  itIfAuth('should show conflict message with prompt when desk is already booked', () => {
    cy.intercept('GET', '/api/v1/areas').as('listAreas');
    cy.intercept('GET', '/api/v1/areas/office_1st_floor/rooms').as('listRooms');
    cy.intercept('GET', '/api/v1/rooms/room_101/desks*').as('listDesks');

    // Mock a 409 Conflict response for booking
    cy.intercept('POST', '/api/v1/bookings', {
      statusCode: 409,
      headers: { 'Content-Type': 'application/vnd.api+json' },
      body: {
        errors: [
          {
            status: '409',
            title: 'Conflict',
            detail: 'Desk is already booked for this date',
            code: 'conflict'
          }
        ]
      }
    }).as('createBookingConflict');

    cy.visit('/oauth/callback');
    openArea('Office 1st Floor');
    openRoom('Room 101');

    cy.wait('@listDesks');

    // Click book on an available desk
    cy.get('[data-cy="desk-item"][data-cy-availability="available"]')
      .first()
      .find('[data-cy="book-desk-btn"]')
      .click();

    cy.wait('@createBookingConflict');

    // Error message should show backend detail + prompt
    cy.get('[data-cy="booking-error"]')
      .should('contain', 'Desk is already booked for this date')
      .and('contain', 'Please choose another desk');

    // Desk list should be refreshed
    cy.wait('@listDesks');
  });

  itIfAuth('should show self-duplicate message when user already has booking', () => {
    cy.intercept('GET', '/api/v1/areas').as('listAreas');
    cy.intercept('GET', '/api/v1/areas/office_1st_floor/rooms').as('listRooms');
    cy.intercept('GET', '/api/v1/rooms/room_101/desks*').as('listDesks');

    // Mock a 409 Conflict response for self-duplicate
    cy.intercept('POST', '/api/v1/bookings', {
      statusCode: 409,
      headers: { 'Content-Type': 'application/vnd.api+json' },
      body: {
        errors: [
          {
            status: '409',
            title: 'Conflict',
            detail: 'You already have this desk booked for this date',
            code: 'conflict'
          }
        ]
      }
    }).as('createBookingSelfDuplicate');

    cy.visit('/oauth/callback');
    openArea('Office 1st Floor');
    openRoom('Room 101');

    cy.wait('@listDesks');

    cy.get('[data-cy="desk-item"][data-cy-availability="available"]')
      .first()
      .find('[data-cy="book-desk-btn"]')
      .click();

    cy.wait('@createBookingSelfDuplicate');

    // Error message should show the self-duplicate message
    cy.get('[data-cy="booking-error"]')
      .should('contain', 'You already have this desk booked for this date')
      .and('contain', 'Please choose another desk');
  });
});
