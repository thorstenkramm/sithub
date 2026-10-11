/**
 * E2E tests for the admin users table (story 37.4).
 *
 * The demo data provides 15 users; the spec seeds 10 more so the table spans two pages.
 * Seeding is idempotent (201 or 409) so reruns against the same database work.
 */

const PAGE_SIZE = 20;
const SEED_PASSWORD = 'E2ePagerPassword2026!';

interface UserListBody {
  data: { attributes: { email: string; display_name: string } }[];
  meta: { total: number; limit: number; offset: number };
  links: { self: string };
}

/** SQLite COLLATE NOCASE only folds ASCII letters; mirror that for the order assertion. */
function asciiLower(value: string) {
  return value.replace(/[A-Z]/g, (char) => char.toLowerCase());
}

function seedPagerUsers() {
  for (let i = 1; i <= 10; i++) {
    const nr = String(i).padStart(2, '0');
    cy.request({
      method: 'POST',
      url: '/api/v1/users',
      headers: { 'Content-Type': 'application/vnd.api+json' },
      body: {
        data: {
          type: 'users',
          attributes: {
            email: `e2e-pager-${nr}@sithub.local`,
            display_name: `E2E Pager ${nr}`,
            password: SEED_PASSWORD
          }
        }
      },
      failOnStatusCode: false
    }).its('status').should('be.oneOf', [201, 409]);
  }
}

function openUsersPage() {
  cy.visit('/');
  cy.get('[data-cy="user-menu-trigger"]').click();
  cy.get('[data-cy="users-menu-btn"]').click();
}

describe('admin users table', () => {
  before(() => {
    cy.login();
    seedPagerUsers();
  });

  beforeEach(() => {
    cy.clearLocalStorage();
    cy.login();
    cy.intercept('GET', '/api/v1/users*').as('listUsers');
  });

  it('should show the Users menu entry for admins and open the table', () => {
    openUsersPage();

    cy.wait('@listUsers').then(({ request, response }) => {
      expect(request.url).to.include('page%5Blimit%5D=20');
      expect(request.url).to.include('page%5Boffset%5D=0');
      expect(response?.statusCode).to.eq(200);
    });
    cy.location('pathname').should('eq', '/admin/users');
    cy.get('[data-cy="users-table"]').find('th').should('have.length.at.least', 5);
    cy.get('[data-cy="user-email"]').should('have.length', PAGE_SIZE);
    cy.get('[data-cy="users-pagination"]').should('be.visible').contains('button', '2');
  });

  it('should page through all users exactly once sorted by name', () => {
    const emails: string[] = [];
    const names: string[] = [];

    openUsersPage();
    cy.wait('@listUsers').then(({ response }) => {
      const body = response?.body as UserListBody;
      const pageCount = Math.ceil(body.meta.total / PAGE_SIZE);
      expect(pageCount, 'seeded data spans more than one page').to.be.greaterThan(1);

      for (let pageNr = 1; pageNr <= pageCount; pageNr++) {
        if (pageNr > 1) {
          cy.get('[data-cy="users-pagination"]').contains('button', String(pageNr)).click();
          cy.wait('@listUsers').its('response.body.meta.offset').should('eq', (pageNr - 1) * PAGE_SIZE);
        }
        const expectedRows = Math.min(PAGE_SIZE, body.meta.total - (pageNr - 1) * PAGE_SIZE);
        cy.get('[data-cy="user-email"]').should('have.length', expectedRows);
        cy.get('[data-cy="user-email"]').each(($el) => { emails.push($el.text().trim()); });
        cy.get('[data-cy="user-name"]').each(($el) => { names.push($el.text().trim()); });
      }

      cy.then(() => {
        expect(emails, 'every user listed').to.have.length(body.meta.total);
        expect(new Set(emails).size, 'no duplicates').to.eq(emails.length);
        const sorted = [...names].sort((a, b) => (asciiLower(a) < asciiLower(b) ? -1 : 1));
        expect(names.map(asciiLower), 'sorted by name').to.deep.equal(sorted.map(asciiLower));
      });
    });
  });

  it('should show a never-logged-in seeded user with "never"', () => {
    openUsersPage();
    cy.wait('@listUsers');

    cy.contains('[data-cy="user-email"]', 'e2e-pager-01@sithub.local')
      .closest('tr')
      .within(() => {
        cy.get('[data-cy="user-last-login"]').should('have.text', 'never');
        cy.get('[data-cy="user-roles"]').should('contain.text', 'User');
      });
  });

  it('should reject invalid pagination parameters', () => {
    cy.request({ url: '/api/v1/users?page[limit]=abc', failOnStatusCode: false }).then((response) => {
      expect(response.status).to.eq(400);
      expect(response.body.errors[0].detail).to.contain('page[limit]');
    });

    cy.request('/api/v1/users?page[limit]=500').then((response) => {
      const body = response.body as UserListBody;
      expect(body.data.length).to.be.at.most(100);
      expect(body.meta.limit).to.eq(100);
      expect(body.links.self).to.contain('page%5Blimit%5D=100');
    });
  });

  it('should hide the Users entry and deny the route for non-admins', () => {
    cy.login('alex@sithub.local', 'SitHubDemo2026!!');
    cy.visit('/');
    cy.get('[data-cy="user-menu-trigger"]').click();
    cy.get('[data-cy="avatar-btn"]').should('exist');
    cy.get('[data-cy="users-menu-btn"]').should('not.exist');

    cy.visit('/admin/users');
    cy.location('pathname').should('eq', '/access-denied');
  });
});
