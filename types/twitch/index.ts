export type GlobalEmotesErrorResponse = {
  error: string;
  message: string;
  status: number;
};

/**
 * @see https://dev.twitch.tv/docs/api/reference#get-global-emotes
 */
export type GlobalEmotesSuccessResponse = {
  data: Array<{
    id: string;
    name: string;
    images: {
      url_1x: string;
      url_2x: string;
      url_4x: string;
    };
    format: Array<'static' | 'animated'>;
    scale: Array<'1.0' | '2.0' | '3.0'>;
    theme_mode: Array<'dark' | 'light'>;
  }>;
  template: string;
};

export type GlobalEmotesResponse =
  | GlobalEmotesErrorResponse
  | GlobalEmotesSuccessResponse;

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
