import { ExternalLinkIcon } from '@heroicons/react/solid';
import Image from 'next/future/image';
import { ReactNode } from 'react';
import { Card, Heading } from '~/components/shared';
import { capitaliseFirstLetter, formatDate, timeFromNow } from '~/lib/helpers';
import { User } from '~/types/twitch';

type CardSectionProps = {
  title: string;
  children: ReactNode;
};

const CardSection = ({ title, children }: CardSectionProps) => (
  <div className='flex flex-col'>
    <p className='text-sm font-medium uppercase tracking-wider text-neutral-400'>
      {title}
    </p>
    {children}
  </div>
);

type Props = {
  user: User;
};

const UserCard = ({ user }: Props) => {
  const {
    broadcaster_type,
    created_at,
    description,
    display_name,
    id,
    login,
    offline_image_url,
    profile_image_url,
    type,
  } = user;

  return (
    <Card className='flex w-full flex-col gap-4 p-4 md:p-6'>
      <div className='flex flex-col items-center gap-2 md:flex-row md:gap-4'>
        <Image
          src={profile_image_url}
          alt='Profile image'
          className='h-20 w-20 rounded-sm md:h-32 md:w-32'
          width={300}
          height={300}
        />
        <div className='text-center md:text-start'>
          <Heading className='text-2xl' variant='h2'>
            {display_name}
          </Heading>
          <a
            className='inline-flex items-center font-medium text-purple-500 no-underline transition-all duration-75 hover:underline hover:underline-offset-4'
            href={`https://twitch.tv/${login}`}
            rel='noopener noreferrer'
            target='_blank'
          >
            {`twitch.tv/${login}`}
            <ExternalLinkIcon className='ml-1 h-4 w-4' />
          </a>
        </div>
      </div>
      <div className='flex flex-col gap-4 md:flex-row md:gap-12'>
        <CardSection title='Id'>
          <p>{id}</p>
        </CardSection>
        <CardSection title='Created'>
          <p>{formatDate(created_at, 'LLL')}</p>
          <p className='text-sm opacity-75'>({timeFromNow(created_at)})</p>
        </CardSection>
        {broadcaster_type && (
          <CardSection title='Broadcaster type'>
            <p>{capitaliseFirstLetter(broadcaster_type)}</p>
          </CardSection>
        )}
        {type && (
          <CardSection title='User type'>
            <p>{capitaliseFirstLetter(type).replace('_', ' ')}</p>
          </CardSection>
        )}
      </div>
      <CardSection title='Channel description'>
        <p>{description || '(empty)'}</p>
      </CardSection>
      {offline_image_url && (
        <CardSection title='Offline image'>
          <div className='aspect-w-16 aspect-h-9'>
            <Image
              src={offline_image_url}
              alt='Offline image'
              className='rounded-sm'
              width={1920}
              height={1080}
            />
          </div>
        </CardSection>
      )}
    </Card>
  );
};

export default UserCard;
