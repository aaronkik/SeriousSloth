import Head from 'next/head';
import { Navigation } from '~/components/home';
import { Heading, MutedText } from '~/components/shared';
import SlothLogo from '~/components/shared/sloth-logo';
import { homeTitle } from '~/constants/titles';

function HomePage() {
  return (
    <>
      <Head>
        <title>{homeTitle}</title>
      </Head>
      <div className='flex flex-col items-center gap-2 py-8'>
        <SlothLogo className='rounded-full' width={80} height={80} />
        <Heading variant='h1'>SeriousSloth</Heading>
        <MutedText>A web app that interacts with the Twitch API</MutedText>
      </div>
      <Navigation />
    </>
  );
}

export default HomePage;
