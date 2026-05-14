import {
  GetStaticPaths,
  GetStaticPropsContext,
  InferGetStaticPropsType,
} from 'next';
import Head from 'next/head';
import { DynamicLastUpdated, EmoteTabs } from '~/components/emotes';
import { Heading } from '~/components/shared';
import {
  Channel,
  getActiveEmotes,
  getChannels,
  getRemovedEmotes,
} from '~/lib/api/emotes-service';

type Params = { channel: string };

export const getStaticPaths: GetStaticPaths<Params> = async () => {
  const channels = await getChannels();

  return {
    paths: channels.map(({ id }) => ({ params: { channel: id } })),
    fallback: 'blocking',
  };
};

export async function getStaticProps(ctx: GetStaticPropsContext<Params>) {
  const channelId = ctx.params?.channel;

  if (!channelId) {
    return { notFound: true } as const;
  }

  const channels = await getChannels();
  const channel = channels.find(({ id }) => id === channelId);

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
    revalidate: 60 * 60,
  };
}

const ChannelEmotesPage = ({
  channel,
  activeEmotes,
  removedEmotes,
  updatedAt,
}: InferGetStaticPropsType<typeof getStaticProps> & { channel: Channel }) => (
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
