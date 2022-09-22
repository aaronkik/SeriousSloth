import { useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import { supabase } from '~/lib/supabase';
import { StreamHistoryWithTotalViewCount } from '~/types/supabase/overrides';
import { Button, MutedText, Spinner } from '../shared';
import TopStreamTimesList from './top-stream-times-list';

const view = 'stream_history_timestamp_with_total_view_count';
const columns = 'id,time,total_viewer_count,total_streams';
const resultLimit = 50;

const TopStreams = () => {
  const [topStreamTimes, setTopStreamTimes] = useState<
    Array<StreamHistoryWithTotalViewCount> | undefined
  >(undefined);
  const [isFetchingMore, setIsFetchingMore] = useState(false);
  const [isMoreResults, setIsMoreResults] = useState(true);

  useEffect(() => {
    supabase
      .from<StreamHistoryWithTotalViewCount>(view)
      .select(columns)
      .order('time', { ascending: false })
      .limit(resultLimit)
      .then(({ data, error }) => {
        if (error) {
          console.error(error);
          toast.error(
            <div data-testid='streamHistoryError'>{error.message}</div>
          );
          return;
        }
        if (data) {
          setTopStreamTimes(data);
          if (data.length < resultLimit) {
            setIsMoreResults(false);
          }
        }
      });
  }, []);

  const fetchMoreTimes = async () => {
    if (!topStreamTimes) return;
    if (!isMoreResults) return;

    try {
      setIsFetchingMore(true);
      const { data, error } = await supabase
        .from<StreamHistoryWithTotalViewCount>(view)
        .select(columns)
        .lt('time', topStreamTimes[topStreamTimes.length - 1].time)
        .order('time', { ascending: false })
        .limit(resultLimit);

      if (error) {
        console.error(error);
        toast.error(error.message);
        return;
      }

      if (data) {
        setTopStreamTimes((state) => [...state!, ...data]);
        if (data.length < resultLimit) {
          setIsMoreResults(false);
        }
      }
    } catch (error: any) {
      console.error(error);
      toast.error(error?.message || 'Unknown error');
    } finally {
      setIsFetchingMore(false);
    }
  };

  return topStreamTimes ? (
    <>
      <TopStreamTimesList history={topStreamTimes} />
      {isMoreResults ? (
        <Button
          className='min-w-[10rem]'
          disabled={isFetchingMore}
          onClick={fetchMoreTimes}
          type='button'
        >
          Load more results{' '}
          {isFetchingMore && <Spinner className='ml-2 h-4 w-4' />}
        </Button>
      ) : (
        <MutedText className='font-medium'>End of results</MutedText>
      )}
    </>
  ) : null;
};

export default TopStreams;
