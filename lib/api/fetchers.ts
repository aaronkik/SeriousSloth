import { UserSearchApiRequest, UserSearchApiResponse } from '~/types/api';

export const fetchUsers = async ({ username }: UserSearchApiRequest) => {
  const usersResponse = await fetch('/api/user-search', {
    method: 'POST',
    headers: new Headers({ 'Content-Type': 'application/json' }),
    body: JSON.stringify({ username }),
  });

  const payload = (await usersResponse.json()) as UserSearchApiResponse;

  if (!usersResponse.ok || !('data' in payload)) {
    return Promise.reject(payload);
  }

  return payload;
};
