export type OAuthClientCredentialsErrorResponse = {
  status: number;
  message: string;
};

/**
 * @see https://dev.twitch.tv/docs/authentication/getting-tokens-oauth/#client-credentials-grant-flow
 */
export type OAuthClientCredentialsSuccessResponse = {
  access_token: string;
  expires_in: number;
  token_type: string;
};

export type OAuthClientCredentialsResponse =
  | OAuthClientCredentialsErrorResponse
  | OAuthClientCredentialsSuccessResponse;
