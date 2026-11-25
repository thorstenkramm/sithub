import type { Interception } from 'cypress/types/net-stubbing';

const weights = [400, 500, 600, 700];
const sample = 'SitHub Inter — Привіт, Україна!';

describe('local fonts', { browser: { family: 'chromium' } }, () => {
  it('loads text and icon fonts with only the SitHub origin reachable', () => {
    const origin = new URL(Cypress.config('baseUrl')!).origin;
    const blocked: string[] = [];
    cy.login();
    cy.then(() => Cypress.automation('remote:debugger:protocol', { command: 'Network.clearBrowserCache' }));
    cy.intercept({ url: '**', middleware: true }, (req) => {
      if (new URL(req.url).origin !== origin) {
        blocked.push(req.url);
        req.destroy();
      }
    });
    cy.intercept({ method: 'GET', pathname: '**/*.woff2' }).as('fonts');
    cy.viewport(1280, 800);
    cy.visit('/');
    cy.get('[data-cy="nav-areas"]').should('be.visible').and(($nav) => {
      expect($nav.css('font-family').split(',')[0].replace(/["']/g, '')).to.equal('Inter');
    });

    cy.document().then(async (doc) => {
      for (const weight of weights) {
        const faces = await doc.fonts.load(`${weight} 16px Inter`, sample);
        expect(faces.some((face) => face.weight === String(weight) && face.status === 'loaded')).to.be.true;
      }

      // Configured icons use MDI classes, while many default icons render as SVG.
      const icon = doc.createElement('i');
      icon.className = 'mdi mdi-home';
      doc.body.append(icon);
      try {
        const glyph = doc.defaultView!.getComputedStyle(icon, '::before');
        expect(glyph.fontFamily).to.contain('Material Design Icons');
        const text = glyph.content.replace(/["']/g, '');
        expect(text).to.have.length.greaterThan(0).and.not.equal('none');
        const faces = await doc.fonts.load('16px "Material Design Icons"', text);
        expect(faces.some((face) => face.status === 'loaded')).to.be.true;
      } finally {
        icon.remove();
      }
    });

    cy.get<Interception[]>('@fonts.all').should((requests) => {
      const paths = requests.map(({ request, response }) => {
        const url = new URL(request.url);
        expect(url.origin).to.equal(origin);
        expect(response?.statusCode, url.pathname).to.equal(200);
        return url.pathname;
      });
      expect(paths.some((p) => p.includes('materialdesignicons-webfont'))).to.be.true;
      for (const weight of weights) {
        for (const subset of ['latin', 'cyrillic']) {
          const name = `inter-${subset}-${weight}-normal`;
          expect(paths.some((p) => p.includes(name)), name).to.be.true;
        }
      }
      expect(blocked, 'attempted third-party requests').to.be.empty;
    });
  });
});
