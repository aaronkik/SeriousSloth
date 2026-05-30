import type { Metadata } from 'next';
import { Suspense } from 'react';
import UserSearchContent from '~/app/user-search/user-search-content';
import UserSearchForm from '~/app/user-search/user-search-form';
import UserSearchSkeleton from '~/app/user-search/user-search-skeleton';
import { Card, Heading } from '~/components/shared';

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
    <div className='flex flex-col gap-4'>
      <Card className='w-full px-2 py-2 md:w-1/2 md:self-center'>
        <UserSearchForm />
      </Card>
      <Suspense fallback={<UserSearchSkeleton />}>
        <UserSearchContent searchParams={searchParams} />
      </Suspense>
    </div>
  </>
);

export default Page;
