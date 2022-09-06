import { usernameMaxLength } from '~/constants/twitch';

export const invalidUsernames = [
  '@bcdef',
  'bad username',
  '&bcdef',
  'abc-def',
  'u$ername',
];
export const longUsername = 'a'.repeat(usernameMaxLength + 1);
export const shortUsername = 'a';
export const validUsername = 'valid_login';
