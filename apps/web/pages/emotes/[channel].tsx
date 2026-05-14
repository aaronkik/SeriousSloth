import {
  GetStaticPaths,
  GetStaticPropsContext,
  InferGetStaticPropsType,
} from 'next';
import Head from 'next/head';
import { DynamicLastUpdated, EmotesList } from '~/components/emotes';
import { Heading } from '~/components/shared';
import {
  Channel,
  getActiveEmotes,
  getChannels,
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

  const activeEmotes = await getActiveEmotes(channel.id);

  return {
    props: {
      channel,
      activeEmotes,
      updatedAt: Date.now(),
    },
    revalidate: 60 * 60,
  };
}

const ChannelEmotesPage = ({
  channel,
  activeEmotes,
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
    <EmotesList emotes={activeEmotes} />
  </>
);

export default ChannelEmotesPage;
