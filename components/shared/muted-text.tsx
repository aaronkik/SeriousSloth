import { DetailedHTMLProps, HTMLAttributes } from 'react';
import { twMerge } from 'tailwind-merge';

type Props = DetailedHTMLProps<
  HTMLAttributes<HTMLParagraphElement>,
  HTMLParagraphElement
>;

const MutedText = ({ className, ...props }: Props) => (
  <p className={twMerge('text-neutral-400', className)} {...props} />
);

export default MutedText;
