import dynamic from 'next/dynamic';
import { Skeleton } from '~/components/shared';

const DynamicLastUpdated = dynamic(() => import('./last-updated'), {
  ssr: false,
  loading: () => <Skeleton className='h-5 w-60' />,
});

export default DynamicLastUpdated;
