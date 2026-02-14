import { createMockItem, createMockItemsResponse } from '../support/flows';

function mockBookingCreate(bookingId: string, itemId: string) {
  cy.intercept('POST', '/api/v1/bookings', {
    statusCode: 201,
    headers: { 'Content-Type': 'application/vnd.api+json' },
    body: {
      data: {
        id: bookingId,
        type: 'bookings',
        attributes: {
          item_id: itemId,
          user_id: 'user-1',
          booking_date: '2026-02-10',
          created_at: '2026-02-09T10:00:00Z',
          note: ''
        }
      }
    }
  }).as('createBooking');
}

function setupItemsAndBooking(itemId: string, itemName: string, bookingId: string) {
  const mockItem = createMockItem(itemId, itemName);
  cy.intercept('GET', '/api/v1/item-groups/*/items*', createMockItemsResponse([mockItem])).as(
    'listItems'
  );
  mockBookingCreate(bookingId, itemId);
}

function bookFirstAvailableItem() {
  cy.get('[data-cy="item-entry"][data-cy-availability="available"]')
    .first()
    .find('[data-cy="book-item-btn"]')
    .click();
}

describe('booking notes', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
    cy.login();
  });

  it('should show add note action after successful booking', () => {
    setupItemsAndBooking('item-note-1', 'Note Test Item', 'booking-note-1');

    cy.visit('/item-groups/test_room/items');
    cy.wait('@listItems');

    bookFirstAvailableItem();
    cy.wait('@createBooking');

    // Success message should appear with add note button
    cy.get('[data-cy="booking-success-text"]').should('contain', 'Note Test Item');
    cy.get('[data-cy="add-note-after-booking"]').should('be.visible');
  });

  it('should add note after booking via dialog', () => {
    setupItemsAndBooking('item-note-2', 'Note Test Item 2', 'booking-note-2');

    cy.intercept('PATCH', '/api/v1/bookings/booking-note-2', {
      statusCode: 200,
      headers: { 'Content-Type': 'application/vnd.api+json' },
      body: {
        data: {
          id: 'booking-note-2',
          type: 'bookings',
          attributes: {
            item_id: 'item-note-2',
            user_id: 'user-1',
            booking_date: '2026-02-10',
            created_at: '2026-02-09T10:00:00Z',
            note: 'Arriving after 2pm'
          }
        }
      }
    }).as('updateNote');

    cy.visit('/item-groups/test_room/items');
    cy.wait('@listItems');

    bookFirstAvailableItem();
    cy.wait('@createBooking');

    // Click add note
    cy.get('[data-cy="add-note-after-booking"]').click();

    // Dialog should open
    cy.get('[data-cy="post-booking-note-input"]').should('be.visible');
    cy.get('[data-cy="post-booking-note-input"]').type('Arriving after 2pm');
    cy.get('[data-cy="post-booking-note-save"]').click();

    cy.wait('@updateNote').then((interception) => {
      expect(interception.request.body.data.attributes.note).to.eq('Arriving after 2pm');
    });

    // Success message should update
    cy.get('[data-cy="booking-success-text"]').should('contain', 'Note Test Item 2');
  });

  it('should display and edit note on booking card in My Bookings', () => {
    cy.intercept('GET', '/api/v1/bookings', {
      statusCode: 200,
      headers: { 'Content-Type': 'application/vnd.api+json' },
      body: {
        data: [
          {
            id: 'booking-with-note',
            type: 'bookings',
            attributes: {
              item_id: 'item-1',
              item_name: 'Corner Desk',
              item_group_id: 'ig-1',
              item_group_name: 'Room 101',
              area_id: 'area-1',
              area_name: 'Main Office',
              booking_date: '2026-02-10',
              created_at: '2026-02-09T10:00:00Z',
              note: 'Will arrive late, after 2pm',
              booked_by_user_id: '',
              booked_by_user_name: '',
              booked_for_me: false
            }
          }
        ]
      }
    }).as('listBookings');

    cy.intercept('PATCH', '/api/v1/bookings/booking-with-note', {
      statusCode: 200,
      headers: { 'Content-Type': 'application/vnd.api+json' },
      body: {
        data: {
          id: 'booking-with-note',
          type: 'bookings',
          attributes: {
            item_id: 'item-1',
            user_id: 'user-1',
            booking_date: '2026-02-10',
            created_at: '2026-02-09T10:00:00Z',
            note: 'Updated: arriving at 3pm'
          }
        }
      }
    }).as('updateNote');

    cy.visit('/my-bookings');
    cy.wait('@listBookings');

    // Note should be displayed
    cy.get('[data-cy="booking-note"]').should('contain', 'Will arrive late');

    // Click edit note
    cy.get('[data-cy="edit-note-btn"]').click();

    // Edit dialog should open with existing note
    cy.get('[data-cy="note-edit-input"]').should('be.visible');
    cy.get('[data-cy="note-edit-input"]').clear().type('Updated: arriving at 3pm');
    cy.get('[data-cy="note-save-btn"]').click();

    cy.wait('@updateNote').then((interception) => {
      expect(interception.request.body.data.attributes.note).to.eq('Updated: arriving at 3pm');
    });
  });

  it('should show add note button when booking has no note', () => {
    cy.intercept('GET', '/api/v1/bookings', {
      statusCode: 200,
      headers: { 'Content-Type': 'application/vnd.api+json' },
      body: {
        data: [
          {
            id: 'booking-no-note',
            type: 'bookings',
            attributes: {
              item_id: 'item-1',
              item_name: 'Corner Desk',
              item_group_id: 'ig-1',
              item_group_name: 'Room 101',
              area_id: 'area-1',
              area_name: 'Main Office',
              booking_date: '2026-02-10',
              created_at: '2026-02-09T10:00:00Z',
              note: '',
              booked_by_user_id: '',
              booked_by_user_name: '',
              booked_for_me: false
            }
          }
        ]
      }
    }).as('listBookings');

    cy.visit('/my-bookings');
    cy.wait('@listBookings');

    // No note should be displayed
    cy.get('[data-cy="booking-note"]').should('not.exist');

    // Add note button should be visible
    cy.get('[data-cy="add-note-btn"]').should('be.visible');
  });

  it('should show note on occupied items in items view', () => {
    const occupiedItem = {
      id: 'item-occupied-note',
      type: 'items',
      attributes: {
        name: 'Occupied Desk',
        equipment: ['Monitor'],
        availability: 'occupied',
        booker_name: 'Alice Smith',
        note: 'Working from home in the morning'
      }
    };

    cy.intercept('GET', '/api/v1/item-groups/*/items*', {
      statusCode: 200,
      body: { data: [occupiedItem] }
    }).as('listItems');

    cy.visit('/item-groups/test_room/items');
    cy.wait('@listItems');

    // Note should be displayed on occupied item
    cy.get('[data-cy="item-note"]').should('contain', 'Working from home');
  });
});
