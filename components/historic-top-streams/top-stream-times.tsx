import { useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import { supabase } from '~/lib/supabase';
import { StreamHistoryWithTotalViewCount } from '~/types/supabase/overrides';
import TopStreamTimesList from './top-stream-times-list';

const TopStreams = () => {
  const [topStreamTimes, setTopStreamTimes] = useState<
    Array<StreamHistoryWithTotalViewCount> | undefined
  >(undefined);

  useEffect(() => {
    supabase
      .from<StreamHistoryWithTotalViewCount>(
        'stream_history_timestamp_with_total_view_count'
      )
      .select('id,time,total_viewer_count,total_streams')
      .order('time', { ascending: false })
      .limit(50)
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
        }
      });
  }, []);

  return topStreamTimes ? (
    <TopStreamTimesList history={topStreamTimes} />
  ) : null;
};

export default TopStreams;
