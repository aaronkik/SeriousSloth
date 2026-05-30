import type { Metadata } from 'next';
import { Suspense } from 'react';
import ChannelEmotes from '~/app/emotes/[channel]/channel-emotes';
import ChannelEmotesSkeleton from '~/app/emotes/[channel]/channel-emotes-skeleton';
import { getChannelEmotes } from '~/app/emotes/[channel]/queries';

type PageProps = {
  params: Promise<{ channel: string }>;
};

export async function generateMetadata({
  params,
}: PageProps): Promise<Metadata> {
  const { channel } = await params;
  const data = await getChannelEmotes(channel);

  if (!data) {
    return {};
  }

  return { title: `${data.channel.displayName} Emotes | SeriousSloth` };
}

const Page = ({ params }: PageProps) => (
  <Suspense fallback={<ChannelEmotesSkeleton />}>
    <ChannelEmotes params={params} />
  </Suspense>
);

export default Page;
