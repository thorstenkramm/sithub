const interWeights = [400, 500, 600, 700] as const;
const fontProbeText = 'SitHub Inter — Привіт, Україна!';

describe('local fonts', () => {
  beforeEach(() => {
    cy.clearCookies();
    cy.clearLocalStorage();
    cy.intercept(/^https?:\/\/fonts\.googleapis\.com\//).as('googleFontStylesheets');
    cy.intercept(/^https?:\/\/fonts\.gstatic\.com\//).as('googleFontFiles');
    cy.intercept('GET', '**/*.woff2').as('fontFiles');
    cy.login();
  });

  it('should load every Inter weight from SitHub without Google requests', () => {
    cy.intercept('GET', '/api/v1/areas').as('listAreas');
    cy.visit('/');
    cy.wait('@listAreas').its('response.statusCode').should('eq', 200);
    cy.get('.v-application').should('exist');

    cy.document().then(async (doc) => {
      const probes = interWeights.map((weight) => {
        const probe = doc.createElement('span');
        probe.textContent = fontProbeText;
        probe.style.fontFamily = 'Inter';
        probe.style.fontWeight = String(weight);
        doc.body.append(probe);
        return probe;
      });

      try {
        const facesByWeight = await Promise.all(
          interWeights.map((weight) => doc.fonts.load(`${weight} 16px "Inter"`, fontProbeText))
        );

        facesByWeight.forEach((faces, index) => {
          const weight = interWeights[index];
          const hasLoadedNormalFace = faces.some(
            (face) =>
              face.family.replaceAll('"', '') === 'Inter' &&
              face.style === 'normal' &&
              face.weight === String(weight) &&
              face.status === 'loaded'
          );
          expect(hasLoadedNormalFace, `normal Inter ${weight} is loaded for English and Ukrainian text`).to.be
            .true;
        });
      } finally {
        probes.forEach((probe) => probe.remove());
      }
    });

    cy.location('origin').then((origin) => {
      cy.get('@fontFiles.all').then((interceptions) => {
        const fontResponses = interceptions as Cypress.Interception[];
        expect(fontResponses.length, 'captured WOFF2 responses').to.be.greaterThan(0);

        fontResponses.forEach((interception) => {
          const url = new URL(interception.request.url);
          expect(url.origin, `response origin for ${url.pathname}`).to.eq(origin);
          expect(interception.response?.statusCode, `response status for ${url.pathname}`).to.eq(200);
        });

        interWeights.forEach((weight) => {
          (['latin', 'cyrillic'] as const).forEach((subset) => {
            const subsetResponses = fontResponses.filter((interception) =>
              new URL(interception.request.url).pathname.includes(
                `inter-${subset}-${weight}-normal`
              )
            );
            expect(subsetResponses, `${subset} Inter ${weight} responses`).to.have.length(1);
          });
        });
      });
    });

    cy.get('@googleFontStylesheets.all').should('have.length', 0);
    cy.get('@googleFontFiles.all').should('have.length', 0);
  });
});
