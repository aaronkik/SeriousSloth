import { cacheLife } from 'next/cache';
import { channelSlug } from '~/lib/api/channels';
import {
  type ActiveEmoteEntry,
  type RemovedEmoteEntry,
  getActiveEmotes,
  getChannels,
  getRemovedEmotes,
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

export const getEmoteData = async (
  channelParam: string,
): Promise<{
  activeEmotes: Record<string, ActiveEmoteEntry[]>;
  activeEmotesCount: number;
  removedEmotes: Record<string, RemovedEmoteEntry[]>;
  removedEmotesCount: number;
}> => {
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

  const activeEmotesCount = rawActiveEmotes.length;
  const removedEmotesCount = rawRemovedEmotes.length;

  setAttributes({
    'emotes.count.active': activeEmotesCount,
    'emotes.count.remote': removedEmotesCount,
  });

  const activeEmotes = rawActiveEmotes.reduce(
    (prev, curr) => {
      const emoteAddedAtDate = new Date(curr.addedAt)
        .toISOString()
        .split('T')[0];

      (prev[emoteAddedAtDate] ??= []).push({
        id: curr.emote.id,
        name: curr.emote.name,
        emoteUrl: buildEmoteUrl(curr.emote),
        addedAt: curr.addedAt,
      });

      return prev;
    },
    {} as Record<string, ActiveEmoteEntry[]>,
  );

  const removedEmotes = rawRemovedEmotes.reduce(
    (prev, curr) => {
      const emoteRemovedDate = new Date(curr.removedAt)
        .toISOString()
        .split('T')[0];

      (prev[emoteRemovedDate] ??= []).push({
        id: curr.emote.id,
        name: curr.emote.name,
        emoteUrl: buildEmoteUrl(curr.emote),
        removedAt: curr.removedAt,
      });

      return prev;
    },
    {} as Record<string, RemovedEmoteEntry[]>,
  );

  const sortedActiveEmotes = Object.fromEntries(
    Object.entries(activeEmotes)
      .sort(([a], [b]) => b.localeCompare(a))
      .map(([date, emotes]) => [
        date,
        emotes.sort((a, b) => a.name.localeCompare(b.name)),
      ]),
  );

  const sortedRemovedEmotes = Object.fromEntries(
    Object.entries(removedEmotes)
      .sort(([a], [b]) => b.localeCompare(a))
      .map(([date, emotes]) => [
        date,
        emotes.sort((a, b) => a.name.localeCompare(b.name)),
      ]),
  );

  return {
    activeEmotes: sortedActiveEmotes,
    activeEmotesCount,
    removedEmotes: sortedRemovedEmotes,
    removedEmotesCount,
  };
};
