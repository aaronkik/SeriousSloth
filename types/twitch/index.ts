/**
 * @see https://dev.twitch.tv/docs/api/reference#get-global-emotes
 */
export type GlobalEmotesResponse = {
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

/**
 * @see https://dev.twitch.tv/docs/authentication/getting-tokens-oauth/#client-credentials-grant-flow
 */
export type OAuthClientCredentialsResponse = {
  access_token: string;
  expires_in: number;
  token_type: string;
};
