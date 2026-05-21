import { InferGetStaticPropsType } from 'next';
import Head from 'next/head';
import { ChannelList } from '~/components/emotes';
import { Heading } from '~/components/shared';
import { emotesTitle } from '~/constants/titles';
import { getChannelListing } from '~/lib/api/emotes-service';

export async function getStaticProps() {
  const channels = await getChannelListing();

  return {
    props: { channels },
    revalidate: 60 * 60,
  };
}

const EmotesPage = ({
  channels,
}: InferGetStaticPropsType<typeof getStaticProps>) => (
  <>
    <Head>
      <title>{emotesTitle}</title>
    </Head>
    <div className='mb-6 flex flex-col items-center gap-2 text-center'>
      <Heading variant='h1'>Emotes</Heading>
      <p>Pick a channel to view its current Twitch emotes.</p>
    </div>
    <ChannelList channels={channels} />
  </>
);

export default EmotesPage;
