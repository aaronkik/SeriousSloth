import { InferGetStaticPropsType } from 'next';
import Head from 'next/head';
import Image from 'next/image';
import { Heading, MutedText } from '~/components/shared';
import { globalEmotesTitle } from '~/constants/titles';
import { timeFromNow } from '~/lib/helpers';
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
      updatedAt: Date.now(),
    },
    revalidate: 60 * 60 * 8, // 8 hours in seconds
  };
}

const GlobalEmotesPage = ({
  globalEmotes,
  updatedAt,
}: InferGetStaticPropsType<typeof getStaticProps>) => (
  <>
    <Head>
      <title>{globalEmotesTitle}</title>
    </Head>
    <div className='py-8 text-center flex flex-col gap-2'>
      <Heading variant='h1'>Global Emotes</Heading>
      <MutedText className='text-sm'>
        Last updated: {timeFromNow(updatedAt)}
      </MutedText>
    </div>
    <ul
      data-testid='globalEmoteList'
      className='flex flex-row gap-4 justify-center flex-wrap py-8'
    >
      {globalEmotes.map((emote, index) => (
        <li
          key={emote.id}
          className='bg-neutral-800 rounded shadow w-40 h-40 p-4'
        >
          <div className='flex flex-col gap-2 w-full h-full items-center'>
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
            <p className='tracking-wide' data-testid={`emoteName${index}`}>
              {emote.name}
            </p>
          </div>
        </li>
      ))}
    </ul>
  </>
);

export default GlobalEmotesPage;
