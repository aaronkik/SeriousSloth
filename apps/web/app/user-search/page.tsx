import type { Metadata } from 'next';
import { Suspense } from 'react';
import UserSearchContent from '~/app/user-search/user-search-content';
import UserSearchSkeleton from '~/app/user-search/user-search-skeleton';
import { Heading } from '~/components/shared';

export const metadata: Metadata = {
  title: 'User Search | SeriousSloth',
  description: 'Search Twitch users by username.',
};

interface PageProps {
  searchParams: Promise<{ username?: string }>;
}

const Page = ({ searchParams }: PageProps) => (
  <>
    <div className='mb-6 flex flex-col items-center gap-2 text-center'>
      <Heading variant='h1'>User Search</Heading>
      <p>Search Twitch users by username.</p>
    </div>
    <Suspense fallback={<UserSearchSkeleton />}>
      <UserSearchContent searchParams={searchParams} />
    </Suspense>
  </>
);

export default Page;
