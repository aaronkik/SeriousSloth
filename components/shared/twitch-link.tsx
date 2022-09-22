import { ExternalLinkIcon } from '@heroicons/react/solid';
import { ComponentPropsWithoutRef } from 'react';
import { twMerge } from 'tailwind-merge';

type Props = {
  className?: ComponentPropsWithoutRef<'a'>['className'];
  loginName: string;
};

const TwitchLink = ({ className, loginName }: Props) => (
  <a
    className={twMerge(
      'inline-flex items-center break-all font-medium text-purple-400 no-underline transition-all duration-75 hover:underline hover:underline-offset-4',
      className
    )}
    href={`https://twitch.tv/${encodeURIComponent(loginName)}`}
    rel='noopener noreferrer'
    target='_blank'
  >
    {`twitch.tv/${loginName}`}
    <div>
      <ExternalLinkIcon className='ml-1 h-4 w-4' />
    </div>
  </a>
);

export default TwitchLink;
