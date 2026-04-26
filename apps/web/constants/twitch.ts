export const clientId = process.env.NEXT_PUBLIC_TWITCH_CLIENT_ID!;
export const clientSecret = process.env.TWITCH_CLIENT_SECRET!;

export const usernameMinLength = 3;
export const usernameMaxLength = 25;
export const usernameRegex = new RegExp(
  `^\\w{${usernameMinLength},${usernameMaxLength}}$`
);

// Endpoints
export const globalEmotesEndpoint =
  'https://api.twitch.tv/helix/chat/emotes/global';
export const oauth2TokenEndpoint = 'https://id.twitch.tv/oauth2/token';
export const usersEndpoint = 'https://api.twitch.tv/helix/users';
