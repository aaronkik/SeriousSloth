import type { NextApiRequest, NextApiResponse } from 'next';
import { GetUsers } from '~/types/twitch';

type ApiErrorResponse = { status: number; message: string };

export interface UserSearchApiRequest {
  username: string;
}

export interface NextUserSearchApiRequest extends NextApiRequest {
  body: UserSearchApiRequest;
}

export type UserSearchApiResponse = GetUsers | ApiErrorResponse;
export type NextUserSearchApiResponse = NextApiResponse<UserSearchApiResponse>;
