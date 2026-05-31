import type { Metadata } from 'next';
import { Heading, MutedText } from '~/components/shared';
import SlothLogo from './sloth-logo';
import { Navigation } from './navigation';

export const metadata: Metadata = {
  title: 'SeriousSloth',
  description: 'A web app that interacts with the Twitch API',
};

const Page = () => {
  return (
    <>
      <div className='mb-12 flex flex-col items-center gap-2'>
        <SlothLogo className='rounded-full' width={80} height={80} />
        <Heading variant='h1'>SeriousSloth</Heading>
        <MutedText>A web app that interacts with the Twitch API</MutedText>
      </div>
      <Navigation />
    </>
  );
};

export default Page;
