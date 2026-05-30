import type { Metadata } from 'next';
import { Suspense } from 'react';
import UserSearchForm from '~/app/user-search/user-search-form';
import UserSearchResult from '~/app/user-search/components/user-search-result';
import { Card, Heading } from '~/components/shared';

export const metadata: Metadata = {
  title: 'User Search | SeriousSloth',
  description: 'Search Twitch users by username.',
};

interface PageProps {
  searchParams: Promise<{ username?: string }>;
}

const Page = ({ searchParams }: PageProps) => {
  return (
    <>
      <div className='mb-6 flex flex-col items-center gap-2 text-center'>
        <Heading variant='h1'>User Search</Heading>
        <p>Search Twitch users by username.</p>
      </div>
      <div className='flex flex-col gap-4'>
        <Card className='w-full px-2 py-2 md:w-1/2 md:self-center'>
          <UserSearchForm />
        </Card>
        <Suspense fallback={null}>
          <UserSearchResult searchParams={searchParams} />
        </Suspense>
      </div>
    </>
  );
};

export default Page;
