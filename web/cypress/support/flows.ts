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

// Positions for the seeded "Büro 2.EG" / "Cube 1" floor plans (cypress/fixtures/floor-plans,
// seed areas in run-all-tests.sh and ci.yml). Created through the admin API when missing, so
// specs can run repeatedly against the same database. Logs in as the (admin) test user.
const SEEDED_FLOOR_PLAN_POSITIONS = [
  { floor_plan: 'e2e-area.svg', item_id: 'cube_1', label: 'Cube 1', x: 10, y: 13, width: 38, height: 34 },
  { floor_plan: 'e2e-room.svg', item_id: 'cube_1_desk_1', label: 'Desk C1-1', x: 12, y: 16, width: 20, height: 17 },
  { floor_plan: 'e2e-room.svg', item_id: 'cube_1_desk_2', label: 'Desk C1-2', x: 40, y: 16, width: 20, height: 17 }
];

export function seedFloorPlanPositions() {
  cy.login();
  SEEDED_FLOOR_PLAN_POSITIONS.forEach((position) => {
    cy.request(`/api/v1/floor-plan-positions?floor_plan=${position.floor_plan}`).then((response) => {
      const existing = (response.body.data as Array<{ attributes: { item_id: string } }>)
        .some((resource) => resource.attributes.item_id === position.item_id);
      if (!existing) {
        cy.request('POST', '/api/v1/floor-plan-positions', {
          data: { type: 'floor-plan-positions', attributes: { ...position, border_width: 2 } }
        });
      }
    });
  });
}

/**
 * Asserts that no word inside the matched elements is split across lines: every word's
 * Range must render as a single line box. Elements that are not ellipsised must also be at
 * least as wide as their widest word ("no element narrower than its longest word").
 */
export function assertNoMidWordBreak(selector: string) {
  cy.get(selector).should('have.length.at.least', 1).each(($el) => {
    const el = $el[0] as HTMLElement;
    const doc = el.ownerDocument;
    const walker = doc.createTreeWalker(el, NodeFilter.SHOW_TEXT);
    let widestWord = 0;
    for (let node = walker.nextNode(); node; node = walker.nextNode()) {
      for (const match of (node.textContent ?? '').matchAll(/\S+/g)) {
        const range = doc.createRange();
        range.setStart(node, match.index);
        range.setEnd(node, match.index + match[0].length);
        const lines = Array.from(range.getClientRects()).filter((rect) => rect.width > 0);
        expect(lines.length, `"${match[0]}" is rendered on one line`).to.eq(1);
        widestWord = Math.max(widestWord, lines[0]!.width);
      }
    }
    if (doc.defaultView!.getComputedStyle(el).textOverflow !== 'ellipsis') {
      expect(el.getBoundingClientRect().width + 0.5, `"${el.textContent?.trim()}" fits its widest word`)
        .to.be.at.least(widestWord);
    }
  });
}

/** Asserts that the text of each matched element is laid out on a single line. */
export function assertSingleLine(selector: string) {
  cy.get(selector).should('have.length.at.least', 1).each(($el) => {
    const range = $el[0]!.ownerDocument.createRange();
    range.selectNodeContents($el[0]!);
    const lineTops = new Set(
      Array.from(range.getClientRects()).filter((rect) => rect.width > 0).map((rect) => Math.round(rect.top))
    );
    expect(lineTops.size, `"${$el.text().trim()}" is rendered on one line`).to.eq(1);
  });
}
