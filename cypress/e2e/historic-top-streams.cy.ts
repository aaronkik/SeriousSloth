describe('Top streams page', () => {
  beforeEach(() => {
    cy.visit('/historic-top-streams');
  });

  it('Displays a heading', () => {
    cy.get('h1').should('be.visible');
  });

  it('Displays a list of historic timestamps for top streams', () => {
    cy.get('ul[data-testid="historicTopStreamTimes"]>li').should('be.visible');
  });
});
