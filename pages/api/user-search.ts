import {
  usernameMaxLength,
  usernameMinLength,
  usernameRegex,
} from '~/constants/twitch';
import { fetchClientCredentials, fetchUsers } from '~/lib/twitch';
import {
  NextUserSearchApiRequest,
  NextUserSearchApiResponse,
} from '~/types/api';

const userSearch = async (
  req: NextUserSearchApiRequest,
  res: NextUserSearchApiResponse
) => {
  const {
    body: { username },
    method,
  } = req;

  if (method !== 'POST') {
    return res.status(405).json({ status: 405, message: 'Method not allowed' });
  }

  if (!username) {
    return res.status(400).json({ status: 400, message: 'No username passed' });
  }

  if (
    username.length < usernameMinLength ||
    username.length > usernameMaxLength
  ) {
    return res.status(400).json({
      status: 400,
      message: `Username must be between ${usernameMinLength} and ${usernameMaxLength}`,
    });
  }

  if (!usernameRegex.test(username)) {
    return res.status(400).json({
      status: 400,
      message: 'Username can only contain alphanumeric characters',
    });
  }

  try {
    const { access_token } = await fetchClientCredentials();
    const users = await fetchUsers(access_token, username);
    return res.status(200).json(users);
  } catch (error: any) {
    console.error(error);
    return res.status(400).json({
      status: error?.status || 400,
      message: error?.message || 'Unknown error',
    });
  }
};

export default userSearch;
