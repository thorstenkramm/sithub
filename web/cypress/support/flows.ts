export function openArea(areaName: string) {
  cy.wait('@listAreas').its('response.statusCode').should('eq', 200);
  cy.contains('[data-cy="area-item"]', areaName).click();
}

export function openItemGroup(itemGroupName: string) {
  cy.wait('@listItemGroups').its('response.statusCode').should('eq', 200);
  cy.contains('[data-cy="item-group-item"]', itemGroupName).click();
}

export function resetAndLogin() {
  cy.clearCookies();
  cy.clearLocalStorage();
  cy.login();
}

export function visitHomeAndWaitForAreas() {
  cy.visit('/');
  cy.waitForVuetify();
  cy.wait('@listAreas').its('response.statusCode').should('eq', 200);
}

export function openFirstAreaAndWaitItemGroups() {
  cy.get('[data-cy="area-item"]').first().click();
  cy.wait('@listItemGroups').its('response.statusCode').should('eq', 200);
  cy.location('pathname').should('match', /\/areas\/.*\/item-groups/);
}

export function openFirstItemGroupAndWaitItems() {
  cy.get('[data-cy="item-group-item"]').first().click();
  cy.wait('@listItems').its('response.statusCode').should('eq', 200);
  cy.location('pathname').should('match', /\/item-groups\/.*\/items/);
}

/** Sets up intercepts for areas, item groups, items, and availability API endpoints */
export function setupItemsPageIntercepts() {
  cy.intercept('GET', '/api/v1/areas').as('listAreas');
  cy.intercept('GET', /\/api\/v1\/areas\/[^/]+\/item-groups\/availability/).as('availability');
  cy.intercept('GET', /\/api\/v1\/areas\/[^/]+\/item-groups$/).as('listItemGroups');
  cy.intercept('GET', '/api/v1/item-groups/*/items*').as('listItems');
  cy.intercept('GET', /\/api\/v1\/item-groups\/[^/]+\/items\?date=/).as('listItemsWithDate');
}

/** Creates a mock item object for testing */
export function createMockItem(
  id: string,
  name: string,
  availability: 'available' | 'occupied' = 'available',
  equipment: string[] = ['Monitor'],
  bookerName?: string
) {
  return {
    id,
    type: 'items',
    attributes: {
      name,
      equipment,
      availability,
      ...(bookerName ? { booker_name: bookerName } : {})
    }
  };
}

/** Creates a mock items API response body */
export function createMockItemsResponse(items: ReturnType<typeof createMockItem>[]) {
  return {
    statusCode: 200,
    body: {
      data: items
    }
  };
}

/**
 * Intercepts POST /api/v1/bookings with a 409 Conflict carrying the given detail
 * message, aliased so the test can cy.wait on it.
 */
export function interceptBookingConflict(detail: string, alias: string) {
  cy.intercept('POST', '/api/v1/bookings', {
    statusCode: 409,
    headers: { 'Content-Type': 'application/vnd.api+json' },
    body: {
      errors: [
        {
          status: '409',
          title: 'Conflict',
          detail,
          code: 'conflict'
        }
      ]
    }
  }).as(alias);
}

/** A free item carrying a warning, for exercising the pre-booking warning confirmation. */
export const WARNED_ITEM = {
  id: 'item-warned-1',
  type: 'items',
  attributes: {
    name: 'Window Desk',
    equipment: ['Monitor'],
    availability: 'available',
    booked_by_me: false,
    warning: 'Near noisy area'
  }
};

/** Intercepts the items list to return a single free, warned item, aliased `@listItems`. */
export function interceptWarnedItemsList() {
  cy.intercept('GET', '/api/v1/item-groups/*/items*', {
    statusCode: 200,
    body: { data: [WARNED_ITEM] }
  }).as('listItems');
}

/** Stubs a successful booking (201) aliased `@createBooking`. */
export function interceptBookingSuccess(itemId = 'item-warned-1') {
  cy.intercept('POST', '/api/v1/bookings', {
    statusCode: 201,
    headers: { 'Content-Type': 'application/vnd.api+json' },
    body: { data: { id: 'booking-1', type: 'bookings', attributes: { item_id: itemId } } }
  }).as('createBooking');
}

/**
 * Asserts the shared warning dialog is shown (with no booking created yet) and
 * confirms it, expecting the booking POST to fire for `itemId`.
 */
export function confirmWarningAndExpectBooking(itemId = 'item-warned-1') {
  cy.get('[data-cy="warning-dialog"]').should('be.visible');
  cy.get('[data-cy="warning-message"]').should('contain', 'Near noisy area');
  cy.get('@createBooking.all').should('have.length', 0);
  cy.get('[data-cy="warning-confirm-btn"]').click();
  cy.wait('@createBooking').then((interception) => {
    expect(interception.response?.statusCode).to.eq(201);
    expect(interception.request.body.data.attributes.item_id).to.eq(itemId);
  });
}
