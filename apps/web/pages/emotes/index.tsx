import { GetServerSidePropsContext, InferGetServerSidePropsType } from 'next';
import Head from 'next/head';
import { ChannelList } from '~/components/emotes';
import { Heading } from '~/components/shared';
import { emotesTitle } from '~/constants/titles';
import { getChannels } from '~/lib/api/emotes-service';

export async function getServerSideProps(ctx: GetServerSidePropsContext) {
  ctx.res.setHeader(
    'Cache-Control',
    'public, s-maxage=300, stale-while-revalidate'
  );

  const channels = await getChannels();

  return { props: { channels } };
}

const EmotesPage = ({
  channels,
}: InferGetServerSidePropsType<typeof getServerSideProps>) => (
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
