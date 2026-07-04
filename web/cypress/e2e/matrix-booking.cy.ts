import {
  resetAndLogin,
  setupItemsPageIntercepts,
  visitHomeAndWaitForAreas,
  openFirstAreaAndWaitItemGroups
} from '../support/flows';

/**
 * Creates a mock matrix API response for E2E testing.
 * Returns a single room with the given items/cells.
 */
function makeMatrixResponse(items: Array<{
  id: string;
  name: string;
  warning?: string;
  reserved?: boolean;
  cells: Array<{
    date: string;
    availability: 'free' | 'occupied';
    booker_name?: string;
    booker_user_id?: string;
    booked_by_me?: boolean;
    booking_id?: string;
  }>;
}>) {
  const days = [
    { date: '2099-06-02', weekday: 'MO' },
    { date: '2099-06-03', weekday: 'TU' },
    { date: '2099-06-04', weekday: 'WE' },
    { date: '2099-06-05', weekday: 'TH' },
    { date: '2099-06-06', weekday: 'FR' }
  ];

  return {
    data: [{
      id: 'ig-e2e',
      type: 'item-group-weekly-matrix',
      attributes: {
        item_group_id: 'ig-e2e',
        item_group_name: 'Test Room',
        days,
        items: items.map(item => ({
          item_id: item.id,
          item_name: item.name,
          equipment: [],
          warning: item.warning ?? null,
          reserved: item.reserved ?? false,
          cells: item.cells.map(c => ({
            date: c.date,
            availability: c.availability,
            booker_name: c.booker_name ?? '',
            booker_user_id: c.booker_user_id ?? '',
            booked_by_me: c.booked_by_me ?? false,
            booking_id: c.booking_id ?? ''
          }))
        }))
      }
    }]
  };
}

/** Navigates to item groups and switches to table view with a mocked matrix. */
function setupMatrixView(matrixBody: ReturnType<typeof makeMatrixResponse>) {
  setupItemsPageIntercepts();
  cy.intercept('GET', /\/api\/v1\/areas\/[^/]+\/item-groups\/matrix/, {
    statusCode: 200,
    body: matrixBody
  }).as('matrix');

  visitHomeAndWaitForAreas();
  openFirstAreaAndWaitItemGroups();

  // Switch to table view
  cy.get('[data-cy="view-switch"]').click();
  cy.wait('@matrix');
  cy.get('[data-cy="area-weekly-matrix"]').should('exist');
}

/** Stubs a successful booking (201) aliased as `@createBooking`. */
function interceptBookingSuccess() {
  cy.intercept('POST', '/api/v1/bookings', {
    statusCode: 201,
    headers: { 'Content-Type': 'application/vnd.api+json' },
    body: {
      data: { id: 'b-new', type: 'bookings', attributes: { item_id: 'desk-1', booking_date: '2099-06-02', note: '' } }
    }
  }).as('createBooking');
}

