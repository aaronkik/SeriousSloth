import type { Metadata } from 'next';
import { Heading } from '~/components/shared';
import { UserSearch } from '~/components/user-search';

export const metadata: Metadata = {
  title: 'User Search | SeriousSloth',
  description: 'Search Twitch users by username.',
};

const Page = () => {
  return (
    <>
      <div className='mb-6 flex flex-col items-center gap-2 text-center'>
        <Heading variant='h1'>User Search</Heading>
        <p>Search Twitch users by username.</p>
      </div>
      <UserSearch />
    </>
  );
};

export default Page;
