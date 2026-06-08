import { Card } from '~/components/ui/card';
import { Empty, EmptyDescription, EmptyHeader } from '~/components/ui/empty';
import type {
  ActiveEmoteEntry,
  RemovedEmoteEntry,
} from '~/lib/api/emotes-service';
import localizedFormat from 'dayjs/plugin/localizedFormat';
import dayjs from 'dayjs';

dayjs.extend(localizedFormat);

type Props = {
  count: number;
  emotes:
    | Record<string, ActiveEmoteEntry[]>
    | Record<string, RemovedEmoteEntry[]>;
  emptyMessage?: string;
};

const EmotesList = ({
  count,
  emotes,
  emptyMessage = 'No emotes to display',
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
      {Object.entries(emotes).map(
        ([groupingDate, emotes]: [
          string,
          (ActiveEmoteEntry | RemovedEmoteEntry)[],
        ]) => (
          <div key={groupingDate} className='flex flex-col'>
            <div className='py-6 flex items-center gap-2'>
              <p className='text-lg font-semibold'>
                {dayjs(groupingDate).format('LL')}
              </p>
              <hr className='w-48 h-1 bg-primary rounded-full' />
            </div>
            <ul className='flex flex-row flex-wrap justify-start gap-4'>
              {emotes.map((emote) => (
                <li key={emote.id}>
                  <Card
                    size='sm'
                    className='size-40 items-center justify-center px-4'
                  >
                    <div className='relative size-full'>
                      <img
                        alt={`${emote.name} emote`}
                        src={emote.emoteUrl}
                        loading='lazy'
                        className='absolute inset-0 size-full object-contain'
                      />
                    </div>
                    <p className='tracking-wide break-all text-center'>
                      {emote.name}
                    </p>
                  </Card>
                </li>
              ))}
            </ul>
          </div>
        ),
      )}
    </div>
  );
};

export default EmotesList;
