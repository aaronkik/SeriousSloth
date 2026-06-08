import { notFound } from 'next/navigation';
import { Suspense } from 'react';
import EmoteTabsSection from '~/app/emotes/[channel]/emote-tabs-section';
import EmoteTabsSkeleton from '~/app/emotes/[channel]/emote-tabs-skeleton';
import { getChannel } from '~/app/emotes/[channel]/queries';
import { Heading, MutedText } from '~/components/shared';
import { timeFromNow } from '~/lib/helpers';

type Props = {
  params: Promise<{ channel: string }>;
};

const ChannelEmotes = async ({ params }: Props) => {
  const { channel } = await params;
  const found = await getChannel(channel);

  if (!found) {
    notFound();
  }

  const currentTime = new Date().toISOString();

  return (
    <>
      <div className='mb-2 flex flex-col items-center gap-2 text-center'>
        <Heading variant='h1'>{`${found.displayName} Emotes`}</Heading>
        <MutedText className='text-sm'>
          Last updated: {timeFromNow(currentTime)}
        </MutedText>
      </div>
      <Suspense fallback={<EmoteTabsSkeleton />}>
        <EmoteTabsSection channelParam={channel} />
      </Suspense>
    </>
  );
};

export default ChannelEmotes;
