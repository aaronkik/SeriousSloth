import { DetailedHTMLProps, HTMLAttributes } from 'react';
import { twMerge } from 'tailwind-merge';

type Props = DetailedHTMLProps<HTMLAttributes<HTMLDivElement>, HTMLDivElement>;

const Card = ({ className, ...props }: Props) => (
  <div className={twMerge('rounded-md bg-neutral-800', className)} {...props} />
);

export default Card;
