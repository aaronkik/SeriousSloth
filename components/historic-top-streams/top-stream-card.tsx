import { BiJoystick } from 'react-icons/bi';
import { MdPersonOutline } from 'react-icons/md';
import { Card, TwitchLink } from '~/components/shared';
import { definitions } from '~/types/supabase';

type Props = {
  stream: definitions['stream_history'];
};

const TopStreamCard = ({ stream }: Props) => {
  const { game_name, user_login, user_name, viewer_count } = stream;

  return (
    <Card className='flex flex-col p-4'>
      <p className='break-words text-xl font-medium tracking-wide'>
        {user_name}
      </p>
      <div className='flex flex-col gap-1 sm:flex-row sm:gap-4'>
        <div className='flex items-center gap-0.5'>
          <div>
            <MdPersonOutline className='h-5 w-5 text-red-400' />
          </div>
          {viewer_count.toLocaleString()}
        </div>
        <div className='flex items-center gap-0.5'>
          <div>
            <BiJoystick className='h-5 w-5 text-blue-400' />
          </div>
          <p className='break-words'>{game_name}</p>
        </div>
      </div>
      <TwitchLink className='mt-1 w-max' loginName={user_login} />
    </Card>
  );
};

export default TopStreamCard;
