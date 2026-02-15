/** Opens the password change dialog from the user menu. */
function openPasswordDialog() {
  cy.get('[data-cy="user-menu-trigger"]').click();
  cy.get('[data-cy="change-password-btn"]').click();
  cy.get('[data-cy="current-password"]').should('exist');
}

/** Fills the password change form fields. */
function fillPasswordForm(currentPassword: string, newPassword: string) {
  cy.get('[data-cy="current-password"]').find('input').type(currentPassword);
  cy.get('[data-cy="new-password"]').find('input').type(newPassword);
}

describe('password change', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
    cy.login();
    cy.visit('/');
  });

  it('should show change password option in user menu for local users', () => {
    cy.get('[data-cy="user-menu-trigger"]').click();
    cy.get('[data-cy="change-password-btn"]').should('exist');
  });

  it('should show change password icon in desktop and mobile menus', () => {
    cy.get('[data-cy="user-menu-trigger"]').click();
    cy.get('[data-cy="change-password-icon"]').should('exist');
    cy.get('body').click(0, 0);

    cy.viewport(375, 667);
    cy.get('[data-cy="mobile-menu-btn"]').click();
    cy.get('[data-cy="mobile-change-password-icon"]').should('exist');
    cy.viewport(1280, 720);
  });

  it('should open and close the password change dialog', () => {
    openPasswordDialog();

    cy.get('[data-cy="new-password"]').should('exist');
    cy.get('[data-cy="password-submit"]').should('exist');

    cy.get('[data-cy="password-cancel"]').click();
    cy.get('[data-cy="current-password"]').should('not.exist');
  });

  it('should show error for wrong current password', () => {
    cy.intercept('PATCH', '/api/v1/me').as('changePassword');

    openPasswordDialog();
    fillPasswordForm('WrongPassword123!', 'NewValidPass2026!!');
    cy.get('[data-cy="password-submit"]').click();

    cy.wait('@changePassword').its('response.statusCode').should('eq', 401);
    cy.get('[data-cy="password-error"]').should('exist');
  });

  it('should show error for short new password', () => {
    cy.intercept('PATCH', '/api/v1/me').as('changePassword');

    openPasswordDialog();
    fillPasswordForm(Cypress.env('testUserPassword'), 'short');
    cy.get('[data-cy="password-submit"]').click();

    cy.wait('@changePassword').its('response.statusCode').should('eq', 400);
    cy.get('[data-cy="password-error"]').should('exist');
  });

  it('should change password successfully', () => {
    const originalPassword = Cypress.env('testUserPassword');
    const newPassword = 'NewTestPass2026!!';

    cy.intercept('PATCH', '/api/v1/me').as('changePassword');

    openPasswordDialog();
    fillPasswordForm(originalPassword, newPassword);
    cy.get('[data-cy="password-submit"]').click();

    cy.wait('@changePassword').its('response.statusCode').should('eq', 200);
    cy.get('[data-cy="password-success"]').should('exist');

    // Close dialog, log out, log in with new password
    cy.get('[data-cy="password-cancel"]').click();
    cy.get('[data-cy="user-menu-trigger"]').click();
    cy.get('[data-cy="logout-btn"]').click();
    cy.location('pathname').should('eq', '/login');

    cy.login(Cypress.env('testUserEmail'), newPassword);
    cy.visit('/');
    cy.get('[data-cy="user-menu-trigger"]').should('exist');

    // Revert password to original so other tests aren't affected
    cy.request({
      method: 'PATCH',
      url: '/api/v1/me',
      body: {
        data: {
          attributes: {
            current_password: newPassword,
            new_password: originalPassword
          }
        }
      }
    }).then((response) => {
      expect(response.status).to.eq(200);
    });
  });
});
