import { cacheLife } from 'next/cache';
import { channelSlug } from '~/lib/api/channels';
import {
  type ActiveEmoteEntry,
  getActiveEmotes,
  getChannels,
  getRemovedEmotes,
  type RemovedEmoteEntry,
} from '~/lib/api/emotes-service';
import { buildEmoteUrl } from '~/lib/helpers';

export const getChannel = async (channelParam: string) => {
  'use cache';
  cacheLife({ stale: 300, revalidate: 300, expire: 300 });

  const channels = await getChannels();

  return channels.find((c) => channelSlug(c) === channelParam) ?? null;
};

export const getEmoteData = async (channelParam: string) => {
  'use cache';
  cacheLife({ stale: 300, revalidate: 300, expire: 300 });

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
    }),
  );

  const removedEmotes: RemovedEmoteEntry[] = rawRemovedEmotes.map(
    ({ emote, removedAt }) => ({
      id: emote.id,
      name: emote.name,
      emoteUrl: buildEmoteUrl(emote),
      removedAt,
    }),
  );

  return { activeEmotes, removedEmotes, updatedAt: Date.now() };
};
