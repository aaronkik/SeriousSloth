import Link from 'next/link';
import { Avatar, AvatarFallback, AvatarImage } from '~/components/ui/avatar';
import { Card } from '~/components/ui/card';
import { channelSlug, type Channel } from '~/lib/api/emotes-service';

type Props = {
  channels: Channel[];
};

const ChannelList = ({ channels }: Props) => (
  <ul
    data-testid='channelList'
    className='grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5'
  >
    {channels.map((channel) => (
      <li key={channel.id}>
        <Link
          href={`/emotes/${channelSlug(channel)}`}
          className='group block'
        >
          <Card className='flex flex-col items-center gap-3 p-5 ring-1 ring-transparent group-hover:ring-primary/60'>
            {channel.type === 'twitch' ? (
              <Avatar className='size-20 sm:size-24'>
                <AvatarImage
                  src={channel.imageUrl}
                  alt={`${channel.displayName} avatar`}
                />
                <AvatarFallback>
                  {channel.displayName.slice(0, 2)}
                </AvatarFallback>
              </Avatar>
            ) : (
              <div
                className='flex size-20 items-center justify-center text-5xl sm:size-24'
                aria-hidden
              >
                {channel.icon}
              </div>
            )}
            <p className='truncate text-base font-semibold text-foreground group-hover:text-primary'>
              {channel.displayName}
            </p>
          </Card>
        </Link>
      </li>
    ))}
  </ul>
);

export default ChannelList;
