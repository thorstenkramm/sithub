/// <reference types="cypress" />

/**
 * Custom Cypress commands for SitHub E2E testing
 */

declare global {
  namespace Cypress {
    interface Chainable {
      /**
       * Verify that Vuetify is properly initialized and rendering components
       */
      verifyVuetifyLoaded(): Chainable<void>;

      /**
       * Verify that no raw Vuetify tags exist in the DOM
       */
      verifyNoRawVuetifyTags(): Chainable<void>;

      /**
       * Wait for the page to be fully loaded with Vuetify
       */
      waitForVuetify(): Chainable<void>;
    }
  }
}

/**
 * Verify that Vuetify framework is properly loaded and initialized
 */
Cypress.Commands.add('verifyVuetifyLoaded', () => {
  // The v-application class is added by Vuetify when properly initialized
  cy.get('.v-application', { timeout: 10000 }).should('exist');

  // Verify Vuetify CSS is loaded by checking for theme CSS variables
  cy.document().then((doc) => {
    const styles = getComputedStyle(doc.documentElement);
    const primaryColor = styles.getPropertyValue('--v-theme-primary');
    expect(primaryColor, 'Vuetify theme CSS should be loaded').to.not.be.empty;
  });
});

/**
 * Verify that no raw/unrendered Vuetify component tags exist in the DOM
 * If Vuetify components aren't imported, they render as raw HTML tags like <v-card>
 */
Cypress.Commands.add('verifyNoRawVuetifyTags', () => {
  cy.document().then((doc) => {
    // Common Vuetify component tags that should never appear raw in the DOM
    const rawTagSelectors = [
      'v-app',
      'v-app-bar',
      'v-toolbar',
      'v-container',
      'v-row',
      'v-col',
      'v-card',
      'v-card-title',
      'v-card-text',
      'v-card-actions',
      'v-btn',
      'v-icon',
      'v-list',
      'v-list-item',
      'v-alert',
      'v-chip',
      'v-text-field',
      'v-select',
      'v-dialog',
      'v-menu',
      'v-avatar',
      'v-spacer',
      'v-divider',
      'v-radio',
      'v-radio-group',
      'v-checkbox',
      'v-progress-linear',
      'v-skeleton-loader',
      'v-navigation-drawer'
    ];

    const rawTags = doc.querySelectorAll(rawTagSelectors.join(', '));

    if (rawTags.length > 0) {
      const tagNames = Array.from(rawTags).map((el) => el.tagName.toLowerCase());
      const uniqueTags = [...new Set(tagNames)];
      throw new Error(
        `Found ${rawTags.length} raw Vuetify tag(s) in DOM: ${uniqueTags.join(', ')}. ` +
          'This indicates Vuetify components are not being properly registered. ' +
          'Check that vuetify/components and vuetify/directives are imported in the Vuetify plugin.'
      );
    }
  });
});

/**
 * Wait for the page to fully load with Vuetify initialized
 */
Cypress.Commands.add('waitForVuetify', () => {
  cy.get('#app').should('exist');
  cy.verifyVuetifyLoaded();
  cy.verifyNoRawVuetifyTags();
});

export {};
