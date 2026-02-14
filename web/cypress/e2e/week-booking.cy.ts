import { createMockItem, createMockItemsResponse } from '../support/flows';

const DESK_A_AVAILABLE = createMockItem('item-1', 'Desk A', 'available');

function mockAllItemsRequests() {
  cy.intercept('GET', '/api/v1/item-groups/*/items*', createMockItemsResponse([DESK_A_AVAILABLE])).as(
    'listItems'
  );
}

function visitItemsAndSwitchToWeek() {
  mockAllItemsRequests();
  cy.visit('/item-groups/test_room/items');
  cy.wait('@listItems');
  cy.get('[data-cy="mode-week-btn"]').click();
}

describe('week booking mode', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
    cy.login();
  });

  it('should switch to week mode and show week selector', () => {
    mockAllItemsRequests();
    cy.visit('/item-groups/test_room/items');
    cy.wait('@listItems');

    cy.get('[data-cy="booking-mode-toggle"]').should('be.visible');
    cy.get('[data-cy="mode-week-btn"]').click();

    cy.get('[data-cy="week-selector"]').should('be.visible');
    cy.get('[data-cy="items-date"]').should('not.exist');
  });

  it('should show per-day breakdown with checkboxes in week mode', () => {
    visitItemsAndSwitchToWeek();

    cy.get('[data-cy="week-items-list"]').should('be.visible');
    cy.get('[data-cy="week-item-entry"]').should('exist');
    cy.get('[data-cy="week-days"]').should('exist');
    cy.get('[data-cy-weekday="MO"]').should('exist');
    cy.get('[data-cy-weekday="FR"]').should('exist');
  });

  it('should show confirm button when days are selected', () => {
    visitItemsAndSwitchToWeek();

    cy.get('[data-cy="week-confirm-btn"]').should('not.exist');
    cy.get('[data-cy="week-day-checkbox"]').first().click();

    cy.get('[data-cy="week-confirm-btn"]').should('be.visible');
    cy.get('[data-cy="week-confirm-btn"]').should('contain', '1 day');
  });

  it('should submit bookings and show results', () => {
    visitItemsAndSwitchToWeek();

    cy.intercept('POST', '/api/v1/bookings', {
      statusCode: 201,
      headers: { 'Content-Type': 'application/vnd.api+json' },
      body: {
        data: {
          id: 'booking-week-1',
          type: 'bookings',
          attributes: {
            item_id: 'item-1',
            user_id: 'user-1',
            booking_date: '2026-02-10',
            created_at: '2026-02-09T10:00:00Z',
            note: ''
          }
        }
      }
    }).as('createBooking');

    cy.get('[data-cy="week-day-checkbox"]').first().click();
    cy.get('[data-cy="week-confirm-btn"]').click();
    cy.wait('@createBooking');

    cy.get('[data-cy="week-booking-results"]', { timeout: 10000 }).should('be.visible');
    cy.get('[data-cy="week-booking-results"]').should('contain', 'Booked');
  });

  it('should switch back to day mode and restore standard UI', () => {
    visitItemsAndSwitchToWeek();
    cy.get('[data-cy="week-selector"]').should('be.visible');

    cy.get('[data-cy="mode-day-btn"]').click();

    cy.get('[data-cy="items-list"]').should('exist');
    cy.get('[data-cy="book-item-btn"]').should('be.visible');
    cy.get('[data-cy="week-selector"]').should('not.exist');
  });

  it('should persist mode in localStorage across page reload', () => {
    visitItemsAndSwitchToWeek();
    cy.get('[data-cy="week-selector"]').should('be.visible');

    cy.window().then((win) => {
      expect(win.localStorage.getItem('sithub_booking_mode')).to.eq('week');
    });

    cy.visit('/item-groups/test_room/items');
    cy.get('[data-cy="week-selector"]').should('be.visible');
  });
});
