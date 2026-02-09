import {
  createMockItem,
  createMockItemsResponse,
  setupItemsPageIntercepts
} from '../support/flows';

/**
 * UI Framework Integration Tests
 *
 * These tests verify that the UI framework (Vuetify) is properly initialized
 * and components render correctly. This catches issues like:
 * - Missing Vuetify component/directive imports
 * - CSS not loading
 * - JavaScript bundle errors
 */

describe('UI Framework Integration', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
    cy.login();
  });

  it('should render Vuetify components (not raw HTML tags)', () => {
    cy.visit('/');

    // Wait for the app to load
    cy.get('#app').should('exist');

    // The v-application wrapper should be rendered by Vuetify
    // If Vuetify isn't initialized, we'd see <v-app> as a raw tag
    cy.get('.v-application').should('exist');

    // Verify the app bar renders as a proper Vuetify component
    cy.get('.v-app-bar').should('exist');
    cy.get('.v-toolbar').should('exist');

    // Verify raw Vuetify tags are NOT in the DOM (they should be rendered as proper components)
    cy.document().then((doc) => {
      const rawVuetifyTags = doc.querySelectorAll(
        'v-app, v-app-bar, v-card, v-btn, v-container, v-row, v-col'
      );
      expect(rawVuetifyTags.length, 'Raw Vuetify tags should not exist in rendered DOM').to.equal(0);
    });
  });

  it('should load Vuetify CSS styles', () => {
    cy.visit('/');

    // Check that Vuetify CSS variables are defined
    cy.document().then((doc) => {
      const styles = getComputedStyle(doc.documentElement);
      // Vuetify defines these CSS custom properties
      const primaryColor = styles.getPropertyValue('--v-theme-primary');
      expect(primaryColor, 'Vuetify theme primary color should be defined').to.not.be.empty;
    });

    // Verify the app bar has proper background color (primary color)
    cy.get('.v-app-bar').should('have.css', 'background-color');
  });

  it('should render navigation with proper Vuetify button components', () => {
    cy.visit('/');

    // Navigation buttons should be rendered as proper Vuetify buttons
    cy.get('.v-app-bar').within(() => {
      // Check for Vuetify button classes
      cy.get('.v-btn').should('have.length.at.least', 1);
    });
  });

  it('should render cards with proper Vuetify card components', () => {
    cy.visit('/');

    // Wait for areas to load
    cy.get('[data-cy="areas-list"], [data-cy="areas-loading"], [data-cy="areas-empty"]', {
      timeout: 10000
    }).should('exist');

    // If areas exist, verify cards render properly
    cy.get('body').then(($body) => {
      if ($body.find('[data-cy="areas-list"]').length > 0) {
        // In redesigned UI, area-item IS the v-card
        cy.get('[data-cy="area-item"]').first().should('have.class', 'v-card');
        // Should have Vuetify card structure inside
        cy.get('[data-cy="area-item"]').first().within(() => {
          cy.get('.v-card-item, .v-card-title, .v-card-text, .v-card-actions')
            .should('have.length.at.least', 1);
        });
      }
    });
  });

  it('should render icons properly (not as text)', () => {
    cy.visit('/');

    // Wait for content to load
    cy.get('.v-application').should('exist');

    // Icons should render as SVG, not as text like "$area" or "mdi-office-building"
    cy.get('.v-icon').should('exist');
    cy.get('.v-icon').first().within(() => {
      cy.get('svg').should('exist');
    });

    // Verify no icon placeholder text is visible
    cy.get('body').should('not.contain.text', '$area');
    cy.get('body').should('not.contain.text', '$room');
    cy.get('body').should('not.contain.text', '$desk');
    cy.get('body').should('not.contain.text', 'mdi-');
  });

  it('should have proper theme colors applied', () => {
    cy.visit('/');

    // App bar should have the primary color background
    cy.get('.v-app-bar').then(($appBar) => {
      const bgColor = $appBar.css('background-color');
      // Primary blue should be applied (not default or transparent)
      expect(bgColor).to.not.equal('rgba(0, 0, 0, 0)');
      expect(bgColor).to.not.equal('transparent');
    });
  });

  it('should render loading states with skeleton loaders', () => {
    // Intercept and delay the API response to see loading state
    cy.intercept('GET', '/api/v1/areas', (req) => {
      req.on('response', (res) => {
        res.setDelay(500);
      });
    }).as('listAreas');

    cy.visit('/');

    // Should show Vuetify skeleton loader, not raw HTML
    cy.get('[data-cy="areas-loading"]').should('exist');
    cy.get('.v-skeleton-loader').should('exist');

    // Wait for data to load
    cy.wait('@listAreas');
  });

  it('should render empty states with proper styling', () => {
    // Mock empty response
    cy.intercept('GET', '/api/v1/areas', {
      statusCode: 200,
      body: { data: [] }
    }).as('listAreasEmpty');

    cy.visit('/');
    cy.wait('@listAreasEmpty');

    // Empty state should render with proper Vuetify components
    cy.get('[data-cy="areas-empty"]').should('exist');
    cy.get('[data-cy="areas-empty"]').within(() => {
      // Should have an icon
      cy.get('.v-icon').should('exist');
    });
  });

  it('should have responsive navigation (mobile drawer)', () => {
    // Set mobile viewport
    cy.viewport(375, 812);
    cy.visit('/');

    // Desktop nav should be hidden
    cy.get('.v-app-bar').within(() => {
      // Mobile menu button should exist
      cy.get('.v-app-bar-nav-icon, [data-cy="mobile-menu-btn"]').should('be.visible');
    });

    // Desktop navigation links should be hidden on mobile
    cy.get('.d-none.d-md-flex').should('not.be.visible');
  });
});

