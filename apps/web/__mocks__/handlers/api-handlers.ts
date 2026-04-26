import { rest } from 'msw';
import {
  usernameMinLength,
  usernameMaxLength,
  usernameRegex,
} from '~/constants/twitch';
import { fetchClientCredentials, fetchUsers } from '~/lib/twitch';
import { UserSearchApiRequest, UserSearchApiResponse } from '~/types/api';

const handlers = [
  rest.post<UserSearchApiRequest, {}, UserSearchApiResponse>(
    '/api/user-search',
    async (req, res, ctx) => {
      const { username } = await req.json<UserSearchApiRequest>();

      if (!username) {
        return res(
          ctx.status(400),
          ctx.json({ status: 400, message: 'No username passed' })
        );
      }

      if (
        username.length < usernameMinLength ||
        username.length > usernameMaxLength
      ) {
        return res(
          ctx.status(400),
          ctx.json({
            status: 400,
            message: `Username must be between ${usernameMinLength} and ${usernameMaxLength}`,
          })
        );
      }

      if (!usernameRegex.test(username)) {
        return res(
          ctx.status(400),
          ctx.json({
            status: 400,
            message: 'Username can only contain alphanumeric characters',
          })
        );
      }

      const { access_token } = await fetchClientCredentials();
      const users = await fetchUsers(access_token, username);

      return res(ctx.status(200), ctx.json(users));
    }
  ),
];

export default handlers;
