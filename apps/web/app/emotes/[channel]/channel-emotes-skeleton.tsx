import EmoteTabsSkeleton from '~/app/emotes/[channel]/emote-tabs-skeleton';
import { Skeleton } from '~/components/shared';

const ChannelEmotesSkeleton = () => (
  <>
    <div className='mb-2 flex flex-col items-center gap-2 text-center'>
      <Skeleton className='h-10 w-72 md:h-14' />
      <Skeleton className='h-5 w-60' />
    </div>
    <EmoteTabsSkeleton />
  </>
);

export default ChannelEmotesSkeleton;
