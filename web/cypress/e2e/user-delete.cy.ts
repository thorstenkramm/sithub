/**
 * E2E coverage for deleting/deactivating users (story 37.5).
 * Uses the real API; intercepts are only aliases for waiting on requests.
 */

const PASSWORD = 'E2eDeletePassword2026!';

function unique(prefix: string) {
  return `${prefix}-${Date.now()}-${Cypress._.random(1000, 9999)}`;
}

function serverFutureWeekday() {
  return cy.request('/api/v1/ping').then((response) => {
    const date = new Date(response.headers.date as string);
    date.setUTCDate(date.getUTCDate() + Cypress._.random(1, 25));
    while (date.getUTCDay() === 0 || date.getUTCDay() === 6) {
      date.setUTCDate(date.getUTCDate() + 1);
    }
    return cy.wrap(date.toISOString().slice(0, 10));
  });
}

function secondFutureWeekday(first: string) {
  const date = new Date(`${first}T00:00:00Z`);
  date.setUTCDate(date.getUTCDate() + 1);
  while (date.getUTCDay() === 0 || date.getUTCDay() === 6) {
    date.setUTCDate(date.getUTCDate() + 1);
  }
  return date.toISOString().slice(0, 10);
}

function createUser(email: string, displayName: string) {
  cy.request({
    method: 'POST',
    url: '/api/v1/users',
    headers: { 'Content-Type': 'application/vnd.api+json' },
    body: {
      data: {
        type: 'users',
        attributes: { email, display_name: displayName, password: PASSWORD }
      }
    }
  }).its('status').should('eq', 201);
}

function openUsersPage() {
  cy.visit('/');
  cy.get('[data-cy="user-menu-trigger"]').click();
  cy.get('[data-cy="users-menu-btn"]').click();
  cy.wait('@listUsers');
}

function rowFor(email: string) {
  return cy.contains('[data-cy="user-email"]', email).closest('tr');
}

function firstAvailableItemId(date: string) {
  return cy.request(`/api/v1/item-groups/test_room/items?date=${date}`).then((response) => {
    const items = response.body.data as Array<{ id: string; attributes: { availability: string } }>;
    const available = items.find((item) => item.attributes.availability === 'available') ?? items[0];
    expect(available, 'available test item').to.exist;
    return cy.wrap(available.id);
  });
}

describe('user deletion', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
    cy.login();
    cy.intercept('GET', '/api/v1/users*').as('listUsers');
    cy.intercept('GET', '/api/v1/users/*/deletion-impact').as('deletionImpact');
    cy.intercept('DELETE', '/api/v1/users/*').as('deleteUser');
  });

  it('should hard-delete a user without bookings and remove the row', () => {
    const id = unique('hard');
    const email = `aaa-${id}@sithub.local`;
    createUser(email, `AAA ${id}`);

    openUsersPage();
    rowFor(email).within(() => {
      cy.get('[data-cy="user-delete-btn"]').click();
    });
    cy.wait('@deletionImpact').its('response.body.data.attributes.mode').should('eq', 'delete');
    cy.get('[data-cy="confirm-dialog"]').should('contain', 'Delete user').and('contain', 'permanently deleted');
    cy.get('[data-cy="confirm-dialog-confirm"]').should('contain', 'Delete').click();
    cy.wait('@deleteUser').its('response.statusCode').should('eq', 204);
    cy.get('[data-cy="users-feedback"]').should('contain', `AAA ${id} was deleted.`);
    cy.contains('[data-cy="user-email"]', email).should('not.exist');
  });

  it('should deactivate a user with future bookings and block auth/session use', () => {
    const id = unique('soft');
    const email = `aaa-${id}@sithub.local`;
    const displayName = `AAA ${id}`;
    createUser(email, displayName);

    cy.request('POST', '/api/v1/auth/login', { email, password: PASSWORD });
    cy.getCookie('sithub_user').then((cookie) => {
      expect(cookie?.value, 'captured user session').to.be.a('string').and.not.be.empty;
      const userCookie = cookie!.value;

      serverFutureWeekday().then((bookingDate) => {
        const fallbackDate = secondFutureWeekday(bookingDate);
        firstAvailableItemId(bookingDate).then((itemID) => {
        cy.request({
          method: 'POST',
          url: '/api/v1/bookings',
          headers: { 'Content-Type': 'application/vnd.api+json' },
          failOnStatusCode: false,
          body: {
            data: {
              type: 'bookings',
              attributes: { item_id: itemID, booking_date: bookingDate }
            }
          }
        }).then((response) => {
          if (response.status === 201) return;
          cy.request({
            method: 'POST',
            url: '/api/v1/bookings',
            headers: { 'Content-Type': 'application/vnd.api+json' },
            body: {
              data: {
                type: 'bookings',
                attributes: { item_id: itemID, booking_date: fallbackDate }
              }
            }
          }).its('status').should('eq', 201);
        });
        });
      });

      cy.login();
      openUsersPage();
      rowFor(email).within(() => {
        cy.get('[data-cy="user-delete-btn"]').click();
      });
      cy.wait('@deletionImpact').then(({ response }) => {
        expect(response?.body.data.attributes.mode).to.eq('deactivate');
        expect(response?.body.data.attributes.future_booking_count).to.be.greaterThan(0);
      });
      cy.get('[data-cy="confirm-dialog"]')
        .should('contain', 'Deactivate user')
        .and('contain', 'upcoming booking')
        .and('not.contain', 'Delete');
      cy.get('[data-cy="confirm-dialog-confirm"]').should('contain', 'Deactivate').click();
      cy.wait('@deleteUser').its('response.statusCode').should('eq', 200);
      cy.get('[data-cy="users-feedback"]').should('contain', `${displayName} was deactivated`);
      rowFor(email).within(() => {
        cy.get('[data-cy="user-deleted-chip"]').should('contain', 'deactivated');
      });

      cy.request('/api/v1/colleagues').then((response) => {
        const names = (response.body.data as Array<{ attributes: { display_name: string } }>)
          .map((resource) => resource.attributes.display_name);
        expect(names).not.to.include(displayName);
      });
      cy.request({
        method: 'POST',
        url: '/api/v1/auth/login',
        body: { email, password: PASSWORD },
        failOnStatusCode: false
      }).then((response) => {
        expect(response.status).to.eq(401);
        expect(response.body.errors[0].code).to.eq('invalid_credentials');
      });
      cy.request({
        url: '/api/v1/bookings',
        headers: { Cookie: `sithub_user=${userCookie}` },
        failOnStatusCode: false
      }).its('status').should('eq', 401);
    });
  });

  it('should not offer a delete action on the admin own row', () => {
    openUsersPage();
    rowFor('anna@sithub.local').within(() => {
      cy.get('[data-cy="user-delete-btn"]').should('be.disabled');
    });
  });
});
