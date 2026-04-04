// Cypress support file.
import './commands';

// Force English locale for all E2E tests to prevent auto-detection
// from the OS/browser locale affecting assertions on visible text.
Cypress.on('window:before:load', (win) => {
  win.localStorage.setItem('sithub_locale', 'en');
});
