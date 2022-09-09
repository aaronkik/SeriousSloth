import { ComponentProps, forwardRef } from 'react';
import { twMerge } from 'tailwind-merge';

const Input = forwardRef<HTMLInputElement, ComponentProps<'input'>>(
  ({ className, ...props }, ref) => (
    <input
      className={twMerge(
        'w-full rounded-md border border-transparent bg-neutral-900 placeholder:text-neutral-500 focus:border-purple-500 focus:bg-transparent focus:ring-purple-500',
        className
      )}
      ref={ref}
      {...props}
    />
  )
);

Input.displayName = 'Input';

export default Input;
