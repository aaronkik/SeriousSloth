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
import { setAttributes, trace } from '~/observability';

export const getChannel = async (channelParam: string) => {
  'use cache';
  cacheLife({ stale: 300, revalidate: 300, expire: 300 });

  const channels = await trace({
    name: 'getChannel/getChannels',
    handler: () => getChannels(),
    attributes: { 'channel.searchQuery': channelParam },
  });

  setAttributes({ 'channels.count': channels.length });

  return channels.find((c) => channelSlug(c) === channelParam) ?? null;
};

export const getEmoteData = async (channelParam: string) => {
  'use cache';
  cacheLife({ stale: 300, revalidate: 300, expire: 300 });

  setAttributes({ 'emotes.channel.searchQuery': channelParam });

  const [rawActiveEmotes, rawRemovedEmotes] = await Promise.all([
    trace({
      name: 'getEmoteData/getActiveEmotes',
      handler: () => getActiveEmotes(channelParam),
    }),
    trace({
      name: 'getEmoteData/getRemovedEmotes',
      handler: () => getRemovedEmotes(channelParam),
    }),
  ]);

  setAttributes({
    'emotes.count.active': rawActiveEmotes.length,
    'emotes.count.remote': rawRemovedEmotes.length,
  });

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
