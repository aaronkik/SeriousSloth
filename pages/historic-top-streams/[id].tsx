import { GetStaticPaths, GetStaticProps, InferGetStaticPropsType } from 'next';
import Head from 'next/head';
import NextLink from 'next/link';
import { TopStreamCard } from '~/components/historic-top-streams';
import { Heading, Link, MutedText } from '~/components/shared';
import { historicTopStreamsTitle } from '~/constants/titles';
import { formatDate } from '~/lib/helpers';
import { supabase } from '~/lib/supabase';
import { definitions } from '~/types/supabase';

export const getStaticPaths: GetStaticPaths = async () => {
  const { data } = await supabase
    .from<definitions['stream_history_timestamp']>('stream_history_timestamp')
    .select('id')
    .order('time', { ascending: false })
    .limit(50);

  return {
    fallback: 'blocking',
    paths: data ? data.map(({ id }) => ({ params: { id: String(id) } })) : [],
  };
};

type StreamHistoryData = definitions['stream_history_timestamp'] & {
  stream_history: Array<definitions['stream_history']>;
};

export const getStaticProps: GetStaticProps<{
  data: StreamHistoryData;
}> = async ({ params }) => {
  const paramId = params?.id;

  if (!paramId || typeof paramId !== 'string') {
    return { notFound: true };
  }

  try {
    const { data } = await supabase
      .from<StreamHistoryData>('stream_history_timestamp')
      .select('time,stream_history(*)')
      .eq('id', paramId)
      // @ts-ignore - false flag
      .order('viewer_count', {
        foreignTable: 'stream_history',
        ascending: false,
      })
      .single();

    if (data) {
      return { props: { data } };
    }

    return { notFound: true };
  } catch (error) {
    console.error(error);
    return { notFound: true };
  }
};

const HistoricTopStreamsIdPage = (
  props: InferGetStaticPropsType<typeof getStaticProps>
) => {
  const {
    data: { time, stream_history },
  } = props;

  return (
    <>
      <Head>
        <title>{historicTopStreamsTitle}</title>
      </Head>
      <div className='mb-2 flex flex-col items-center gap-2 text-center'>
        <Heading variant='h1'>Historic Top Streams</Heading>
        <MutedText>
          The top {stream_history.length} streams at {formatDate(time, 'LLL')}
        </MutedText>
        <NextLink href='/historic-top-streams' passHref>
          <Link className='text-sm' data-testid='historicTopStreamsLink'>
            Back to stream times
          </Link>
        </NextLink>
      </div>
      <div className='flex justify-center py-8'>
        <ul
          className='flex w-full max-w-2xl flex-col gap-6'
          data-testid='topStreams'
        >
          {stream_history.map((stream) => (
            <li key={stream.id}>
              <TopStreamCard stream={stream} />
            </li>
          ))}
        </ul>
      </div>
    </>
  );
};

export default HistoricTopStreamsIdPage;
