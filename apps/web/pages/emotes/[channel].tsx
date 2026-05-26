import { GetServerSidePropsContext, InferGetServerSidePropsType } from 'next';
import Head from 'next/head';
import { DynamicLastUpdated, EmoteTabs } from '~/components/emotes';
import { Heading } from '~/components/shared';
import {
  getActiveEmotes,
  getChannels,
  getRemovedEmotes,
  type ActiveEmoteEntry,
  type RemovedEmoteEntry,
} from '~/lib/api/emotes-service';
import { channelSlug, type Channel } from '~/lib/api/channels';
import { buildEmoteUrl } from '~/lib/helpers';

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

  const [rawActiveEmotes, rawRemovedEmotes] = await Promise.all([
    getActiveEmotes(channelParam),
    getRemovedEmotes(channelParam),
  ]);

  const activeEmotes: ActiveEmoteEntry[] = rawActiveEmotes.map(
    ({ emote, addedAt }) => ({
      id: emote.id,
      name: emote.name,
      emoteUrl: buildEmoteUrl(emote),
      addedAt,
    })
  );

  const removedEmotes: RemovedEmoteEntry[] = rawRemovedEmotes.map(
    ({ emote, removedAt }) => ({
      id: emote.id,
      name: emote.name,
      emoteUrl: buildEmoteUrl(emote),
      removedAt,
    })
  );

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
