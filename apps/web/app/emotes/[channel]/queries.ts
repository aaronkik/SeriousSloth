import { cacheLife } from 'next/cache';

const REVALIDATE_SECONDS = 300;
import { channelSlug } from '~/lib/api/channels';
import {
  getActiveEmotes,
  getChannels,
  getRemovedEmotes,
  type ActiveEmoteEntry,
  type RemovedEmoteEntry,
} from '~/lib/api/emotes-service';
import { buildEmoteUrl } from '~/lib/helpers';

export const getChannel = async (channelParam: string) => {
  'use cache: remote';
  cacheLife({ revalidate: REVALIDATE_SECONDS });

  const channels = await getChannels();

  return channels.find((c) => channelSlug(c) === channelParam) ?? null;
};

export const getEmoteData = async (channelParam: string) => {
  'use cache: remote';
  cacheLife({ revalidate: REVALIDATE_SECONDS });

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

  // Align to the revalidate window so the value is identical for every
  // request served from the same cache entry (i.e. when the emotes were
  // last fetched), rather than tracking the current request time.
  const windowMs = REVALIDATE_SECONDS * 1000;
  const updatedAt = Math.floor(Date.now() / windowMs) * windowMs;

  return { activeEmotes, removedEmotes, updatedAt };
};
