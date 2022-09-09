import { ComponentProps } from 'react';
import { twMerge } from 'tailwind-merge';

const Button = ({ className, ...props }: ComponentProps<'button'>) => (
  <button
    className={twMerge(
      'inline-flex select-none items-center justify-center rounded-md bg-purple-500 px-4 py-1.5 font-medium leading-4 text-neutral-100 transition-colors duration-150 focus-visible:bg-purple-700 active:bg-purple-600',
      className
    )}
    {...props}
  />
);

export default Button;
