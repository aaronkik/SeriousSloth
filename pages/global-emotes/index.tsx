import { InferGetStaticPropsType } from 'next';
import Head from 'next/head';
import {
  DynamicLastUpdated,
  GlobalEmotesList,
} from '~/components/global-emotes';
import { Heading } from '~/components/shared';
import { globalEmotesTitle } from '~/constants/titles';
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
}: InferGetStaticPropsType<typeof getStaticProps>) => {
  return (
    <>
      <Head>
        <title>{globalEmotesTitle}</title>
      </Head>
      <div className='py-8 text-center flex items-center flex-col gap-2'>
        <Heading variant='h1'>Global Emotes</Heading>
        <DynamicLastUpdated lastUpdated={updatedAt} />
      </div>
      <GlobalEmotesList globalEmotes={globalEmotes} />
    </>
  );
};

export default GlobalEmotesPage;
