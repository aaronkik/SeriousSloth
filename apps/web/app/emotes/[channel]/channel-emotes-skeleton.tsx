import { Skeleton } from '~/components/shared';

const ChannelEmotesSkeleton = () => (
  <>
    <div className='mb-2 flex flex-col items-center gap-2 text-center'>
      <Skeleton className='h-10 w-72 md:h-14' />
      <Skeleton className='h-5 w-60' />
    </div>
    <div className='mt-6 flex flex-col gap-2'>
      <Skeleton className='mx-auto h-9 w-56 rounded-full' />
      <ul className='flex flex-row flex-wrap justify-center gap-6 py-12'>
        {Array.from({ length: 12 }).map((_, index) => (
          <li
            key={index}
            className='flex flex-col items-center gap-2'
          >
            <Skeleton className='size-36 rounded-xl' />
            <Skeleton className='h-4 w-20' />
          </li>
        ))}
      </ul>
    </div>
  </>
);

export default ChannelEmotesSkeleton;
