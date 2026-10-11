/**
 * E2E coverage for the data export table (story 37.7).
 * Real backend, intercepts only wait for and inspect requests. Anna (admin) is granted the data exporter
 * role through the admin API for the tests that need it and loses it again afterwards.
 */

import { toLocalISODate as localISODate } from '../../src/utils/date';

const DEMO_PASSWORD = 'SitHubDemo2026!!';
const JSON_API = 'application/vnd.api+json';
const PAGE_SIZE = 20;

// Seeded items of the E2E areas config (run-all-tests.sh / CI): one in "Büro 2.EG", two in "Test Area".
const OWN_ITEM = 'cube_1_desk_1';
const COLLEAGUE_ITEM = 'desk_1';
const GUEST_ITEM = 'desk_2';
const GUEST_NAME = `E2E Guest ${Date.now()}`;

interface ExportBody {
  data: Array<{ id: string; attributes: { user_name: string; booked_item: string; item_id: string } }>;
  meta: { total: number };
}

/** Server-side "today" is the UTC date, so bookings are created for it. */
function utcToday() {
  return new Date().toISOString().slice(0, 10);
}

function setAnnaDataExporter(flag: boolean) {
  cy.login();
  cy.request('/api/v1/me').its('body.data.id').then((id) => {
    cy.request({
      method: 'PATCH',
      url: `/api/v1/users/${id}`,
      headers: { 'Content-Type': JSON_API },
      body: { data: { type: 'users', id, attributes: { is_data_exporter: flag } } }
    }).its('status').should('eq', 200);
  });
}

function createBooking(attributes: Record<string, unknown>) {
  return cy.request({
    method: 'POST',
    url: '/api/v1/bookings',
    headers: { 'Content-Type': JSON_API },
    body: { data: { type: 'bookings', attributes: { booking_date: utcToday(), ...attributes } } }
  }).then((response) => {
    expect(response.status).to.eq(201);
    return response.body.data.id as string;
  });
}

/** Cancels every booking of the given date so item availability does not depend on other specs. */
function clearBookingsOn(date: string): Cypress.Chainable<void> {
  // Always re-fetch offset 0: each DELETE shifts the remaining rows, so paging forward would skip some.
  return cy.request(`/api/v1/bookings/export?filter[from]=${date}&filter[to]=${date}&page[limit]=100&page[offset]=0`)
    .then((response) => {
      const bookings = (response.body as ExportBody).data;
      bookings.forEach((booking) => {
        cy.request({ method: 'DELETE', url: `/api/v1/bookings/${booking.id}` });
      });
      if (bookings.length > 0) {
        return clearBookingsOn(date);
      }
    });
}

function openDataExport() {
  cy.visit('/');
  cy.get('[data-cy="user-menu-trigger"]').click();
  cy.get('[data-cy="data-export-btn"]').click();
}

