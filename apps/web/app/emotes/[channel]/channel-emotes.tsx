import { notFound } from 'next/navigation';
import DynamicLastUpdated from '~/app/emotes/components/dynamic-last-updated';
import EmoteTabs from '~/app/emotes/components/emote-tabs';
import { Heading } from '~/components/shared';
import { getChannelEmotes } from '~/app/emotes/[channel]/queries';

type Props = {
  params: Promise<{ channel: string }>;
};

const ChannelEmotes = async ({ params }: Props) => {
  const { channel } = await params;
  const data = await getChannelEmotes(channel);

  if (!data) {
    notFound();
  }

  const { channel: emotesChannel, activeEmotes, removedEmotes, updatedAt } =
    data;

  return (
    <>
      <div className='mb-2 flex flex-col items-center gap-2 text-center'>
        <Heading variant='h1'>{`${emotesChannel.displayName} Emotes`}</Heading>
        <DynamicLastUpdated lastUpdated={updatedAt} />
      </div>
      <EmoteTabs activeEmotes={activeEmotes} removedEmotes={removedEmotes} />
    </>
  );
};

export default ChannelEmotes;
