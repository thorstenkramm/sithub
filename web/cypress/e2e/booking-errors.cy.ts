// Story 37.1: failed API calls show a translated prefix plus the untranslated backend detail.

const toIsoDate = (date: Date): string => {
  const offset = date.getTimezoneOffset() * 60_000;
  return new Date(date.getTime() - offset).toISOString().slice(0, 10);
};

// Mirrors the backend booking horizon (internal/bookings/handler.go): the first rejected
// day is the Monday after the current week plus `weeksInAdvanced` weeks. The Tuesday after
// that Monday is a weekday safely beyond the limit.
const dayBeyondAdvanceLimit = (weeksInAdvanced: number): string => {
  const today = new Date();
  const daysUntilMonday = (8 - today.getDay()) % 7;
  const target = new Date(today);
  target.setDate(today.getDate() + daysUntilMonday + weeksInAdvanced * 7 + 1);
  return toIsoDate(target);
};

describe('booking errors', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
    cy.login();
  });

  it('should show the backend detail when booking beyond the advance limit', () => {
    cy.request('/api/v1/settings').then((response) => {
      const weeksInAdvanced = response.body.data.attributes.weeks_in_advanced as number;
      const bookingDate = dayBeyondAdvanceLimit(weeksInAdvanced);

      cy.intercept('GET', '/api/v1/item-groups/*/items*').as('listItems');
      cy.intercept('POST', '/api/v1/bookings*').as('createBooking');

      cy.visit('/item-groups/test_room/items', {
        onBeforeLoad(win) {
          win.localStorage.setItem('sithub_booking_mode', 'day');
          win.sessionStorage.setItem('sithub_selected_day', bookingDate);
        }
      });
      cy.wait('@listItems');

      cy.get('[data-cy="book-item-btn"]').first().click();
      cy.get('[data-cy="tile-booking-confirm"]').click();

      cy.wait('@createBooking').its('response.statusCode').should('eq', 400);
      cy.get('[data-cy="booking-error-text"]')
        .should('contain', 'Unable to book item. API error:')
        .and('contain', 'booking_date is too far in the future')
        .and('contain', bookingDate);
      cy.location('pathname').should('eq', '/item-groups/test_room/items');
    });
  });
});
