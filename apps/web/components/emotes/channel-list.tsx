import Link from 'next/link';
import { Avatar, AvatarFallback, AvatarImage } from '~/components/ui/avatar';
import { Card } from '~/components/ui/card';
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
          <Card className='flex flex-col items-center gap-3 p-5 ring-1 ring-transparent transition-colors group-hover:ring-primary/60'>
            {profileImageUrl ? (
              <Avatar className='size-20 sm:size-24'>
                <AvatarImage src={profileImageUrl} alt={`${displayName} avatar`} />
                <AvatarFallback>{displayName.slice(0, 2)}</AvatarFallback>
              </Avatar>
            ) : icon ? (
              <div className='flex size-20 items-center justify-center text-5xl sm:size-24' aria-hidden>
                {icon}
              </div>
            ) : null}
            <p className='truncate text-base font-semibold text-foreground transition-colors group-hover:text-primary'>
              {displayName}
            </p>
          </Card>
        </Link>
      </li>
    ))}
  </ul>
);

export default ChannelList;
