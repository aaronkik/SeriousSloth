import { GetUsers } from '~/types/twitch';

interface Props {
  userResponse: GetUsers;
}

const User = ({ userResponse }: Props) =>
  userResponse.data.length ? (
    <div data-testid='userResult'>
      <pre>{JSON.stringify(userResponse.data[0], undefined, 2)}</pre>
    </div>
  ) : (
    <div data-testid='userNotFound'>User not found</div>
  );

export default User;
