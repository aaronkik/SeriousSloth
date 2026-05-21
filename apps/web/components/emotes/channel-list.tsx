import Image from 'next/image';
import Link from 'next/link';
import { Card } from '~/components/shared';
import type { Channel } from '~/lib/api/emotes-service';

type Props = {
  channels: Channel[];
};

const ChannelList = ({ channels }: Props) => (
  <ul
    data-testid='channelList'
    className='grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5'
  >
    {channels.map(({ id, displayName, profileImageUrl, icon }) => (
      <li key={id}>
        <Link href={`/emotes/${id}`} className='group block'>
          <Card className='flex flex-col items-center gap-3 p-5 ring-1 ring-transparent group-hover:ring-purple-500/60'>
            <div className='relative h-20 w-20 sm:h-24 sm:w-24'>
              {profileImageUrl ? (
                <Image
                  src={profileImageUrl}
                  alt={`${displayName} avatar`}
                  width={96}
                  height={96}
                  className='h-full w-full rounded-full object-cover ring-2 ring-neutral-700'
                />
              ) : icon ? (
                <div
                  className='flex h-full w-full items-center justify-center text-5xl'
                  aria-hidden
                >
                  {icon}
                </div>
              ) : null}
            </div>
            <p className='truncate text-base font-semibold text-neutral-100 group-hover:text-purple-300'>
              {displayName}
            </p>
          </Card>
        </Link>
      </li>
    ))}
  </ul>
);

export default ChannelList;
