import Link from 'next/link';
import { formatDate, timeFromNow } from '~/lib/helpers';
import { StreamHistoryWithTotalViewCount } from '~/types/supabase/overrides';
import { Card, MutedText } from '../shared';

type Props = {
  history: Array<StreamHistoryWithTotalViewCount>;
};

const TopStreamTimesList = ({ history }: Props) =>
  history.length ? (
    <ul
      className='flex w-full max-w-2xl flex-col justify-center gap-6'
      data-testid='historicTopStreamTimes'
    >
      {history.map(({ id, time, total_streams, total_viewer_count }) => (
        <li key={id}>
          <Link
            href={`/historic-top-streams/${encodeURIComponent(id)}`}
            passHref
          >
            <a>
              <Card className='flex flex-col p-4 transition-all duration-150 hover:shadow-md hover:shadow-purple-500/10'>
                <p>
                  View the Top {total_streams.toLocaleString()} streams{' '}
                  {timeFromNow(time)}{' '}
                  <span className='text-neutral-400'>
                    ({formatDate(time, 'LLL')})
                  </span>
                </p>
                <MutedText className='text-sm'>
                  Total views: {total_viewer_count.toLocaleString()}{' '}
                </MutedText>
              </Card>
            </a>
          </Link>
        </li>
      ))}
    </ul>
  ) : (
    <p className='text-center text-xl font-semibold tracking-wide'>
      No results returned
    </p>
  );

export default TopStreamTimesList;
