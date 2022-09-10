import Head from 'next/head';
import Image from 'next/future/image';
import { Navigation } from '~/components/home';
import { Heading, MutedText } from '~/components/shared';
import { homeTitle } from '~/constants/titles';

import sloth from '~/public/assets/sloth-face-square.png';
function HomePage() {
  return (
    <>
      <Head>
        <title>{homeTitle}</title>
      </Head>
      <div className='flex flex-col items-center gap-2 py-8'>
        <Image
          alt='Sloth face'
          src={sloth}
          width={80}
          height={80}
          placeholder='blur'
          className='rounded-full'
        />
        <Heading variant='h1'>SeriousSloth</Heading>
        <MutedText>A web app that interacts with the Twitch API</MutedText>
      </div>
      <Navigation />
    </>
  );
}

export default HomePage;
