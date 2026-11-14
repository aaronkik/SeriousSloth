import type { ReactNode } from 'react';
import { Card } from '~/components/ui/card';
import { Empty, EmptyDescription, EmptyHeader } from '~/components/ui/empty';
import type { EmoteListEntry, EmotesByDate } from '~/lib/api/emotes-service';
import { formatDate } from '~/lib/helpers';

export const EmoteGrid = ({ emotes }: { emotes: EmoteListEntry[] }) => (
  <ul className='flex flex-row flex-wrap justify-start gap-4'>
    {emotes.map((emote) => (
      <li key={emote.id}>
        <Card size='sm' className='size-40 items-center justify-center px-4'>
          <div className='relative size-full'>
            <img
              alt={`${emote.name} emote`}
              src={emote.emoteUrl}
              loading='lazy'
              className='absolute inset-0 size-full object-contain'
            />
          </div>
          <p className='tracking-wide break-all text-center'>{emote.name}</p>
        </Card>
      </li>
    ))}
  </ul>
);

export const EmoteDateGroups = ({ emotes }: { emotes: EmotesByDate }) =>
  Object.entries(emotes).map(([date, dateEmotes]) => (
    <div key={date} className='flex flex-col'>
      <div className='py-6 flex items-center gap-2'>
        <p className='text-lg font-semibold'>{formatDate(date, 'LL')}</p>
        <hr className='w-48 h-1 bg-primary rounded-full' />
      </div>
      <EmoteGrid emotes={dateEmotes} />
    </div>
  ));

type Props = {
  count: number;
  emptyMessage?: string;
  children: ReactNode;
};

const EmotesList = ({
  count,
  emptyMessage = 'No emotes to display',
  children,
}: Props) => {
  if (count === 0) {
    return (
      <Empty className='min-h-100 border-none'>
        <EmptyHeader>
          <EmptyDescription>{emptyMessage}</EmptyDescription>
        </EmptyHeader>
      </Empty>
    );
  }

  return (
    <div className='flex flex-col flex-wrap justify-center py-6'>
      {children}
    </div>
  );
};

export default EmotesList;
