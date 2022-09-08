import Head from 'next/head';
import Image from 'next/image';
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
      <div className='py-8 flex flex-col gap-2 items-center'>
        <div className='w-20 h-20 relative rounded-full'>
          <Image
            alt='Sloth face'
            src={sloth}
            layout='fill'
            placeholder='blur'
            className='rounded-full'
          />
        </div>
        <Heading variant='h1'>SeriousSloth</Heading>
        <MutedText>A web app that interacts with the Twitch API</MutedText>
      </div>
      <Navigation />
    </>
  );
}

export default HomePage;
