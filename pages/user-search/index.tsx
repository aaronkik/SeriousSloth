import Head from 'next/head';
import { UserSearch } from '~/components/user-search';
import { userSearchTitle } from '~/constants/titles';

const UserSearchPage = () => {
  return (
    <>
      <Head>
        <title>{userSearchTitle}</title>
      </Head>
      <UserSearch />
    </>
  );
};

export default UserSearchPage;
