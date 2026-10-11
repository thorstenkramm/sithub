/**
 * E2E coverage for the data exporter role (story 37.6).
 * Uses the real backend; intercepts are only aliases for waiting on requests.
 * Felix (non-admin demo user) is the subject; his flag is reset around every test.
 */

const FELIX_EMAIL = 'felix@sithub.local';
const ALEX_EMAIL = 'alex@sithub.local';
const DEMO_PASSWORD = 'SitHubDemo2026!!';
const PAGE_SIZE = 20;

interface ListedUser {
  id: string;
  attributes: { email: string };
}

/** Resolves the user id and the 1-based table page (20 per page) of a user via the admin API. */
function locateUser(email: string, offset = 0): Cypress.Chainable<{ id: string; page: number }> {
  return cy.request(`/api/v1/users?page[limit]=100&page[offset]=${offset}`).then((response) => {
    const users = response.body.data as ListedUser[];
    const index = users.findIndex((user) => user.attributes.email === email);
    if (index >= 0) {
      return cy.wrap({ id: users[index].id, page: Math.floor((offset + index) / PAGE_SIZE) + 1 });
    }
    expect(response.body.links.next, `${email} exists`).to.be.a('string');
    return locateUser(email, offset + 100);
  });
}

function setFelixDataExporter(flag: boolean) {
  cy.login();
  locateUser(FELIX_EMAIL).then(({ id }) => {
    cy.request({
      method: 'PATCH',
      url: `/api/v1/users/${id}`,
      headers: { 'Content-Type': 'application/vnd.api+json' },
      body: { data: { type: 'users', id, attributes: { is_data_exporter: flag } } }
    }).its('status').should('eq', 200);
  });
}

/** Opens the Users page on the table page that contains the given user and returns its row. */
function openUsersPageFor(email: string) {
  cy.visit('/');
  cy.get('[data-cy="user-menu-trigger"]').click();
  cy.get('[data-cy="users-menu-btn"]').click();
  cy.wait('@listUsers');
  locateUser(email).then(({ page }) => {
    if (page > 1) {
      cy.get('[data-cy="users-pagination"]').contains('button', String(page)).click();
      cy.wait('@listUsers');
    }
  });
  return cy.contains('[data-cy="user-email"]', email).closest('tr');
}

describe('data exporter role', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
    cy.intercept('GET', '/api/v1/users*').as('listUsers');
    cy.intercept('PATCH', '/api/v1/users/*').as('updateUser');
    setFelixDataExporter(false);
  });

  after(() => {
    setFelixDataExporter(false);
  });

  it('should show Admin or User in the role column', () => {
    cy.login();
    openUsersPageFor('anna@sithub.local').within(() => {
      cy.get('[data-cy="user-roles"]').should('contain.text', 'Admin');
    });
    openUsersPageFor(FELIX_EMAIL).within(() => {
      cy.get('[data-cy="user-roles"]').should('contain.text', 'User');
    });
  });

  it('should grant and revoke the data exporter role via the toggle', () => {
    cy.login();
    openUsersPageFor(FELIX_EMAIL).within(() => {
      cy.get('[data-cy="user-data-exporter-toggle"] input').check({ force: true });
    });
    cy.wait('@updateUser').then(({ request, response }) => {
      expect(request.body.data.attributes).to.deep.eq({ is_data_exporter: true });
      expect(response?.statusCode).to.eq(200);
      expect(response?.body.data.attributes.is_data_exporter).to.eq(true);
    });
    cy.get('[data-cy="users-feedback"]').should('contain.text', 'granted');
    cy.contains('[data-cy="user-email"]', FELIX_EMAIL).closest('tr').within(() => {
      cy.get('[data-cy="user-roles"]').should('contain.text', 'Data exporter');
      cy.get('[data-cy="user-data-exporter-toggle"] input').uncheck({ force: true });
    });
    cy.wait('@updateUser').its('response.body.data.attributes.is_data_exporter').should('eq', false);
    cy.contains('[data-cy="user-email"]', FELIX_EMAIL).closest('tr').within(() => {
      cy.get('[data-cy="user-roles"]').should('not.contain.text', 'Data exporter');
    });
  });

  it('should expose is_data_exporter on /me without a new login', () => {
    cy.request('POST', '/api/v1/auth/login', { email: FELIX_EMAIL, password: DEMO_PASSWORD });
    cy.getCookie('sithub_user').then((cookie) => {
      const felixCookie = cookie!.value;
      const meAsFelix = () => cy.request({
        url: '/api/v1/me',
        headers: { Cookie: `sithub_user=${felixCookie}` }
      }).its('body.data.attributes.is_data_exporter');

      meAsFelix().should('eq', false);
      setFelixDataExporter(true);
      meAsFelix().should('eq', true);
      setFelixDataExporter(false);
      meAsFelix().should('eq', false);
    });
  });

  it('should show Access Denied on /data-export for an admin without the flag', () => {
    cy.login();
    cy.visit('/data-export');
    cy.location('pathname').should('eq', '/access-denied');
  });

  it('should show Access Denied on /data-export for a regular user', () => {
    cy.login(ALEX_EMAIL, DEMO_PASSWORD);
    cy.visit('/data-export');
    cy.location('pathname').should('eq', '/access-denied');
  });

  it('should open /data-export for a data exporter', () => {
    setFelixDataExporter(true);
    cy.login(FELIX_EMAIL, DEMO_PASSWORD);
    cy.visit('/data-export');
    cy.get('[data-cy="data-export-page"]').should('be.visible');
    cy.location('pathname').should('eq', '/data-export');
  });
});
