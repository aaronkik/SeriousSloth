import { ComponentPropsWithRef, forwardRef } from 'react';
import { cn } from '~/lib/utils';

const Link = forwardRef<HTMLAnchorElement, ComponentPropsWithRef<'a'>>(
  ({ className, ...props }, ref) => (
    <a
      {...props}
      className={cn(
        'inline-flex items-center text-foreground no-underline transition-all duration-150 hover:text-primary hover:underline hover:underline-offset-4',
        className
      )}
      ref={ref}
    />
  )
);

Link.displayName = 'Link';

export default Link;
