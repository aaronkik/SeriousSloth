import { searchTwitchUser } from '~/app/user-search/queries';
import { usernameSchema } from '~/app/user-search/schemas';
import User from '~/components/user-search/user';

interface Props {
  searchParams: Promise<{ username?: string }>;
}

const UserSearchContent = async ({ searchParams }: Props) => {
  const { username } = await searchParams;
  const parsed = usernameSchema.safeParse(username);

  const result = parsed?.success
    ? await searchTwitchUser(parsed.data)
    : undefined;

  const user = result && !('error' in result) ? result : undefined;

  return user ? <User userResponse={user} /> : null;
};

export default UserSearchContent;
