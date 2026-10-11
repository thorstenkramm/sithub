// Story 37.2: week dropdown plus previous/next arrows in the tile views and the floor plan.
import { getISOWeekNumber, getISOWeekString, getMondayOfWeek } from '../../src/composables/useWeekSelector';
import { resetAndLogin, seedFloorPlanPositions } from '../support/flows';
import { toLocalISODate as isoDate } from '../../src/utils/date';

const AREA_ID = 'buero_2eg';
const ITEM_GROUP_ID = 'cube_1';

// Frontend "now" is frozen to Monday 10:00 of the real current week, so the week options
// start on a known weekday while staying in sync with the backend's booking horizon.
const currentMonday = getMondayOfWeek(new Date());
const frozenNow = new Date(currentMonday).setHours(10, 0, 0, 0);

const addDays = (date: Date, days: number): Date => {
  const result = new Date(date);
  result.setDate(date.getDate() + days);
  return result;
};

const weekAt = (offset: number): string => getISOWeekString(addDays(currentMonday, offset * 7));
const weekLabelAt = (offset: number): string => `Week ${getISOWeekNumber(addDays(currentMonday, offset * 7))}`;

// Loads weeks_in_advanced from the real settings response the page requested.
function waitForWeeksInAdvanced(): Cypress.Chainable<number> {
  return cy.wait('@settings').then((interception) => {
    return interception.response?.body.data.attributes.weeks_in_advanced as number;
  });
}

// Registers an intercept for the NEXT request with the given query, so earlier
// (initial-load) requests on the shared alias cannot satisfy the wait.
function nextAvailabilityFor(week: string): string {
  cy.intercept({ method: 'GET', pathname: `/api/v1/areas/${AREA_ID}/item-groups/availability`, query: { week } })
    .as(`availability-${week}`);
  return `@availability-${week}`;
}

function nextItemsFor(date: Date): string {
  const day = isoDate(date);
  cy.intercept({ method: 'GET', pathname: `/api/v1/item-groups/${ITEM_GROUP_ID}/items`, query: { date: day } })
    .as(`items-${day}`);
  return `@items-${day}`;
}

function expectBounds(prevDisabled: boolean, nextDisabled: boolean) {
  cy.get('[data-cy="week-prev"]').should(prevDisabled ? 'be.disabled' : 'not.be.disabled');
  cy.get('[data-cy="week-next"]').should(nextDisabled ? 'be.disabled' : 'not.be.disabled');
}

// Clicks "next week" until the last bookable week, waiting for the request each step
// produces and checking the dropdown label, then steps back to the current week.
function stepThroughAllWeeks(weeks: number, nextRequestFor: (offset: number) => string) {
  expectBounds(true, false);
  Cypress._.times(weeks, (i) => {
    const request = nextRequestFor(i + 1);
    cy.get('[data-cy="week-next"]').click();
    cy.wait(request);
    cy.get('[data-cy="week-selector"]').should('contain', weekLabelAt(i + 1));
  });
  expectBounds(false, true);

  Cypress._.times(weeks, () => cy.get('[data-cy="week-prev"]').click());
  cy.get('[data-cy="week-selector"]').should('contain', weekLabelAt(0));
  expectBounds(true, false);
}

// In an open floor plan: selects Wednesday, then steps to the next week and waits until
// availability is loaded for Wednesday of that week (the weekday is kept).
function selectWednesdayAndStepToNextWeek() {
  const wednesday = nextItemsFor(addDays(currentMonday, 2));
  cy.get('[data-cy="fp-day-WE"]').click();
  cy.wait(wednesday);

  const nextWednesday = nextItemsFor(addDays(currentMonday, 9));
  cy.get('[data-cy="week-next"]').click();
  cy.wait(nextWednesday);
  cy.get('[data-cy="fp-day-WE"]').should('have.class', 'v-btn--variant-flat');
}

function visitItemsPage(mode: 'day' | 'week', locale = 'en') {
  cy.visit(`/item-groups/${ITEM_GROUP_ID}/items`, {
    onBeforeLoad(win) {
      win.localStorage.setItem('sithub_booking_mode', mode);
      win.localStorage.setItem('sithub_locale', locale);
      win.sessionStorage.setItem('sithub_selected_day', isoDate(currentMonday));
    }
  });
}

