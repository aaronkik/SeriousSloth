import { render, screen, waitFor } from '~/testing-library-setup';
import userEvent from '@testing-library/user-event';
import { rest } from 'msw';
import { validUsername } from '~/__mocks__/data/twitch/constants';
import { userNotFound } from '~/__mocks__/data/twitch/users';
import { server } from '~/__mocks__/msw-server';
import UserSearch from './user-search';

const setup = () => {
  const user = userEvent.setup();
  render(<UserSearch />);
  const input = screen.getByLabelText('user-search');
  const submitButton = screen.getByRole('button', { name: /search/i });
  return { input, submitButton, user };
};

describe('<UserSearch />', () => {
  it('Displays user not found when no user is returned from API', async () => {
    const { input, submitButton, user } = setup();

    await user.type(input, validUsername);
    expect(input).toHaveValue(validUsername);

    server.use(
      rest.post('/api/user-search', (req, res, ctx) =>
        res(ctx.status(200), ctx.json(userNotFound))
      )
    );

    await user.click(submitButton);

    await waitFor(() => {
      expect(screen.getByTestId('userNotFound')).toBeInTheDocument();
    });
  });

  it('Displays user when user is returned from API', async () => {
    const { input, submitButton, user } = setup();

    await user.type(input, validUsername);
    expect(input).toHaveValue(validUsername);

    await user.click(submitButton);

    await waitFor(() => {
      expect(screen.getByTestId('userResult')).toBeInTheDocument();
    });
  });

  it('Displays error message when user is not returned from API', async () => {
    const { input, submitButton, user } = setup();

    await user.type(input, validUsername);
    expect(input).toHaveValue(validUsername);

    server.use(
      rest.post('/api/user-search', (req, res, ctx) =>
        res(ctx.status(400), ctx.json({ status: 400, message: 'Error' }))
      )
    );

    await user.click(submitButton);

    await waitFor(() => {
      expect(screen.getByTestId('userError')).toBeInTheDocument();
    });
  });
});
