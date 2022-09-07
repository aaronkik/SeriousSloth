import { InferGetStaticPropsType } from 'next';
import Image from 'next/image';
import {
  fetchClientCredentials,
  fetchGlobalEmotes,
  formatEmoteCDNUrl,
} from '~/lib/twitch';

export async function getStaticProps() {
  const { access_token } = await fetchClientCredentials();
  const { data, template } = await fetchGlobalEmotes(access_token);

  /**
   * Reduce JSON payload/only send neccessary emote information
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
              placeholder='blur'
              blurDataURL={emote.image}
            />
          </div>
          <p data-testid={`emoteName${index}`}>{emote.name}</p>
        </div>
      </li>
    ))}
  </ul>
);

export default GlobalEmotesPage;
