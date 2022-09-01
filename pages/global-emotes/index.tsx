import { InferGetStaticPropsType } from 'next';
import Image from 'next/image';
import {
  clientId,
  clientSecret,
  globalEmotesEndpoint,
  oauth2TokenEndpoint,
} from '~/constants/twitch';
import { formatEmoteCDNUrl } from '~/lib/twitch';
import {
  GlobalEmotesResponse,
  OAuthClientCredentialsResponse,
} from '~/types/twitch';

export async function getStaticProps() {
  const twitchAccessTokenResponse = await fetch(oauth2TokenEndpoint, {
    method: 'POST',
    headers: new Headers({
      'Content-Type': 'application/x-www-form-urlencoded',
    }),
    body: `client_id=${clientId}&client_secret=${clientSecret}&grant_type=client_credentials`,
  });

  if (!twitchAccessTokenResponse.ok) {
    throw new Error('Failed to get access_token');
  }

  const { access_token } =
    (await twitchAccessTokenResponse.json()) as OAuthClientCredentialsResponse;

  const globalEmotesResponse = await fetch(globalEmotesEndpoint, {
    method: 'GET',
    headers: new Headers({
      Authorization: `Bearer ${access_token}`,
      'Client-Id': clientId,
    }),
  });

  if (!globalEmotesResponse.ok) {
    throw new Error('Failed to get global emotes');
  }

  const { data, template } =
    (await globalEmotesResponse.json()) as GlobalEmotesResponse;

  /**
   * Reduce JSON payload and only send the information needed
   */
  const globalEmotes = data.map((emote) => {
    const { id, name } = emote;
    return {
      id,
      name,
      image: formatEmoteCDNUrl(template, {
        id,
        format: 'default',
        theme_mode: 'dark',
        scale: '3.0',
      }),
    };
  });

  return {
    props: {
      globalEmotes,
    },
  };
}

const GlobalEmotesPage = (
  props: InferGetStaticPropsType<typeof getStaticProps>
) => (
  <ul data-testid='globalEmoteList' className='flex flex-row gap-4 flex-wrap'>
    {props.globalEmotes.map((emote, index) => (
      <li key={emote.id} className='bg-gray-300 w-40 h-40 p-8'>
        <div className='flex flex-col gap-2 w-full h-full'>
          <div className='relative w-full h-full'>
            <Image
              alt={`${emote.name} emote`}
              src={emote.image}
              layout='fill'
              objectFit='contain'
              data-testid={`emoteImage${index}`}
            />
          </div>
          <p data-testid={`emoteName${index}`}>{emote.name}</p>
        </div>
      </li>
    ))}
  </ul>
);

export default GlobalEmotesPage;
