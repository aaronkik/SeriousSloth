import { getEnvOrThrow } from '~/lib/helpers';

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

const fetchEmotes = async <T>(url: string, fn: string): Promise<T> => {
  const res = await fetch(url, { headers: { 'x-api-key': apiKey } });

  if (!res.ok) {
    const body = await res.text();
    throw new Error(
      `${fn} ${url} failed: ${res.status} ${res.statusText} — ${body}`
    );
  }

  return res.json();
};

export const getActiveEmotes = (channel: string): Promise<ActiveEmote[]> =>
  fetchEmotes(`${apiUrl}/v1/emotes/${channel}`, 'getActiveEmotes');

export const getRemovedEmotes = (channel: string): Promise<RemovedEmote[]> =>
  fetchEmotes(`${apiUrl}/v1/emotes/${channel}/removed`, 'getRemovedEmotes');

export interface GlobalChannel {
  type: 'global';
  id: 'global';
  displayName: string;
  icon: string;
}

export interface TwitchChannel {
  type: 'twitch';
  id: string;
  twitchId: string;
  displayName: string;
  imageUrl: string;
}

export type Channel = GlobalChannel | TwitchChannel;

export const channelSlug = (channel: Channel): string =>
  channel.type === 'global' ? 'global' : channel.twitchId;

interface ChannelDto {
  id: string;
  twitchId: string;
  displayName: string;
  imageUrl: string;
}

const GLOBAL_CHANNEL: GlobalChannel = {
  type: 'global',
  id: 'global',
  displayName: 'Global',
  icon: '🌐',
};

export const getChannels = async (): Promise<Channel[]> => {
  const channels = await fetchEmotes<ChannelDto[]>(
    `${apiUrl}/v1/channels`,
    'getChannels'
  );
  const twitchChannels: TwitchChannel[] = channels.map((c) => ({
    type: 'twitch',
    id: c.id,
    twitchId: c.twitchId,
    displayName: c.displayName,
    imageUrl: c.imageUrl,
  }));
  return [GLOBAL_CHANNEL, ...twitchChannels];
};
