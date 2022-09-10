import { GetUsers } from '~/types/twitch';
import UserCard from './user-card';

interface Props {
  userResponse: GetUsers;
}

const User = ({ userResponse }: Props) =>
  userResponse.data.length ? (
    <div className='flex flex-col items-center gap-4' data-testid='userResult'>
      <UserCard user={userResponse.data[0]} />
    </div>
  ) : (
    <div data-testid='userNotFound'>
      <p className='text-center text-xl font-semibold tracking-wide'>
        User not found
      </p>
    </div>
  );

export default User;
