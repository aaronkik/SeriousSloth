import { Card, Skeleton } from '~/components/shared';

const SectionSkeleton = () => (
  <div className='flex flex-col gap-2'>
    <Skeleton className='h-3 w-20' />
    <Skeleton className='h-4 w-32' />
  </div>
);

const UserSearchSkeleton = () => (
  <div className='flex flex-col gap-4'>
    <Card className='w-full px-2 py-2 md:w-1/2 md:self-center'>
      <div className='flex items-center gap-2'>
        <Skeleton className='h-9 flex-1' />
        <Skeleton className='h-9 w-20' />
      </div>
    </Card>
    <Card className='flex w-full flex-col gap-4 p-4 md:p-6'>
      <div className='flex flex-col items-center gap-2 md:flex-row md:gap-4'>
        <Skeleton className='h-20 w-20 md:h-32 md:w-32' />
        <div className='flex flex-col items-center gap-2 md:items-start'>
          <Skeleton className='h-7 w-40' />
          <Skeleton className='h-4 w-32' />
        </div>
      </div>
      <div className='flex flex-col gap-4 md:flex-row md:gap-12'>
        <SectionSkeleton />
        <SectionSkeleton />
        <SectionSkeleton />
      </div>
      <div className='flex flex-col gap-2'>
        <Skeleton className='h-3 w-32' />
        <Skeleton className='h-4 w-full' />
        <Skeleton className='h-4 w-3/4' />
      </div>
    </Card>
  </div>
);

export default UserSearchSkeleton;
