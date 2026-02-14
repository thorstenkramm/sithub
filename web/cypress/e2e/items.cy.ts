import {
  createMockItem,
  createMockItemsResponse,
  setupItemsPageIntercepts
} from '../support/flows';

describe('items', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
    cy.login();
  });

  it('should show items with equipment for an item group', () => {
    setupItemsPageIntercepts();

    cy.visit('/');

    // Verify Vuetify is properly loaded (catches component import issues)
    cy.waitForVuetify();

    // Click first area (works with any test data)
    cy.wait('@listAreas').its('response.statusCode').should('eq', 200);
    cy.get('[data-cy="area-item"]').first().click();

    // Click first item group
    cy.wait('@listItemGroups').its('response.statusCode').should('eq', 200);
    cy.get('[data-cy="item-group-item"]').first().click();

    cy.wait('@listItems').then((interception) => {
      expect(interception.response?.statusCode).to.eq(200);
      expect(interception.request.url).to.include('date=');
    });
    cy.location('pathname').should('match', /\/item-groups\/.*\/items/);
    // Check item entries exist (name depends on config)
    cy.get('[data-cy="item-entry"]').should('have.length.at.least', 1);

    // Verify item cards render as proper Vuetify components
    cy.get('[data-cy="item-entry"]').first().should('have.class', 'v-card');

    // Verify status chip renders as Vuetify chip
    cy.get('[data-cy="item-status"]').first().should('have.class', 'v-chip');
  });

  it('should book an available item and show success message', () => {
    // Mock items response with an available item
    const mockItem = createMockItem('item-available-1', 'Available Item');
    cy.intercept('GET', '/api/v1/item-groups/*/items*', createMockItemsResponse([mockItem])).as(
      'listItems'
    );

    // Mock successful booking response
    cy.intercept('POST', '/api/v1/bookings', {
      statusCode: 201,
      headers: { 'Content-Type': 'application/vnd.api+json' },
      body: {
        data: {
          id: 'booking-1',
          type: 'bookings',
          attributes: {
            date: '2026-01-20',
            item_id: 'item-available-1'
          }
        }
      }
    }).as('createBooking');

    cy.visit('/item-groups/test_room/items');

    cy.wait('@listItems');
    cy.location('pathname').should('match', /\/item-groups\/.*\/items/);

    // Click book button on the available item
    cy.get('[data-cy="item-entry"][data-cy-availability="available"]')
      .first()
      .find('[data-cy="book-item-btn"]')
      .click();

    cy.wait('@createBooking').then((interception) => {
      expect(interception.response?.statusCode).to.eq(201);
      expect(interception.request.body.data.type).to.eq('bookings');
    });

    // Success message should appear
    cy.get('[data-cy="booking-success"]').should('contain', 'booked successfully');
  });

  it('should show conflict message with prompt when item is already booked', () => {
    // Mock items response with an available item
    const mockItem = createMockItem('item-mock-1', 'Mock Item 1');
    cy.intercept('GET', '/api/v1/item-groups/*/items*', createMockItemsResponse([mockItem])).as(
      'listItems'
    );

    // Mock a 409 Conflict response for booking
    cy.intercept('POST', '/api/v1/bookings', {
      statusCode: 409,
      headers: { 'Content-Type': 'application/vnd.api+json' },
      body: {
        errors: [
          {
            status: '409',
            title: 'Conflict',
            detail: 'Item is already booked for this date',
            code: 'conflict'
          }
        ]
      }
    }).as('createBookingConflict');

    cy.visit('/item-groups/test_room/items');

    cy.wait('@listItems');

    // Click book on an available item
    cy.get('[data-cy="item-entry"][data-cy-availability="available"]')
      .first()
      .find('[data-cy="book-item-btn"]')
      .click();

    cy.wait('@createBookingConflict');

    // Error message should show backend detail + prompt
    cy.get('[data-cy="booking-error"]')
      .should('contain', 'Item is already booked for this date')
      .and('contain', 'Please choose another item');

    // Item list should be refreshed
    cy.wait('@listItems');
  });

  it('should show booker name on occupied items', () => {
    const occupiedItem = {
      id: 'item-occupied-1',
      type: 'items',
      attributes: {
        name: 'Occupied Desk',
        equipment: ['Monitor'],
        availability: 'occupied',
        booker_name: 'Alice Smith'
      }
    };
    const availableItem = createMockItem('item-free-1', 'Free Desk');

    cy.intercept('GET', '/api/v1/item-groups/*/items*', {
      statusCode: 200,
      body: { data: [occupiedItem, availableItem] }
    }).as('listItems');

    cy.visit('/item-groups/test_room/items');
    cy.wait('@listItems');

    // Occupied item should show booker name
    cy.get('[data-cy="item-entry"][data-cy-availability="occupied"]')
      .find('[data-cy="item-booker"]')
      .should('contain', 'Alice Smith');

    // Available item should not show booker name
    cy.get('[data-cy="item-entry"][data-cy-availability="available"]')
      .find('[data-cy="item-booker"]')
      .should('not.exist');
  });

  it('should show self-duplicate message when user already has booking', () => {
    // Mock items response with an available item
    const mockItem = createMockItem('item-mock-2', 'Mock Item 2');
    cy.intercept('GET', '/api/v1/item-groups/*/items*', createMockItemsResponse([mockItem])).as(
      'listItems'
    );

    // Mock a 409 Conflict response for self-duplicate
    cy.intercept('POST', '/api/v1/bookings', {
      statusCode: 409,
      headers: { 'Content-Type': 'application/vnd.api+json' },
      body: {
        errors: [
          {
            status: '409',
            title: 'Conflict',
            detail: 'You already have this item booked for this date',
            code: 'conflict'
          }
        ]
      }
    }).as('createBookingSelfDuplicate');

    cy.visit('/item-groups/test_room/items');

    cy.wait('@listItems');

    cy.get('[data-cy="item-entry"][data-cy-availability="available"]')
      .first()
      .find('[data-cy="book-item-btn"]')
      .click();

    cy.wait('@createBookingSelfDuplicate');

    // Error message should show the self-duplicate message
    cy.get('[data-cy="booking-error"]')
      .should('contain', 'You already have this item booked for this date')
      .and('contain', 'Please choose another item');
  });
});
