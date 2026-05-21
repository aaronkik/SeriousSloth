import { ExternalLink } from 'lucide-react';
import { ComponentPropsWithoutRef } from 'react';
import { cn } from '~/lib/utils';

type Props = {
  className?: ComponentPropsWithoutRef<'a'>['className'];
  loginName: string;
};

const TwitchLink = ({ className, loginName }: Props) => (
  <a
    className={cn(
      'inline-flex items-center gap-1 break-all font-medium text-primary no-underline transition-all duration-75 hover:underline hover:underline-offset-4',
      className
    )}
    href={`https://twitch.tv/${encodeURIComponent(loginName)}`}
    rel='noopener noreferrer'
    target='_blank'
  >
    {`twitch.tv/${loginName}`}
    <ExternalLink className='size-4' />
  </a>
);

export default TwitchLink;
