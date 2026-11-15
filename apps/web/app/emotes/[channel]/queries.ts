import { cacheLife } from 'next/cache';
import { channelSlug } from '~/lib/api/channels';
import {
  type EmoteListEntry,
  type EmotesByDate,
  getActiveEmotes,
  getAddedEmotes,
  getChannels,
  getRemovedEmotes,
} from '~/lib/api/emotes-service';
import { groupEmotesByDate, sortEmotesByName } from '~/lib/helpers';
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
  activeEmotes: EmoteListEntry[];
  addedEmotes: EmotesByDate;
  removedEmotes: EmotesByDate;
}> => {
  'use cache';
  cacheLife({ stale: 300, revalidate: 300, expire: 300 });

  setAttributes({ 'emotes.channel.searchQuery': channelParam });

  const [rawActiveEmotes, rawAddedEmotes, rawRemovedEmotes] = await Promise.all(
    [
      trace({
        name: 'getEmoteData/getActiveEmotes',
        handler: () => getActiveEmotes(channelParam),
      }),
      trace({
        name: 'getEmoteData/getAddedEmotes',
        handler: () => getAddedEmotes(channelParam),
      }),
      trace({
        name: 'getEmoteData/getRemovedEmotes',
        handler: () => getRemovedEmotes(channelParam),
      }),
    ],
  );

  setAttributes({
    'emotes.count.active': rawActiveEmotes.length,
    'emotes.count.added': rawAddedEmotes.length,
    'emotes.count.remote': rawRemovedEmotes.length,
  });

  return {
    activeEmotes: sortEmotesByName(rawActiveEmotes.map(({ emote }) => emote)),
    addedEmotes: groupEmotesByDate(rawAddedEmotes, 'addedAt'),
    removedEmotes: groupEmotesByDate(rawRemovedEmotes, 'removedAt'),
  };
};
