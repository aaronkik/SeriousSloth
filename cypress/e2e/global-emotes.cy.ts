describe('Global emotes page', () => {
  beforeEach(() => {
    cy.visit('/global-emotes');
  });

  it('Displays list of emotes', () => {
    cy.get('ul[data-testid="globalEmoteList"]>li').should('be.visible');
  });

  it('Each emote list item includes an image with a corresponding name', () => {
    cy.get('ul[data-testid="globalEmoteList"]>li').each(($el, index) => {
      cy.get(`[data-testid="emoteImage${index}"]`).should('be.visible');
      cy.get(`[data-testid="emoteName${index}"]`).should('be.visible');
    });
  });
});
