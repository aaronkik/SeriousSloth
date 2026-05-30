import type { Metadata } from 'next';
import { Suspense } from 'react';
import ChannelEmotes from '~/app/emotes/[channel]/channel-emotes';
import ChannelEmotesSkeleton from '~/app/emotes/[channel]/channel-emotes-skeleton';
import { getChannel } from '~/app/emotes/[channel]/queries';
import { channelSlug } from '~/lib/api/channels';
import { getChannels } from '~/lib/api/emotes-service';

type PageProps = {
  params: Promise<{ channel: string }>;
};

export async function generateStaticParams() {
  const channels = await getChannels();

  return channels.map((channel) => ({ channel: channelSlug(channel) }));
}

export async function generateMetadata({
  params,
}: PageProps): Promise<Metadata> {
  const { channel } = await params;
  const found = await getChannel(channel);

  if (!found) {
    return {};
  }

  return { title: `${found.displayName} Emotes | SeriousSloth` };
}

const Page = ({ params }: PageProps) => (
  <Suspense fallback={<ChannelEmotesSkeleton />}>
    <ChannelEmotes params={params} />
  </Suspense>
);

export default Page;
