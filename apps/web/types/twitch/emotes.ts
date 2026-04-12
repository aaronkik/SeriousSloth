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