describe('matrix booking popover', () => {
  beforeEach(() => {
    resetAndLogin();
  });

  it('should open booking popover on free cell click and book successfully', () => {
    const matrix = makeMatrixResponse([{
      id: 'desk-1',
      name: 'Corner Desk',
      cells: [
        { date: '2099-06-02', availability: 'free' },
        { date: '2099-06-03', availability: 'free' }
      ]
    }]);
    setupMatrixView(matrix);

    // Mock successful booking
    interceptBookingSuccess();

    // Re-intercept matrix for the refresh after booking
    cy.intercept('GET', /\/api\/v1\/areas\/[^/]+\/item-groups\/matrix/, {
      statusCode: 200,
      body: matrix
    }).as('matrixRefresh');

    // Click a free cell
    cy.get('[data-cy="matrix-cell-free"]').first().click();

    // Popover should open with booking controls
    cy.get('[data-cy="matrix-booking-card"]').should('be.visible');
    cy.get('[data-cy="matrix-book-self-radio"]').should('exist');
    cy.get('[data-cy="matrix-book-colleague-radio"]').should('exist');
    cy.get('[data-cy="matrix-booking-note"]').should('exist');

    // Type a note and confirm
    cy.get('[data-cy="matrix-booking-note"]').find('input').type('Arriving late');
    cy.get('[data-cy="matrix-booking-confirm"]').click();

    // Verify the booking request includes the note
    cy.wait('@createBooking').then((interception) => {
      expect(interception.response?.statusCode).to.eq(201);
      const attrs = interception.request.body.data.attributes;
      expect(attrs.item_id).to.eq('desk-1');
      expect(attrs.note).to.eq('Arriving late');
    });

    // Popover should close, snackbar should show
    cy.get('[data-cy="matrix-booking-card"]').should('not.exist');
    cy.get('[data-cy="matrix-snackbar"]').should('contain', 'Booking confirmed');
  });

  it('should show only the cancel popover after switching from a free to an occupied cell', () => {
    // Regression: both popovers shared one activator and the booking popover
    // stayed mounted, so clicking an occupied cell showed the booking AND the
    // cancel popover at once. Only one popover may ever be visible.
    const matrix = makeMatrixResponse([
      {
        id: 'desk-1',
        name: 'Free Desk',
        cells: [{ date: '2099-06-02', availability: 'free' }]
      },
      {
        id: 'desk-2',
        name: 'My Desk',
        cells: [{
          date: '2099-06-02',
          availability: 'occupied',
          booker_name: 'Test User',
          booker_user_id: 'me',
          booked_by_me: true,
          booking_id: 'b-mine'
        }]
      }
    ]);
    setupMatrixView(matrix);

    // Open the booking popover on the free cell.
    cy.get('[data-cy="matrix-cell-free"]').first().click();
    cy.get('[data-cy="matrix-booking-card"]').should('be.visible');

    // Switching to the occupied cell must swap to the cancel popover only.
    cy.get('[data-cy="matrix-cell-occupied"]').first().click();
    cy.get('[data-cy="matrix-cancel-card"]').should('be.visible');
    cy.get('[data-cy="matrix-booking-card"]').should('not.exist');
  });

  it('should show inline error on 409 conflict and keep popover open', () => {
    const matrix = makeMatrixResponse([{
      id: 'desk-1',
      name: 'Corner Desk',
      cells: [{ date: '2099-06-02', availability: 'free' }]
    }]);
    setupMatrixView(matrix);

    cy.intercept('POST', '/api/v1/bookings', {
      statusCode: 409,
      headers: { 'Content-Type': 'application/vnd.api+json' },
      body: {
        errors: [{ status: '409', title: 'Conflict', detail: 'Item is already booked for this date' }]
      }
    }).as('createBookingConflict');

    cy.intercept('GET', /\/api\/v1\/areas\/[^/]+\/item-groups\/matrix/, {
      statusCode: 200,
      body: matrix
    }).as('matrixRefresh');

    cy.get('[data-cy="matrix-cell-free"]').first().click();
    cy.get('[data-cy="matrix-booking-confirm"]').click();
    cy.wait('@createBookingConflict');

    // Popover stays open with inline error
    cy.get('[data-cy="matrix-booking-card"]').should('be.visible');
    cy.get('[data-cy="matrix-booking-error"]').should('be.visible');
  });

  it('surfaces the uniform warning confirmation when booking a warned item', () => {
    const matrix = makeMatrixResponse([{
      id: 'desk-1',
      name: 'Window Desk',
      warning: 'Near noisy area',
      cells: [{ date: '2099-06-02', availability: 'free' }]
    }]);
    setupMatrixView(matrix);
    interceptBookingSuccess();

    cy.get('[data-cy="matrix-cell-free"]').first().click();
    // The old inline in-popover warning is gone; warnings now use the uniform dialog.
    cy.get('[data-cy="matrix-booking-warning"]').should('not.exist');

    // Clicking Book surfaces the shared confirmation dialog before any booking.
    cy.get('[data-cy="matrix-booking-confirm"]').click();
    cy.get('[data-cy="warning-dialog"]').should('be.visible');
    cy.get('[data-cy="warning-message"]').should('contain', 'Near noisy area');

    // Confirming the warning proceeds with the booking.
    cy.get('[data-cy="warning-confirm-btn"]').click();
    cy.wait('@createBooking').its('response.statusCode').should('eq', 201);
  });
});

describe('matrix cancellation popover', () => {
  beforeEach(() => {
    resetAndLogin();
  });

  it('should open cancel popover on own booking and cancel successfully', () => {
    const matrix = makeMatrixResponse([{
      id: 'desk-1',
      name: 'My Desk',
      cells: [{
        date: '2099-06-02',
        availability: 'occupied',
        booker_name: 'Test User',
        booker_user_id: 'me',
        booked_by_me: true,
        booking_id: 'b-mine'
      }]
    }]);
    setupMatrixView(matrix);

    cy.intercept('DELETE', '/api/v1/bookings/b-mine', { statusCode: 204 }).as('cancelBooking');
    cy.intercept('GET', /\/api\/v1\/areas\/[^/]+\/item-groups\/matrix/, {
      statusCode: 200,
      body: matrix
    }).as('matrixRefresh');

    // Click occupied cell
    cy.get('[data-cy="matrix-cell-occupied"]').first().click();

    // Cancel popover should show person, desk, date
    cy.get('[data-cy="matrix-cancel-card"]').should('be.visible');
    cy.get('[data-cy="matrix-cancel-person"]').should('contain', 'Test User');
    cy.get('[data-cy="matrix-cancel-desk"]').should('contain', 'My Desk');
    cy.get('[data-cy="matrix-cancel-date"]').should('contain', '2099-06-02');

    // Confirm cancellation
    cy.get('[data-cy="matrix-cancel-confirm"]').click();
    cy.wait('@cancelBooking').its('response.statusCode').should('eq', 204);

    // Popover should close, snackbar should show
    cy.get('[data-cy="matrix-cancel-card"]').should('not.exist');
    cy.get('[data-cy="matrix-snackbar"]').should('contain', 'cancelled');
  });

  it('should not open popup on non-admin other-user occupied cell', () => {
    // Login as a non-admin user
    cy.clearCookies();
    cy.login('alex@sithub.local', 'SitHubDemo2026!!');

    const matrix = makeMatrixResponse([{
      id: 'desk-1',
      name: 'Other Desk',
      cells: [{
        date: '2099-06-02',
        availability: 'occupied',
        booker_name: 'Someone Else',
        booker_user_id: 'other-user',
        booked_by_me: false
      }]
    }]);
    setupMatrixView(matrix);

    // Click the inert occupied cell
    cy.get('[data-cy="matrix-cell-occupied"]').first().click();

    // No popover should open
    cy.get('[data-cy="matrix-cancel-card"]').should('not.exist');
    cy.get('[data-cy="matrix-booking-card"]').should('not.exist');
  });
});
