import { DetailedHTMLProps, HTMLAttributes } from 'react';
import { cn } from '~/lib/utils';

type Props = DetailedHTMLProps<
  HTMLAttributes<HTMLParagraphElement>,
  HTMLParagraphElement
>;

const MutedText = ({ className, ...props }: Props) => (
  <p className={cn('text-muted-foreground', className)} {...props} />
);

export default MutedText;
