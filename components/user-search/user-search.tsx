import { useState } from 'react';
import { GetUsers } from '~/types/twitch';
import SearchForm from './search-form';
import User from './user';

const UserSearch = () => {
  const [userResponse, setUserResponse] = useState<GetUsers | undefined>(
    undefined
  );

  return (
    <>
      <SearchForm setUserResponse={setUserResponse} />
      {userResponse && <User userResponse={userResponse} />}
    </>
  );
};

export default UserSearch;
