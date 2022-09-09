import dynamic from 'next/dynamic';
import { ComponentProps, Suspense } from 'react';
import { Skeleton } from '~/components/shared';

/**
 * Component used to address hydration errors
 */

const LastUpdated = dynamic(() => import('./last-updated'), {
  suspense: true,
  ssr: false,
});

const DynamicLastUpdated = ({
  lastUpdated,
}: ComponentProps<typeof LastUpdated>) => (
  <Suspense fallback={<Skeleton className='w-60 h-5' />}>
    <LastUpdated lastUpdated={lastUpdated} />
  </Suspense>
);

export default DynamicLastUpdated;
