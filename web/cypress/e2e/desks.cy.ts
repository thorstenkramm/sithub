import {
  createMockDesk,
  createMockDesksResponse,
  setupDesksPageIntercepts
} from '../support/flows';

describe('desks', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
    cy.login();
  });

  it('should show desks with equipment for a room', () => {
    setupDesksPageIntercepts();

    cy.visit('/');

    // Verify Vuetify is properly loaded (catches component import issues)
    cy.waitForVuetify();

    // Click first area (works with any test data)
    cy.wait('@listAreas').its('response.statusCode').should('eq', 200);
    cy.get('[data-cy="area-item"]').first().click();

    // Click first room
    cy.wait('@listRooms').its('response.statusCode').should('eq', 200);
    cy.get('[data-cy="room-item"]').first().click();

    cy.wait('@listDesks').then((interception) => {
      expect(interception.response?.statusCode).to.eq(200);
      expect(interception.request.url).to.include('date=');
    });
    cy.location('pathname').should('match', /\/rooms\/.*\/desks/);
    // Check desk items exist (name depends on config)
    cy.get('[data-cy="desk-item"]').should('have.length.at.least', 1);

    // Verify desk cards render as proper Vuetify components
    // Note: desk-item IS the v-card in the redesigned UI
    cy.get('[data-cy="desk-item"]').first().should('have.class', 'v-card');

    // Verify status chip renders as Vuetify chip
    // StatusChip component renders as a v-chip directly
    cy.get('[data-cy="desk-status"]').first().should('have.class', 'v-chip');
  });

  it('should book an available desk and show success message', () => {
    // Mock desks response with an available desk
    const mockDesk = createMockDesk('desk-available-1', 'Available Desk');
    cy.intercept('GET', '/api/v1/rooms/*/desks*', createMockDesksResponse([mockDesk])).as(
      'listDesks'
    );

    // Mock successful booking response
    cy.intercept('POST', '/api/v1/bookings', {
      statusCode: 201,
      headers: { 'Content-Type': 'application/vnd.api+json' },
      body: {
        data: {
          id: 'booking-1',
          type: 'bookings',
          attributes: {
            date: '2026-01-20',
            desk_id: 'desk-available-1'
          }
        }
      }
    }).as('createBooking');

    cy.visit('/rooms/test_room/desks');

    cy.wait('@listDesks');
    cy.location('pathname').should('match', /\/rooms\/.*\/desks/);

    // Click book button on the available desk
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
  });

  it('should show conflict message with prompt when desk is already booked', () => {
    // Mock desks response with an available desk
    const mockDesk = createMockDesk('desk-mock-1', 'Mock Desk 1');
    cy.intercept('GET', '/api/v1/rooms/*/desks*', createMockDesksResponse([mockDesk])).as(
      'listDesks'
    );

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

    cy.visit('/rooms/test_room/desks');

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

  it('should show self-duplicate message when user already has booking', () => {
    // Mock desks response with an available desk
    const mockDesk = createMockDesk('desk-mock-2', 'Mock Desk 2');
    cy.intercept('GET', '/api/v1/rooms/*/desks*', createMockDesksResponse([mockDesk])).as(
      'listDesks'
    );

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

    cy.visit('/rooms/test_room/desks');

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
