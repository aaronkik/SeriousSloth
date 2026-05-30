import type { Metadata } from 'next';
import ChannelList from '~/app/emotes/components/channel-list';
import { Heading } from '~/components/shared';
import { emotesTitle } from '~/constants/titles';
import { getCachedChannels } from '~/app/emotes/queries';

export const metadata: Metadata = {
  title: emotesTitle,
  description: 'Pick a channel to view its current Twitch emotes.',
};

const Page = async () => {
  const channels = await getCachedChannels();

  return (
    <>
      <div className='mb-6 flex flex-col items-center gap-2 text-center'>
        <Heading variant='h1'>Emotes</Heading>
        <p>Pick a channel to view its current Twitch emotes.</p>
      </div>
      <ChannelList channels={channels} />
    </>
  );
};

export default Page;
