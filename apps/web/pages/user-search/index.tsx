import Head from 'next/head';
import { Heading } from '~/components/shared';
import { UserSearch } from '~/components/user-search';
import { userSearchTitle } from '~/constants/titles';

const UserSearchPage = () => {
  return (
    <>
      <Head>
        <title>{userSearchTitle}</title>
      </Head>
      <div className='mb-6 flex flex-col items-center gap-2 text-center'>
        <Heading variant='h1'>User Search</Heading>
        <p>Search Twitch users by username.</p>
      </div>
      <UserSearch />
    </>
  );
};

export default UserSearchPage;
