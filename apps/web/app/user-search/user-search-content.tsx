import UserSearchForm from '~/app/user-search/user-search-form';
import { searchTwitchUser } from '~/app/user-search/queries';
import { usernameSchema } from '~/app/user-search/schemas';
import { Card } from '~/components/shared';
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

  return (
    <div className='flex flex-col gap-4'>
      <Card className='w-full px-2 py-2 md:w-1/2 md:self-center'>
        <UserSearchForm defaultUsername={parsed.data} />
      </Card>
      {user ? <User userResponse={user} /> : null}
    </div>
  );
};

export default UserSearchContent;
