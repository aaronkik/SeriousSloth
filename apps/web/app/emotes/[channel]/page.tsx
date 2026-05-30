import type { Metadata } from 'next';
import { Suspense } from 'react';
import ChannelEmotes from '~/app/emotes/[channel]/channel-emotes';
import ChannelEmotesSkeleton from '~/app/emotes/[channel]/channel-emotes-skeleton';
import { getChannel } from '~/app/emotes/[channel]/queries';

type PageProps = {
  params: Promise<{ channel: string }>;
};

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
