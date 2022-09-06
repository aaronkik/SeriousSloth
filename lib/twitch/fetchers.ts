import {
  clientId,
  clientSecret,
  globalEmotesEndpoint,
  oauth2TokenEndpoint,
  usersEndpoint,
} from '~/constants/twitch';
import {
  GetUsers,
  GetUsersResponse,
  GlobalEmotesResponse,
  GlobalEmotesSuccessResponse,
  OAuthClientCredentialsResponse,
  OAuthClientCredentialsSuccessResponse,
} from '~/types/twitch';

export const fetchClientCredentials =
  async (): Promise<OAuthClientCredentialsSuccessResponse> => {
    const twitchAccessTokenResponse = await fetch(oauth2TokenEndpoint, {
      method: 'POST',
      headers: new Headers({
        'Content-Type': 'application/x-www-form-urlencoded',
      }),
      body: new URLSearchParams({
        client_id: clientId,
        client_secret: clientSecret,
        grant_type: 'client_credentials',
      }).toString(),
    });

    const payload =
      (await twitchAccessTokenResponse.json()) as OAuthClientCredentialsResponse;

    if (!twitchAccessTokenResponse.ok || !('access_token' in payload)) {
      return Promise.reject(payload);
    }

    return payload;
  };

export const fetchGlobalEmotes = async (
  accessToken: string
): Promise<GlobalEmotesSuccessResponse> => {
  const globalEmotesResponse = await fetch(globalEmotesEndpoint, {
    method: 'GET',
    headers: new Headers({
      Authorization: `Bearer ${accessToken}`,
      'Client-Id': clientId,
    }),
  });

  const payload = (await globalEmotesResponse.json()) as GlobalEmotesResponse;

  if (!globalEmotesResponse.ok || !('data' in payload)) {
    return Promise.reject(payload);
  }

  return payload;
};

export const fetchUsers = async (
  accessToken: string,
  username: string
): Promise<GetUsers> => {
  const usernameSearchParams = new URLSearchParams({
    login: username,
  }).toString();

  const usersResponse = await fetch(
    `${usersEndpoint}?${usernameSearchParams}`,
    {
      method: 'GET',
      headers: new Headers({
        Authorization: `Bearer ${accessToken}`,
        'Client-Id': clientId,
      }),
    }
  );

  const payload = (await usersResponse.json()) as GetUsersResponse;

  if (!usersResponse.ok || !('data' in payload)) {
    return Promise.reject(payload);
  }

  return payload;
};
