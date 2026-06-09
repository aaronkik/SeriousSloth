import { notFound } from 'next/navigation';
import { Suspense } from 'react';
import EmoteTabsSection from '~/app/emotes/[channel]/emote-tabs-section';
import EmoteTabsSkeleton from '~/app/emotes/[channel]/emote-tabs-skeleton';
import { getChannel } from '~/app/emotes/[channel]/queries';
import { Heading } from '~/components/shared';

type Props = {
  params: Promise<{ channel: string }>;
};

const ChannelEmotes = async ({ params }: Props) => {
  const { channel } = await params;
  const found = await getChannel(channel);

  if (!found) {
    notFound();
  }

  return (
    <>
      <div className='mb-2 flex flex-col items-center gap-2 text-center'>
        <Heading variant='h1'>{`${found.displayName} Emotes`}</Heading>
      </div>
      <Suspense fallback={<EmoteTabsSkeleton />}>
        <EmoteTabsSection channelParam={channel} />
      </Suspense>
    </>
  );
};

export default ChannelEmotes;
