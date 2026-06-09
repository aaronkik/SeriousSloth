import { getEnvOrThrow } from '~/lib/helpers';
import {
  GLOBAL_CHANNEL,
  type Channel,
  type TwitchChannel,
} from '~/lib/api/channels';
import { get } from '~/lib/fetch';

export type { Channel } from '~/lib/api/channels';

const apiUrl = getEnvOrThrow('EMOTES_SERVICE_API_URL');
const apiKey = getEnvOrThrow('EMOTES_SERVICE_API_KEY');

export interface Emote {
  id: string;
  name: string;
  images: {
    url_1x: string;
    url_2x: string;
    url_4x: string;
  };
  format: Array<'static' | 'animated'>;
  scale: Array<'1.0' | '2.0' | '3.0'>;
  theme_mode: Array<'dark' | 'light'>;
}

export interface ActiveEmote {
  emote: Emote;
  addedAt: string;
}

export interface RemovedEmote {
  emote: Emote;
  removedAt: string;
}

export interface EmoteListEntry {
  id: string;
  name: string;
  emoteUrl: string;
}

export interface ActiveEmoteEntry extends EmoteListEntry {
  addedAt: string;
}

export interface RemovedEmoteEntry extends EmoteListEntry {
  removedAt: string;
}

export const getActiveEmotes = (channel: string): Promise<ActiveEmote[]> =>
  get(`${apiUrl}/v1/emotes/${channel}`, {
    headers: { 'x-api-key': apiKey },
  });

export const getRemovedEmotes = (channel: string): Promise<RemovedEmote[]> =>
  get(`${apiUrl}/v1/emotes/${channel}/removed`, {
    headers: { 'x-api-key': apiKey },
  });

interface ChannelDto {
  id: string;
  twitchId: string;
  displayName: string;
  imageUrl: string;
}

export const getChannels = async (): Promise<Channel[]> => {
  const channels = await get<ChannelDto[]>(`${apiUrl}/v1/channels`, {
    headers: { 'x-api-key': apiKey },
  });

  const twitchChannels: TwitchChannel[] = channels.map((c) => ({
    type: 'twitch',
    id: c.id,
    twitchId: c.twitchId,
    displayName: c.displayName,
    imageUrl: c.imageUrl,
  }));

  const sortedTwitchChannels = twitchChannels.toSorted((a, b) =>
    a.displayName.localeCompare(b.displayName),
  );

  return [GLOBAL_CHANNEL, ...sortedTwitchChannels];
};
