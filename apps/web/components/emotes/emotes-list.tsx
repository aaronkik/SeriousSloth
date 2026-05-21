import { Card } from '~/components/ui/card';
import type { Emote } from '~/lib/api/emotes-service';

type Props = {
  emotes: Array<{ emote: Emote }>;
};

const EmotesList = ({ emotes }: Props) => (
  <ul
    data-testid='emoteList'
    className='flex flex-row flex-wrap justify-center gap-6 py-12'
  >
    {emotes.map(({ emote }, index) => (
      <li key={emote.id}>
        <Card size='sm' className='size-36 items-center justify-center gap-2 p-4'>
          <div className='relative size-full'>
            <img
              alt={`${emote.name} emote`}
              data-testid={`emoteImage${index}`}
              src={emote.images.url_4x}
              className='absolute inset-0 size-full object-contain'
            />
          </div>
          <p className='tracking-wide' data-testid={`emoteName${index}`}>
            {emote.name}
          </p>
        </Card>
      </li>
    ))}
  </ul>
);

export default EmotesList;
