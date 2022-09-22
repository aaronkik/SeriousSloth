import Head from 'next/head';
import { TopStreamTimes } from '~/components/historic-top-streams';
import { Heading } from '~/components/shared';
import { historicTopStreamsTitle } from '~/constants/titles';

const HistoricTopStreamsPage = () => {
  return (
    <>
      <Head>
        <title>{historicTopStreamsTitle}</title>
      </Head>
      <Heading
        className='mb-2 flex flex-col items-center gap-2 text-center'
        variant='h1'
      >
        Historic Top Streams
      </Heading>
      <div className='flex justify-center py-8'>
        <TopStreamTimes />
      </div>
    </>
  );
};

export default HistoricTopStreamsPage;
