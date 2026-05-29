import { fetchClientCredentials, fetchUsers } from '~/lib/twitch';
import type { GetUsers } from '~/types/twitch';

export type SearchTwitchUserResult = GetUsers | { error: string };

export const searchTwitchUser = async (
  username: string,
): Promise<SearchTwitchUserResult> => {
  'use cache';

  try {
    const { access_token } = await fetchClientCredentials();
    return await fetchUsers(access_token, username);
  } catch (error) {
    console.error(error);
    return { error: 'Unknown error' };
  }
};
