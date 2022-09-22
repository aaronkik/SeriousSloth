describe('Top streams page', () => {
  beforeEach(() => {
    cy.visit('/historic-top-streams');

    cy.get('ul[data-testid="historicTopStreamTimes"]').within(() => {
      cy.get('a').first().click();
    });

    cy.location('pathname').should('match', /^\/historic-top-streams\/.+$/);
  });

  it('Displays a heading', () => {
    cy.get('h1').should('be.visible');
  });

  it('Displays a list of top streams', () => {
    cy.get('ul[data-testid="topStreams"]>li').should('be.visible');
  });

  it('Clicking link goes back to /historic-top-streams', () => {
    cy.get('a[data-testid="historicTopStreamsLink"]').click();
    cy.location('pathname').should('eq', '/historic-top-streams');
  });
});
