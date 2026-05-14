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
  fetchEmotes(`${apiUrl}/emotes/${channel}`, 'getActiveEmotes');

export const getRemovedEmotes = (channel: string): Promise<RemovedEmote[]> =>
  fetchEmotes(`${apiUrl}/emotes/${channel}/removed`, 'getRemovedEmotes');

export interface Channel {
  id: string;
  displayName: string;
}

// TODO: replace body with `fetch(`${apiUrl}/channels`, { headers: { 'x-api-key': apiKey } })`
// once the emotes-service exposes a list-channels endpoint.
export const getChannels = async (): Promise<Channel[]> => [
  { id: 'global', displayName: 'Global' },
];
