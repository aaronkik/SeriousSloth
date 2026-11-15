import { http, HttpResponse, type PathParams } from 'msw';
import {
  usernameMinLength,
  usernameMaxLength,
  usernameRegex,
} from '~/constants/twitch';
import { fetchClientCredentials, fetchUsers } from '~/lib/twitch';
import { UserSearchApiRequest, UserSearchApiResponse } from '~/types/api';

const handlers = [
  http.post<PathParams, UserSearchApiRequest, UserSearchApiResponse>(
    '/api/user-search',
    async ({ request }) => {
      const { username } = await request.json();

      if (!username) {
        return HttpResponse.json(
          { status: 400, message: 'No username passed' },
          { status: 400 },
        );
      }

      if (
        username.length < usernameMinLength ||
        username.length > usernameMaxLength
      ) {
        return HttpResponse.json(
          {
            status: 400,
            message: `Username must be between ${usernameMinLength} and ${usernameMaxLength}`,
          },
          { status: 400 },
        );
      }

      if (!usernameRegex.test(username)) {
        return HttpResponse.json(
          {
            status: 400,
            message: 'Username can only contain alphanumeric characters',
          },
          { status: 400 },
        );
      }

      const { access_token } = await fetchClientCredentials();
      const users = await fetchUsers(access_token, username);

      return HttpResponse.json(users, { status: 200 });
    },
  ),
];

export default handlers;
