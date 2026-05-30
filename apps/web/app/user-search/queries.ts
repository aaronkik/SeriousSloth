import { fetchClientCredentials, fetchUsers } from '~/lib/twitch';
import { User } from '~/types/twitch';

export type SearchTwitchUserResult = User | null | { error: string };

export const searchTwitchUser = async (
  username: string,
): Promise<SearchTwitchUserResult> => {
  'use cache';

  try {
    const { access_token } = await fetchClientCredentials();
    const { data } = await fetchUsers(access_token, username);
    return data[0] ?? null;
  } catch (error) {
    console.error(error);
    return { error: 'Internal error' };
  }
};
