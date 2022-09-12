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
      <Heading className='mb-2 text-center' variant='h1'>
        User Search
      </Heading>
      <UserSearch />
    </>
  );
};

export default UserSearchPage;
