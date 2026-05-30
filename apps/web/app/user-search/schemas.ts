import * as z from 'zod';

const USERNAME_MIN_LENGTH = 3;
const USERNAME_MAX_LENGTH = 25;

export const usernameSchema = z
  .string()
  .min(
    USERNAME_MIN_LENGTH,
    `Username must be more than ${USERNAME_MIN_LENGTH - 1} characters`,
  )
  .max(
    USERNAME_MAX_LENGTH,
    `Username must be less than ${USERNAME_MAX_LENGTH + 1} characters`,
  )
  .regex(/^\w+$/, 'Username can only contain alphanumeric characters');