describe('data export', () => {
  const seededIds: string[] = [];

  before(() => {
    setAnnaDataExporter(true);
    clearBookingsOn(utcToday());
    cy.request('/api/v1/colleagues').then((response) => {
      const alex = (response.body.data as Array<{ id: string; attributes: { display_name: string } }>)
        .find((colleague) => colleague.attributes.display_name === 'Alex Employee');
      expect(alex, 'colleague Alex Employee').to.exist;

      createBooking({ item_id: OWN_ITEM }).then((id) => seededIds.push(id));
      createBooking({ item_id: COLLEAGUE_ITEM, for_user_id: alex!.id }).then((id) => seededIds.push(id));
      createBooking({
        item_id: GUEST_ITEM, is_guest: true, for_user_name: GUEST_NAME, guest_email: 'guest@example.com'
      }).then((id) => seededIds.push(id));
    });
  });

  after(() => {
    setAnnaDataExporter(true);
    seededIds.forEach((id) => cy.request({ method: 'DELETE', url: `/api/v1/bookings/${id}`, failOnStatusCode: false }));
    setAnnaDataExporter(false);
  });

  beforeEach(() => {
    cy.clearLocalStorage();
    cy.login();
    cy.intercept('GET', '/api/v1/bookings/export*').as('listExport');
  });

  afterEach(() => {
    // Restore the role even when a test that revoked it fails midway, so later tests do not cascade.
    setAnnaDataExporter(true);
  });

  it('should show the Data export entry to data exporters', () => {
    openDataExport();

    cy.location('pathname').should('eq', '/data-export');
    cy.get('[data-cy="data-export-page"]').should('be.visible');
  });

  it('should hide the entry and deny the route for an admin without the role', () => {
    setAnnaDataExporter(false);
    cy.visit('/');
    cy.get('[data-cy="user-menu-trigger"]').click();
    cy.get('[data-cy="users-menu-btn"]').should('exist');
    cy.get('[data-cy="data-export-btn"]').should('not.exist');

    cy.visit('/data-export');
    cy.location('pathname').should('eq', '/access-denied');
  });

  it('should default to the last 5 days and list the seeded bookings', () => {
    openDataExport();

    const today = new Date();
    const start = new Date(today);
    start.setDate(start.getDate() - 4);
    cy.wait('@listExport').then(({ request, response }) => {
      const url = new URL(request.url);
      expect(url.searchParams.get('filter[from]')).to.eq(localISODate(start));
      expect(url.searchParams.get('filter[to]')).to.eq(localISODate(today));
      expect(url.searchParams.get('page[limit]')).to.eq(String(PAGE_SIZE));
      expect(url.searchParams.get('page[offset]')).to.eq('0');
      expect(response?.statusCode).to.eq(200);
    });

    cy.contains('[data-cy="export-row"]', GUEST_NAME).should('contain.text', 'Test Room / Desk 2');
    cy.contains('[data-cy="export-row"]', 'Alex Employee').should('contain.text', 'Test Room / Desk 1');
    cy.contains('[data-cy="export-row"]', 'Anna Admin').should('contain.text', 'Cube 1 / Desk C1-1');
    cy.contains('[data-cy="export-row"]', 'Alex Employee').find('[data-cy="export-email"]')
      .should('have.text', 'alex@sithub.local');
  });

  it('should filter by area and restart from page 1', () => {
    openDataExport();
    cy.wait('@listExport');
    cy.contains('[data-cy="export-row"]', 'Alex Employee').should('exist');

    cy.get('[data-cy="export-area"]').click();
    cy.get('.v-overlay-container').contains('.v-list-item', 'Büro 2.EG').click();

    cy.wait('@listExport').then(({ request, response }) => {
      const url = new URL(request.url);
      expect(url.searchParams.get('filter[area]')).to.eq('buero_2eg');
      expect(url.searchParams.get('page[offset]')).to.eq('0');
      const body = response?.body as ExportBody;
      expect(body.data.every((booking) => booking.attributes.item_id.startsWith('cube_1'))).to.eq(true);
    });
    cy.contains('[data-cy="export-row"]', 'Anna Admin').should('exist');
    cy.contains('[data-cy="export-row"]', 'Alex Employee').should('not.exist');
    cy.contains('[data-cy="export-row"]', GUEST_NAME).should('not.exist');
  });

  it('should page with 20 rows and report the total', () => {
    openDataExport();

    cy.wait('@listExport').then(({ response }) => {
      const body = response?.body as ExportBody;
      expect(body.meta.total).to.be.at.least(3);
      expect(body.data.length).to.eq(Math.min(PAGE_SIZE, body.meta.total));
      cy.get('[data-cy="export-row"]').should('have.length', body.data.length);
    });
    cy.get('[data-cy="export-pagination"]').should('exist');
  });

  it('should reject invalid requests with a JSON:API error', () => {
    const today = utcToday();
    cy.request({
      url: `/api/v1/bookings/export?filter[from]=${today}&filter[to]=2000-01-01`,
      failOnStatusCode: false
    }).then((response) => {
      expect(response.status).to.eq(400);
      expect(response.body.errors[0].detail).to.eq('"filter[from]" must not be after "filter[to]"');
    });
    cy.request({ url: `/api/v1/bookings/export?filter[from]=${today}`, failOnStatusCode: false })
      .its('status').should('eq', 400);
    cy.request({ url: `/api/v1/bookings/export?filter[from]=${today}&filter[to]=${today}&filter[area]=nowhere`, failOnStatusCode: false })
      .its('body.errors.0.detail').should('eq', 'unknown area: nowhere');
  });

  it('should answer 403 for a user without the role', () => {
    cy.login('alex@sithub.local', DEMO_PASSWORD);
    cy.request({
      url: `/api/v1/bookings/export?filter[from]=${utcToday()}&filter[to]=${utcToday()}`,
      failOnStatusCode: false
    }).its('status').should('eq', 403);
  });
});
