import { InferGetStaticPropsType } from 'next';
import Image from 'next/image';
import { Card } from '~/components/shared';
import { getStaticProps } from '~/pages/global-emotes';

type Props = {
  globalEmotes: InferGetStaticPropsType<typeof getStaticProps>['globalEmotes'];
};

const GlobalEmotesList = ({ globalEmotes }: Props) => (
  <ul
    data-testid='globalEmoteList'
    className='flex flex-row flex-wrap justify-center gap-6 py-8'
  >
    {globalEmotes.map((emote, index) => (
      <li key={emote.id}>
        <Card className='h-36 w-36 p-4'>
          <div className='flex h-full w-full flex-col items-center gap-2'>
            <div className='relative h-full w-full'>
              <Image
                alt={`${emote.name} emote`}
                blurDataURL={emote.blurDataUrl}
                data-testid={`emoteImage${index}`}
                fill
                placeholder='blur'
                src={emote.largeImageUrl}
                style={{ objectFit: 'contain' }}
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

export default GlobalEmotesList;
