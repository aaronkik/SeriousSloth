export const clientId = process.env.NEXT_PUBLIC_TWITCH_CLIENT_ID!;
export const clientSecret = process.env.TWITCH_CLIENT_SECRET!;

// Endpoints
export const globalEmotesEndpoint =
  'https://api.twitch.tv/helix/chat/emotes/global';
export const oauth2TokenEndpoint = 'https://id.twitch.tv/oauth2/token';
