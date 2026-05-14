import { Card } from '~/components/shared';
import type { ActiveEmote } from '~/lib/api/emotes-service';

type Props = {
  emotes: ActiveEmote[];
};

const EmotesList = ({ emotes }: Props) => (
  <ul
    data-testid='emoteList'
    className='flex flex-row flex-wrap justify-center gap-6 py-8'
  >
    {emotes.map(({ emote }, index) => (
      <li key={emote.id}>
        <Card className='h-36 w-36 p-4'>
          <div className='flex h-full w-full flex-col items-center gap-2'>
            <div className='relative h-full w-full'>
              <img
                alt={`${emote.name} emote`}
                data-testid={`emoteImage${index}`}
                src={emote.images.url_4x}
                className='absolute inset-0 h-full w-full object-contain'
              />
            </div>
            <p className='tracking-wide' data-testid={`emoteName${index}`}>
              {emote.name}
            </p>
          </div>
        </Card>
      </li>
    ))}
  </ul>
);

export default EmotesList;
