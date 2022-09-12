import { ComponentPropsWithRef, forwardRef } from 'react';
import { twMerge } from 'tailwind-merge';

const Link = forwardRef<HTMLAnchorElement, ComponentPropsWithRef<'a'>>(
  ({ className, ...props }, ref) => (
    <a
      {...props}
      className={twMerge(
        'inline-flex items-center text-neutral-100 no-underline transition-all duration-150 hover:text-purple-400 hover:underline hover:underline-offset-4',
        className
      )}
      ref={ref}
    />
  )
);

Link.displayName = 'Link';

export default Link;
