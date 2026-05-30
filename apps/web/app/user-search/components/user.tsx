import { searchTwitchUser } from '~/app/user-search/queries';
import { usernameSchema } from '~/app/user-search/schemas';
import UserNotFound from '~/app/user-search/components/user-not-found';
import UserCard from '~/app/user-search/components/user-card';

interface Props {
  username?: string;
}

const User = async ({ username }: Props) => {
  const parsed = usernameSchema.safeParse(username);

  if (!parsed.success) {
    return null;
  }

  const twitchUserResponse = await searchTwitchUser(parsed.data);

  if (!twitchUserResponse || 'error' in twitchUserResponse) {
    return <UserNotFound />;
  }

  return <UserCard user={twitchUserResponse} />;
};

export default User;