describe('week navigation', () => {
  before(() => seedFloorPlanPositions());

  beforeEach(() => {
    resetAndLogin();
    cy.clock(frozenNow, ['Date']);
    cy.intercept('GET', '/api/v1/settings').as('settings');
    cy.intercept('GET', '/api/v1/item-groups/*/items*').as('listItems');
  });

  it('should step through all bookable weeks in the area tile view', () => {
    cy.visit(`/areas/${AREA_ID}/item-groups`);
    waitForWeeksInAdvanced().then((weeks) => {
      cy.get('[data-cy="week-selector-card"]').within(() => {
        stepThroughAllWeeks(weeks, (offset) => nextAvailabilityFor(weekAt(offset)));
      });
    });
  });

  it('should step through all bookable weeks in the item tile view week mode', () => {
    visitItemsPage('week');
    waitForWeeksInAdvanced().then((weeks) => {
      // Week mode loads every weekday; Monday of the new week is one of them.
      stepThroughAllWeeks(weeks, (offset) => nextItemsFor(addDays(currentMonday, offset * 7)));
    });
  });

  it('should step weeks in the floor plan opened from the item view in day mode and keep the weekday', () => {
    visitItemsPage('day');
    waitForWeeksInAdvanced().then((weeks) => {
      cy.wait('@listItems');
      cy.get('[data-cy="item-group-floor-plan-btn"]').click();

      cy.get('[data-cy="item-group-floor-plan-dialog"]').within(() => {
        cy.get('[data-cy="fp-week-selector"]').should('exist');
        cy.get('[data-cy="fp-week-label"]').should('contain', 'Week');
        expectBounds(true, false);
        selectWednesdayAndStepToNextWeek();

        Cypress._.times(weeks - 1, () => cy.get('[data-cy="week-next"]').click());
        expectBounds(false, true);
        Cypress._.times(weeks, () => cy.get('[data-cy="week-prev"]').click());
        expectBounds(true, false);
      });
    });
  });

  it('should step weeks in the floor plan opened from the item view in week mode', () => {
    visitItemsPage('week');
    waitForWeeksInAdvanced();
    cy.get('[data-cy="item-group-floor-plan-btn"]').click();

    cy.get('[data-cy="item-group-floor-plan-dialog"]').within(() => {
      cy.get('[data-cy="fp-week-selector"]').should('exist');
      expectBounds(true, false);
      selectWednesdayAndStepToNextWeek();
      expectBounds(false, false);
    });
  });

  it('should keep the drill-down and the weekday when stepping weeks in the area floor plan', () => {
    cy.visit(`/areas/${AREA_ID}/item-groups`);
    waitForWeeksInAdvanced();
    cy.get('[data-cy="area-floor-plan-btn"]').click();

    cy.get('[data-cy="floor-plan-dialog"]').within(() => {
      cy.get('[data-cy="fp-week-selector"]').should('exist');
      expectBounds(true, false);

      cy.get(`[data-cy="fp-area-${ITEM_GROUP_ID}"]`).click();
      cy.get('[data-cy="fp-breadcrumb-current"]').should('contain', 'Cube 1');

      // The drilled item group's items are reloaded for the same weekday one week later.
      selectWednesdayAndStepToNextWeek();
      cy.get('[data-cy="fp-breadcrumb-current"]').should('contain', 'Cube 1');
    });
  });

  it('should fit the floor plan week selector into a 360px wide screen', () => {
    cy.viewport(360, 740);
    visitItemsPage('day');
    waitForWeeksInAdvanced();
    cy.get('[data-cy="item-group-floor-plan-btn"]').click();

    cy.get('[data-cy="fp-header"]').then(($header) => {
      const header = $header[0]!;
      expect(header.scrollWidth, 'header does not overflow').to.be.at.most(header.clientWidth);
    });
    cy.get('[data-cy="week-next"]').should('be.visible').then(($next) => {
      expect($next[0]!.getBoundingClientRect().right).to.be.at.most(360);
    });
  });

  it('should label the floor plan week with "KW" in German', () => {
    visitItemsPage('day', 'de');
    waitForWeeksInAdvanced();
    cy.get('[data-cy="item-group-floor-plan-btn"]').click();

    cy.get('[data-cy="item-group-floor-plan-dialog"]').within(() => {
      cy.get('[data-cy="fp-week-label"]').should('contain', 'KW').and('not.contain', 'Week');
      cy.get('[data-cy="week-prev"]').should('have.attr', 'aria-label', 'Vorherige Woche');
      cy.get('[data-cy="week-next"]').should('have.attr', 'aria-label', 'Nächste Woche');
    });
  });
});
