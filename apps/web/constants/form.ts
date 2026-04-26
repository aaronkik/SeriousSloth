import { RegisterOptions } from 'react-hook-form';
import { usernameMaxLength, usernameMinLength, usernameRegex } from './twitch';

export const usernameRequired = 'Username is required';
export const usernameLengthMessage = `Username must be between ${usernameMinLength} and ${usernameMaxLength}`;
export const usernamePatternMessage = `Username can only contain alphanumeric characters`;

export const usernameFormRules: RegisterOptions = {
  required: usernameRequired,
  minLength: {
    value: usernameMinLength,
    message: usernameLengthMessage,
  },
  maxLength: {
    value: usernameMaxLength,
    message: usernameLengthMessage,
  },
  pattern: {
    value: usernameRegex,
    message: usernamePatternMessage,
  },
};
