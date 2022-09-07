import { DetailedHTMLProps, HTMLAttributes } from 'react';
import { twMerge } from 'tailwind-merge';

type Props = DetailedHTMLProps<HTMLAttributes<HTMLDivElement>, HTMLDivElement>;

const Container = ({ className, ...props }: Props) => (
  <div className={twMerge('mx-auto max-w-6xl', className)} {...props} />
);

export default Container;
