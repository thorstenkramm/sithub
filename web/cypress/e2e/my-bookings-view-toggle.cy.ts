function mockMyBookings() {
  cy.intercept('GET', '/api/v1/bookings', {
    statusCode: 200,
    headers: { 'Content-Type': 'application/vnd.api+json' },
    body: {
      data: [
        {
          id: 'booking-a',
          type: 'bookings',
          attributes: {
            item_id: 'desk-a',
            item_name: 'Corner Desk',
            item_group_id: 'ig-1',
            item_group_name: 'Room 101',
            area_id: 'area-1',
            area_name: 'Main Office',
            booking_date: '2026-02-10',
            created_at: '2026-02-09T10:00:00Z',
            booked_by_user_id: '',
            booked_by_user_name: '',
            booked_for_me: false,
            note: ''
          }
        }
      ]
    }
  }).as('listBookings');
}

describe('my bookings view toggle', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
    cy.login();
    mockMyBookings();
  });

  it('should default to the table view on desktop and persist a tile choice across reloads', () => {
    cy.viewport(1280, 800);
    cy.visit('/my-bookings');
    cy.wait('@listBookings');

    // Table is the desktop default.
    cy.get('[data-cy="bookings-table"]').should('be.visible');
    cy.get('[data-cy="bookings-list"]').should('not.exist');

    // Toggle to tiles.
    cy.get('[data-cy="view-switch"]').click();
    cy.get('[data-cy="bookings-list"]').should('be.visible');
    cy.get('[data-cy="bookings-table"]').should('not.exist');

    // Preference persists across a reload, overriding the desktop default.
    cy.reload();
    cy.wait('@listBookings');
    cy.get('[data-cy="bookings-list"]').should('be.visible');
    cy.get('[data-cy="bookings-table"]').should('not.exist');
  });
});
