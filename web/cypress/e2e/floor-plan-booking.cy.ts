import {
  confirmWarningAndExpectBooking,
  interceptBookingSuccess,
  interceptWarnedItemsList,
  resetAndLogin
} from '../support/flows';

/**
 * The warning confirmation on the interactive floor plan: clicking a free warned
 * item opens the date dialog, and confirming it must surface the shared warning
 * dialog before any booking is created (FR159/FR164). Responses are intercepted
 * (matching the matrix-booking spec) so a warned, free item exists deterministically.
 */
function setupFloorPlanIntercepts() {
  cy.intercept('GET', '/api/v1/areas', {
    statusCode: 200,
    body: { data: [{ id: 'test_area', type: 'areas', attributes: { name: 'Test Area' } }] }
  }).as('listAreas');
  cy.intercept('GET', /\/api\/v1\/areas\/[^/]+\/item-groups(\?.*)?$/, {
    statusCode: 200,
    body: { data: [{ id: 'test_room', type: 'item-groups', attributes: { name: 'Test Room', floor_plan: 'room.svg' } }] }
  }).as('listItemGroups');
  cy.intercept('GET', /\/api\/v1\/floor-plan-positions/, {
    statusCode: 200,
    body: {
      data: [{
        id: 'pos-1',
        type: 'floor-plan-positions',
        attributes: {
          floor_plan: 'room.svg', item_id: 'item-warned-1', label: 'Window Desk',
          x: 20, y: 20, width: 30, height: 30, border_width: 2
        }
      }]
    }
  }).as('positions');
  interceptWarnedItemsList();
  cy.intercept('GET', '/api/v1/floor-plans/*', {
    statusCode: 200,
    headers: { 'Content-Type': 'image/svg+xml' },
    body: '<svg xmlns="http://www.w3.org/2000/svg" width="800" height="600"></svg>'
  }).as('floorPlanImage');
}

function openFloorPlanAndSelectWarnedItem() {
  // Freeze the date to a weekday so "today" is bookable and pre-selected.
  cy.clock(new Date('2026-06-01T10:00:00').getTime(), ['Date']);
  cy.visit('/item-groups/test_room/items');
  cy.wait('@listItems');

  cy.get('[data-cy="item-group-floor-plan-btn"]').click();
  cy.wait('@positions');

  cy.get('[data-cy="fp-item-item-warned-1"]').click();
  cy.get('[data-cy="fp-booking-dialog"]').should('be.visible');
  cy.get('[data-cy="fp-booking-confirm"]').click();
}

describe('floor plan warning confirmation', () => {
  beforeEach(() => {
    resetAndLogin();
  });

  it('should surface the warning confirmation before booking a warned item and book on confirm', () => {
    setupFloorPlanIntercepts();
    interceptBookingSuccess();

    openFloorPlanAndSelectWarnedItem();

    // The shared warning dialog appears before any booking is created.
    confirmWarningAndExpectBooking();
  });

  it('should abort the floor-plan booking when the warning confirmation is cancelled', () => {
    setupFloorPlanIntercepts();
    cy.intercept('POST', '/api/v1/bookings', { statusCode: 201, body: {} }).as('createBooking');

    openFloorPlanAndSelectWarnedItem();

    cy.get('[data-cy="warning-dialog"]').should('be.visible');
    cy.get('[data-cy="warning-cancel-btn"]').click();

    cy.get('[data-cy="warning-dialog"]').should('not.exist');
    cy.get('@createBooking.all').should('have.length', 0);
  });
});
