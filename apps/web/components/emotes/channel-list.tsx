import Link from 'next/link';
import { Card, Heading } from '~/components/shared';
import type { Channel } from '~/lib/api/emotes-service';

type Props = {
  channels: Channel[];
};

const ChannelList = ({ channels }: Props) => (
  <ul
    data-testid='channelList'
    className='grid grid-cols-1 grid-rows-1 gap-6 sm:grid-cols-2'
  >
    {channels.map(({ id, displayName }) => (
      <li key={id}>
        <Link href={`/emotes/${id}`}>
          <Card className='flex flex-col items-center p-4 transition-all duration-150 hover:shadow-md hover:shadow-purple-500/10'>
            <Heading
              className='text-xl text-purple-500 md:text-2xl'
              variant='h2'
            >
              {displayName}
            </Heading>
          </Card>
        </Link>
      </li>
    ))}
  </ul>
);

export default ChannelList;
