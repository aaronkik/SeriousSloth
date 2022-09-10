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
      <div className='flex flex-col items-center gap-2 py-8 text-center'>
        <Heading variant='h1'>User Search</Heading>
      </div>
      <UserSearch />
    </>
  );
};

export default UserSearchPage;
