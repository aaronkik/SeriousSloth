import { InferGetStaticPropsType } from 'next';
import Head from 'next/head';
import sharp from 'sharp';
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
   * Create base64 image urls to make use of the nextjs Image component
   */
  const globalEmotes = await Promise.all(
    data.map(async (emote) => {
      const { id, name } = emote;

      const largeImageUrl = formatEmoteCDNUrl(template, {
        id,
        format: 'default',
        theme_mode: 'dark',
        scale: '3.0',
      });

      const smallImageUrl = formatEmoteCDNUrl(template, {
        id,
        format: 'default',
        theme_mode: 'dark',
        scale: '1.0',
      });

      const response = await fetch(smallImageUrl, {
        method: 'GET',
      });
      const contentType = response.headers.get('Content-Type');
      const imageBuffer = await response.arrayBuffer();
      const resizedBuffer = await sharp(Buffer.from(imageBuffer))
        .resize(5, 5)
        .toBuffer();
      const base64Image = Buffer.from(resizedBuffer).toString('base64');
      // https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/Data_URLs
      const blurDataUrl = `data:${contentType};base64,${base64Image}`;

      return {
        id,
        name,
        largeImageUrl,
        blurDataUrl,
      };
    })
  );

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
    <div className='flex flex-col items-center gap-2 py-8 text-center'>
      <Heading variant='h1'>Global Emotes</Heading>
      <DynamicLastUpdated lastUpdated={updatedAt} />
    </div>
    <GlobalEmotesList globalEmotes={globalEmotes} />
  </>
);

export default GlobalEmotesPage;
