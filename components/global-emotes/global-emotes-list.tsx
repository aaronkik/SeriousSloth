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
    className='flex flex-row gap-6 justify-center flex-wrap py-8'
  >
    {globalEmotes.map((emote, index) => (
      <li key={emote.id}>
        <Card className='w-36 h-36 p-4'>
          <div className='flex flex-col gap-2 w-full h-full items-center'>
            <div className='relative w-full h-full'>
              <Image
                alt={`${emote.name} emote`}
                src={emote.image}
                layout='fill'
                objectFit='contain'
                data-testid={`emoteImage${index}`}
                placeholder='blur'
                blurDataURL={emote.image}
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
