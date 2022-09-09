import { useState } from 'react';
import { Card } from '~/components/shared';
import { GetUsers } from '~/types/twitch';
import SearchForm from './search-form';
import User from './user';

const UserSearch = () => {
  const [userResponse, setUserResponse] = useState<GetUsers | undefined>(
    undefined
  );

  return (
    <>
      <div className='flex flex-col items-center'>
        <Card className='w-full px-2 py-2 md:w-1/2'>
          <SearchForm setUserResponse={setUserResponse} />
        </Card>
        {userResponse && <User userResponse={userResponse} />}
      </div>
    </>
  );
};

export default UserSearch;
