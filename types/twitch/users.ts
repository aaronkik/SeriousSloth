export type GetUsersErrorResponse = {
  error: string;
  status: number;
  message: string;
};

type User = {
  broadcaster_type: '' | 'partner' | 'affiliate';
  description: string;
  display_name: string;
  id: string;
  login: string;
  offline_image_url: string;
  profile_image_url: string;
  type: 'staff' | 'admin' | 'global_mod' | '';
  /** @deprecated */
  view_count: number;
  email?: string;
  created_at: string;
};

/**
 * @see https://dev.twitch.tv/docs/api/reference#get-users
 */
export type GetUsers = {
  data: Array<User>;
};

export type GetUsersResponse = GetUsers | GetUsersErrorResponse;
