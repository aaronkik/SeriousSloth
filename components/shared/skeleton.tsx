import { DetailedHTMLProps, HTMLAttributes } from 'react';
import { twMerge } from 'tailwind-merge';

type Props = DetailedHTMLProps<HTMLAttributes<HTMLDivElement>, HTMLDivElement>;

const Skeleton = ({ className, ...props }: Props) => (
  <div
    className={twMerge('animate-pulse rounded-md bg-neutral-500', className)}
    {...props}
  />
);

export default Skeleton;
