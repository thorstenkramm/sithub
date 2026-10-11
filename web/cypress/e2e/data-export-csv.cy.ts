const JSON_API = 'application/vnd.api+json';

describe('CSV booking export', () => {
  const ids: string[] = [];
  let from = '';
  let to = '';
  let anna = '';

  function setRole(flag: boolean) {
    return cy.request({ method: 'PATCH', url: `/api/v1/users/${anna}`,
      headers: { 'Content-Type': JSON_API },
      body: { data: { type: 'users', id: anna, attributes: { is_data_exporter: flag } } } });
  }

  function selectDate(field: string, date: string) {
    cy.get(`[data-cy="${field}"]`).click();
    // The end of the seeded range may fall in the following month.
    cy.get(`[data-cy="${field}-calendar"]`).then(($calendar) => {
      if (!$calendar.find(`[data-v-date="${date}"] button`).length) {
        cy.wrap($calendar).find('button[aria-label="Next month"]').click();
      }
    });
    cy.get(`[data-cy="${field}-calendar"] [data-v-date="${date}"] button`).click();
  }

  before(() => {
    cy.login();
    cy.request('/api/v1/me').then(({ body }) => {
      anna = body.data.id;
      setRole(true);
    });
    cy.request('/api/v1/ping').then((response) => {
      const date = new Date(response.headers.date as string);
      date.setUTCDate(date.getUTCDate() + 1);
      for (let i = 0; i < 26; i++) {
        if (i > 0 && i % 2 === 0) date.setUTCDate(date.getUTCDate() + 1);
        while ([0, 6].includes(date.getUTCDay())) date.setUTCDate(date.getUTCDate() + 1);
        const day = date.toISOString().slice(0, 10);
        if (i === 0) from = day;
        to = day;
        cy.request({ method: 'POST', url: '/api/v1/bookings', headers: { 'Content-Type': JSON_API },
          body: { data: { type: 'bookings', attributes: {
            item_id: i % 2 === 0 ? 'desk_1' : 'cube_1_desk_1', booking_date: day,
            is_guest: true, for_user_name: ['Jürgen Müßig', '=HYPERLINK("x")', '+1', '-1', '@cmd'][i] ?? `CSV Guest ${i}`,
            guest_email: 'csv@example.com'
          } } }
        }).then(({ body }) => ids.push(body.data.id));
      }
    });
  });

  after(() => {
    cy.login();
    ids.forEach((id) => cy.request('DELETE', `/api/v1/bookings/${id}`));
    if (anna) setRole(false);
  });

  it('should download every filtered row with BOM, umlauts, quoting and formula protection', () => {
    cy.login();
    cy.intercept({ method: 'GET', pathname: '/api/v1/bookings/export' }).as('table');
    cy.intercept({ method: 'GET', pathname: '/api/v1/bookings/export.csv' }).as('csv');
    cy.visit('/data-export');
    cy.wait('@table');
    selectDate('export-to', to);
    cy.wait('@table');
    selectDate('export-from', from);
    cy.wait('@table').then(({ response }) => {
      const total = response!.body.meta.total;
      expect(total).to.be.greaterThan(20);
      cy.get('[data-cy="data-export-csv-button"]').click();
      const filename = `sithub-bookings-${from}_${to}.csv`;
      cy.wait('@csv').then(({ response: download }) => {
        expect(download!.statusCode).to.eq(200);
        expect(download!.headers['content-disposition']).to.contain(filename);
      });
      cy.readFile(`${Cypress.config('downloadsFolder')}/${filename}`, null).then((buffer) => {
        expect([...buffer.subarray(0, 3)]).to.deep.eq([239, 187, 191]);
        const text = buffer.toString('utf8');
        expect(text).to.contain('Name,Email,Booked item,Booking date\r\n');
        expect(text).to.contain('Jürgen Müßig');
        expect(text).to.contain('"\'=HYPERLINK(""x"")"');
        ['+1', '-1', '@cmd'].forEach((name) => expect(text).to.contain(`'${name},`));
        const lines = text.trimEnd().split('\r\n');
        expect(lines).to.have.length(total + 1);
        lines.slice(1).forEach((line: string) => expect(line.split(',').at(-1)).to.match(/^\d{4}-\d{2}-\d{2}$/));
      });
    });
    cy.get('[data-cy="export-area"]').click();
    cy.get('.v-overlay-container').contains('.v-list-item', 'Test Area').click();
    cy.wait('@table').then(({ response }) => {
      expect(response!.body.meta.total).to.eq(13);
      cy.get('[data-cy="data-export-csv-button"]').click();
      cy.wait('@csv').then(({ request, response: download }) => {
        expect(new URL(request.url).searchParams.get('filter[area]')).to.eq('test_area');
        expect(download!.body).to.contain('Jürgen Müßig');
        expect(download!.body).not.to.contain('Cube 1 /');
        expect(download!.body).not.to.contain('HYPERLINK');
        expect(download!.body.trimEnd().split('\r\n')).to.have.length(14);
      });
    });
  });

  it('should deny the CSV endpoint to an admin without the exporter role', () => {
    cy.login();
    setRole(false);
    cy.request({ url: '/api/v1/bookings/export.csv', failOnStatusCode: false }).its('status').should('eq', 403);
  });
});
