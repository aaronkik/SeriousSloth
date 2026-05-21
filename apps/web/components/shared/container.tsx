import { DetailedHTMLProps, HTMLAttributes } from 'react';
import { cn } from '~/lib/utils';

type Props = DetailedHTMLProps<HTMLAttributes<HTMLDivElement>, HTMLDivElement>;

const Container = ({ className, ...props }: Props) => (
  <div className={cn('mx-auto max-w-6xl px-4', className)} {...props} />
);

export default Container;
