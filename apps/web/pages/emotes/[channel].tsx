import { GetServerSidePropsContext, InferGetServerSidePropsType } from 'next';
import Head from 'next/head';
import { DynamicLastUpdated, EmoteTabs } from '~/components/emotes';
import { Heading } from '~/components/shared';
import {
  getActiveEmotes,
  getChannels,
  getRemovedEmotes,
} from '~/lib/api/emotes-service';
import { channelSlug, type Channel } from '~/lib/api/channels';

type Params = { channel: string };

export async function getServerSideProps(
  ctx: GetServerSidePropsContext<Params>
) {
  ctx.res.setHeader(
    'Cache-Control',
    'public, s-maxage=300, stale-while-revalidate'
  );

  const channelParam = ctx.params?.channel;

  if (!channelParam) {
    return { notFound: true } as const;
  }

  const channels = await getChannels();
  const channel = channels.find((c) => channelSlug(c) === channelParam);

  if (!channel) {
    return { notFound: true } as const;
  }

  const [activeEmotes, removedEmotes] = await Promise.all([
    getActiveEmotes(channel.id),
    getRemovedEmotes(channel.id),
  ]);

  return {
    props: {
      channel,
      activeEmotes,
      removedEmotes,
      updatedAt: Date.now(),
    },
  };
}

const ChannelEmotesPage = ({
  channel,
  activeEmotes,
  removedEmotes,
  updatedAt,
}: InferGetServerSidePropsType<typeof getServerSideProps> & {
  channel: Channel;
}) => (
  <>
    <Head>
      <title>{`${channel.displayName} Emotes | SeriousSloth`}</title>
    </Head>
    <div className='mb-2 flex flex-col items-center gap-2 text-center'>
      <Heading variant='h1'>{`${channel.displayName} Emotes`}</Heading>
      <DynamicLastUpdated lastUpdated={updatedAt} />
    </div>
    <EmoteTabs activeEmotes={activeEmotes} removedEmotes={removedEmotes} />
  </>
);

export default ChannelEmotesPage;
