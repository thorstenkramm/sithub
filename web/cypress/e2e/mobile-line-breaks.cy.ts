// Story 37.3: at 360px no label is split inside a word (German UI, longest labels).
import { assertNoMidWordBreak, assertSingleLine, resetAndLogin, seedFloorPlanPositions } from '../support/flows';

const AREA_ID = 'buero_2eg';
const ITEM_GROUP_ID = 'cube_1';

// cypress/support/e2e.ts forces English before every load; this runs afterwards and wins.
const germanUi = {
  onBeforeLoad(win: Cypress.AUTWindow) {
    win.localStorage.setItem('sithub_locale', 'de');
  }
};

const BREADCRUMB_LABELS = '[data-cy^="breadcrumb-link-"], [data-cy^="breadcrumb-text-"]';
const FOOTER_LABELS = [
  '[data-cy="fp-show-labels"] label',
  '[data-cy="fp-show-avatars"] label',
  '[data-cy="fp-close-btn"]'
].join(', ');

// Each crumb entry stays whole: no break inside a word and none between its words.
function expectWholeBreadcrumbEntries() {
  assertNoMidWordBreak(BREADCRUMB_LABELS);
  assertSingleLine(BREADCRUMB_LABELS);
}

function expectFooterInsideViewport() {
  cy.get('[data-cy="fp-footer"]').then(($footer) => {
    const footer = $footer[0]!;
    expect(footer.scrollWidth, 'footer does not overflow').to.be.at.most(footer.clientWidth);
  });
  cy.get('[data-cy="fp-close-btn"]').should('be.visible').then(($btn) => {
    expect($btn[0]!.getBoundingClientRect().right).to.be.at.most(360);
  });
}

describe('mobile line breaks', () => {
  before(() => seedFloorPlanPositions());

  beforeEach(() => {
    resetAndLogin();
    cy.viewport(360, 740);
    cy.intercept('GET', /\/api\/v1\/floor-plan-positions/).as('positions');
  });

  it('should keep breadcrumb entries on word boundaries at 360px', () => {
    cy.intercept('GET', `/api/v1/item-groups/${ITEM_GROUP_ID}/bookings*`).as('itemGroupBookings');
    cy.visit(`/item-groups/${ITEM_GROUP_ID}/bookings?areaId=${AREA_ID}`, germanUi);
    cy.wait('@itemGroupBookings');

    cy.get('[data-cy="breadcrumbs"]').should('contain', 'Startseite').and('contain', 'Büro 2.EG')
      .and('contain', 'Cube 1');
    expectWholeBreadcrumbEntries();
  });

  it('should keep breadcrumb entries on word boundaries on the items and presence pages', () => {
    cy.intercept('GET', `/api/v1/item-groups/${ITEM_GROUP_ID}/items*`).as('listItems');
    cy.visit(`/item-groups/${ITEM_GROUP_ID}/items?areaId=${AREA_ID}`, germanUi);
    cy.wait('@listItems');
    cy.get('[data-cy="breadcrumbs"]').should('contain', 'Büro 2.EG');
    expectWholeBreadcrumbEntries();

    cy.intercept('GET', `/api/v1/areas/${AREA_ID}/presence*`).as('presence');
    cy.visit(`/areas/${AREA_ID}/presence`, germanUi);
    cy.wait('@presence');
    cy.get('[data-cy="breadcrumbs"]').should('contain', 'Büro 2.EG');
    expectWholeBreadcrumbEntries();
  });

  it('should keep floor-plan footer labels on word boundaries at 360px', () => {
    cy.visit(`/item-groups/${ITEM_GROUP_ID}/items?areaId=${AREA_ID}`, germanUi);
    cy.get('[data-cy="item-group-floor-plan-btn"]').click();
    cy.wait('@positions');

    cy.get('[data-cy="item-group-floor-plan-dialog"]').within(() => {
      cy.get('[data-cy="fp-show-labels"]').should('contain', 'Beschriftungen anzeigen');
      assertNoMidWordBreak(FOOTER_LABELS);
      expectFooterInsideViewport();
    });
  });

  it('should keep the area floor-plan footer and the view switch labels on word boundaries at 360px', () => {
    cy.visit(`/areas/${AREA_ID}/item-groups`, germanUi);
    cy.get('[data-cy="view-switch-container"]').should('be.visible');
    assertNoMidWordBreak('[data-cy="view-switch-container"] span');

    cy.get('[data-cy="area-floor-plan-btn"]').click();
    cy.wait('@positions');

    cy.get('[data-cy="floor-plan-dialog"]').within(() => {
      cy.get('[data-cy="floor-plan-area-drill-down-toggle"]').should('contain', 'Bereichs-Drilldown');
      assertNoMidWordBreak(`${FOOTER_LABELS}, [data-cy="floor-plan-area-drill-down-toggle"] label`);
      expectFooterInsideViewport();
    });
  });
});
