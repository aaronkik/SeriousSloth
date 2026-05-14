import { getEnvOrThrow } from '~/lib/helpers';

const apiUrl = getEnvOrThrow('EMOTES_SERVICE_API_URL');
const apiKey = getEnvOrThrow('EMOTES_SERVICE_API_KEY');

export interface ActiveEmote {
  emote: {
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
  };
  addedAt: string;
}

export const getActiveEmotes = async (
  channel: string
): Promise<ActiveEmote[]> => {
  const url = `${apiUrl}/emotes/${channel}`;
  const res = await fetch(url, { headers: { 'x-api-key': apiKey } });

  if (!res.ok) {
    const body = await res.text();
    throw new Error(
      `getActiveEmotes ${url} failed: ${res.status} ${res.statusText} — ${body}`
    );
  }

  return res.json();
};

export interface Channel {
  id: string;
  displayName: string;
}

// TODO: replace body with `fetch(`${apiUrl}/channels`, { headers: { 'x-api-key': apiKey } })`
// once the emotes-service exposes a list-channels endpoint.
export const getChannels = async (): Promise<Channel[]> => [
  { id: 'global', displayName: 'Global' },
];
