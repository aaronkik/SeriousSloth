import { cacheLife } from 'next/cache';
import { channelSlug } from '~/lib/api/channels';
import {
  getActiveEmotes,
  getChannels,
  getRemovedEmotes,
  type ActiveEmoteEntry,
  type RemovedEmoteEntry,
} from '~/lib/api/emotes-service';
import { buildEmoteUrl } from '~/lib/helpers';

export const getChannelEmotes = async (channelParam: string) => {
  'use cache';
  cacheLife({ revalidate: 300 });

  const channels = await getChannels();
  const channel = channels.find((c) => channelSlug(c) === channelParam);

  if (!channel) {
    return null;
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

  return { channel, activeEmotes, removedEmotes, updatedAt: Date.now() };
};
