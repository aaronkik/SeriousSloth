import {
  invalidLoginNames,
  longLoginName,
  shortLoginName,
  validLoginName,
} from '__mocks__/data/twitch/constants';
import { userNotFound, singleUser } from '__mocks__/data/twitch/users';
import {
  usernameLengthMessage,
  usernameRequired,
  usernamePatternMessage,
} from '~/constants/form';

describe('Global emotes page', () => {
  beforeEach(() => {
    cy.visit('/user-search');
  });

  it('Displays error on submit when username has not been typed', () => {
    cy.get('form').within(() => {
      cy.get('input[type="search"]').should('have.value', '');
      cy.get('button[type="submit"]').click();
      cy.get('[data-testid="searchErrorMessage"]').contains(usernameRequired);
    });
  });

  it('Displays error on submit when username is too short', () => {
    cy.get('form').within(() => {
      cy.get('input[type="search"]').type(shortLoginName);
      cy.get('button[type="submit"]').click();
      cy.get('[data-testid="searchErrorMessage"]').contains(
        usernameLengthMessage
      );
    });
  });

  it('Displays error on submit when username is too long', () => {
    cy.get('form').within(() => {
      cy.get('input[type="search"]').type(longLoginName);
      cy.get('button[type="submit"]').click();
      cy.get('[data-testid="searchErrorMessage"]').contains(
        usernameLengthMessage
      );
    });
  });

  it('Displays error on submit when username pattern does not match', () => {
    cy.get('form').within(() => {
      for (const loginName of invalidLoginNames) {
        cy.get('input[type="search"]').type(loginName);
        cy.get('button[type="submit"]').click();
        cy.get('[data-testid="searchErrorMessage"]').contains(
          usernamePatternMessage
        );
        cy.get('input[type="search"]').clear();
      }
    });
  });

  it('Displays user not found message when no user is returned from a successful API response', () => {
    cy.intercept('POST', '/api/user-search', userNotFound);
    cy.get('form').within(() => {
      cy.get('input[type="search"]').type(validLoginName);
      cy.get('button[type="submit"]').click();
    });
    cy.get('[data-testid="userNotFound"]').should('be.visible');
  });

  it('Displays user from the API', () => {
    cy.intercept('POST', '/api/user-search', singleUser);
    cy.get('form').within(() => {
      cy.get('input[type="search"]').type(validLoginName);
      cy.get('button[type="submit"]').click();
    });
    cy.get('[data-testid="userResult"]').should('be.visible');
  });

  it('Displays error response when an error is returned from the API', () => {
    cy.intercept('POST', '/api/user-search', {
      statusCode: 400,
      body: { status: 400, message: 'Error' },
    });
    cy.get('form').within(() => {
      cy.get('input[type="search"]').type(validLoginName);
      cy.get('button[type="submit"]').click();
    });
    cy.get('[data-testid="userError"]').should('be.visible');
  });
});
