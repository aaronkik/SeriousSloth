import { Suspense } from 'react';
import User from '~/app/user-search/components/user';
import UserSearchSkeleton from '~/app/user-search/user-search-skeleton';

interface Props {
  searchParams: Promise<{ username?: string }>;
}

const UserSearchResult = async ({ searchParams }: Props) => {
  const { username } = await searchParams;

  if (!username) {
    return null;
  }

  return (
    <Suspense fallback={<UserSearchSkeleton />}>
      <User username={username} />
    </Suspense>
  );
};

export default UserSearchResult;
