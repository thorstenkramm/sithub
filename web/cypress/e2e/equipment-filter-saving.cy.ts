describe('equipment filter saving', () => {
  const SAVED_FILTERS_KEY = 'sithub_saved_filters';

  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
    cy.login();
  });

  /**
   * Helper: extract the SVG path data from the save/delete v-icon inside
   * a `data-cy`-targeted button. A registered Vuetify icon alias resolves to
   * a real path (>= 50 chars of M/L/C/A commands). An unregistered/literal
   * icon name like `mdi-content-save` would land in the `d` attribute as the
   * literal string (~16 chars) — invisible to the user. This helper lets
   * tests assert the icon actually renders.
   */
  const expectIconPathRenders = (selector: string) => {
    cy.get(`[data-cy="${selector}"] .v-icon__svg path`)
      .invoke('attr', 'd')
      .should((value) => {
        expect(value, `<${selector}> icon path data`).to.be.a('string');
        expect(value!.length, `<${selector}> icon path length`).to.be.greaterThan(40);
        expect(value!, `<${selector}> icon path content`).not.to.equal('mdi-content-save');
      });
  };

  it('shows a save icon (registered alias) when typing a filter on the item-groups page', () => {
    cy.visit('/');
    cy.get('[data-cy="area-item"]').first().click();

    // Filter is empty — save action is in the DOM but invisible (placeholder).
    cy.get('[data-cy="ig-equipment-filter-save"]')
      .should('have.class', 'filter-action-placeholder');

    cy.get('[data-cy="ig-equipment-filter"] input').type('webcam').blur();

    // After typing + blurring, the save button becomes visible and renders
    // a real SVG icon (NOT the literal string `mdi-content-save`).
    cy.get('[data-cy="ig-equipment-filter-save"]')
      .should('not.have.class', 'filter-action-placeholder');
    expectIconPathRenders('ig-equipment-filter-save');
  });

  it('persists a saved filter to localStorage and toggles to a delete icon', () => {
    cy.visit('/');
    cy.get('[data-cy="area-item"]').first().click();

    cy.get('[data-cy="ig-equipment-filter"] input').type('webcam').blur();
    cy.get('[data-cy="ig-equipment-filter-save"]').click();

    // localStorage gets the persisted filter.
    cy.window().its('localStorage')
      .invoke('getItem', SAVED_FILTERS_KEY)
      .should('eq', JSON.stringify(['webcam']));

    // The icon toggles from save → delete and the delete icon also renders.
    cy.get('[data-cy="ig-equipment-filter-delete"]').should('exist');
    expectIconPathRenders('ig-equipment-filter-delete');
  });

  it('survives a page reload and removes the filter on delete-click', () => {
    cy.visit('/', {
      onBeforeLoad(win) {
        win.localStorage.setItem(SAVED_FILTERS_KEY, JSON.stringify(['webcam']));
      }
    });
    cy.get('[data-cy="area-item"]').first().click();

    // Re-typing the same filter recognizes it as already saved → delete icon.
    cy.get('[data-cy="ig-equipment-filter"] input').type('webcam').blur();
    cy.get('[data-cy="ig-equipment-filter-delete"]').should('exist');
    expectIconPathRenders('ig-equipment-filter-delete');

    // Deleting clears the input and removes the entry from localStorage.
    cy.get('[data-cy="ig-equipment-filter-delete"]').click();
    cy.window().its('localStorage')
      .invoke('getItem', SAVED_FILTERS_KEY)
      .should('eq', JSON.stringify([]));
  });

  it('shows a save icon (registered alias) when typing a filter on the items page', () => {
    cy.visit('/');
    cy.get('[data-cy="area-item"]').first().click();
    cy.get('[data-cy="item-group-item"]').first().click();

    cy.get('[data-cy="equipment-filter-save"]')
      .should('have.class', 'filter-action-placeholder');

    cy.get('[data-cy="equipment-filter-input"] input').type('webcam').blur();

    cy.get('[data-cy="equipment-filter-save"]')
      .should('not.have.class', 'filter-action-placeholder');
    expectIconPathRenders('equipment-filter-save');
  });
});