describe('Page-specific Vuetify Rendering', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
    cy.login();
  });

  it('should render date picker as Vuetify component on items page', () => {
    setupItemsPageIntercepts();

    cy.visit('/item-groups/test_room/items');

    // Date picker should be a proper Vuetify text field
    cy.get('[data-cy="items-date"]').should('exist');
    cy.get('.v-text-field').should('exist');

    // Should have a calendar icon
    cy.get('.v-field__prepend-inner .v-icon').should('exist');
  });

  it('should render radio buttons as Vuetify components on items page', () => {
    cy.intercept('GET', '/api/v1/item-groups/*/items*').as('listItems');

    cy.visit('/item-groups/test_room/items');

    // Radio group should render as Vuetify component
    cy.get('.v-radio-group').should('exist');
    cy.get('.v-radio').should('have.length.at.least', 3);

    // Labels should be visible
    cy.contains('Book for myself').should('be.visible');
    cy.contains('Book for colleague').should('be.visible');
    cy.contains('Book for guest').should('be.visible');
  });

  it('should render status chips as Vuetify components', () => {
    const mockItem = createMockItem('item-1', 'Test Item');
    cy.intercept('GET', '/api/v1/item-groups/*/items*', createMockItemsResponse([mockItem])).as(
      'listItems'
    );

    cy.visit('/item-groups/test_room/items');
    cy.wait('@listItems');

    // Status chip should render as Vuetify chip
    cy.get('[data-cy="item-status"]').should('exist');
    cy.get('[data-cy="item-status"]').should('have.class', 'v-chip');
  });

  it('should render alerts as Vuetify components', () => {
    // Mock an error response to trigger alert
    cy.intercept('GET', '/api/v1/areas', {
      statusCode: 500,
      body: { errors: [{ status: '500', title: 'Server Error' }] }
    }).as('listAreasError');

    cy.visit('/');
    cy.wait('@listAreasError');

    // Error alert should render as Vuetify alert
    cy.get('.v-alert').should('exist');
    cy.get('[data-cy="areas-error"]').should('exist');
  });
});
