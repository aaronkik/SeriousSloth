import { render, screen, waitFor } from '~/testing-library-setup';
import userEvent from '@testing-library/user-event';
import {
  usernameLengthMessage,
  usernamePatternMessage,
  usernameRequired,
} from '~/constants/form';
import {
  invalidUsernames,
  longUsername,
  shortUsername,
} from '~/__mocks__/data/twitch/constants';
import SearchForm from './search-form';

const mockSetUserResponse = jest.fn();

const setup = () => {
  const user = userEvent.setup();
  render(<SearchForm setUserResponse={mockSetUserResponse} />);
  const input = screen.getByLabelText('user-search');
  const submitButton = screen.getByRole('button', { name: /search/i });
  return { input, submitButton, user };
};

describe('<SearchForm />', () => {
  it('Displays required message when no username is submitted', async () => {
    const { user, input, submitButton } = setup();

    expect(input).toHaveValue('');

    await user.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(usernameRequired)).toBeInTheDocument();
    });
  });

  it('Displays username length message when short username is submitted', async () => {
    const { user, input, submitButton } = setup();

    await user.type(input, shortUsername);
    expect(input).toHaveValue(shortUsername);

    await user.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(usernameLengthMessage)).toBeInTheDocument();
    });
  });

  it('Displays username length message when long username is submitted', async () => {
    const { user, input, submitButton } = setup();

    await user.type(input, longUsername);
    expect(input).toHaveValue(longUsername);

    await user.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(usernameLengthMessage)).toBeInTheDocument();
    });
  });

  it('Displays username pattern message when invalid username is submitted', async () => {
    const { user, input, submitButton } = setup();

    for (const Username of invalidUsernames) {
      await user.type(input, Username);
      expect(input).toHaveValue(Username);

      await user.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText(usernamePatternMessage)).toBeInTheDocument();
      });

      await user.clear(input);
    }
  });
});
