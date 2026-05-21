import { createElement, DetailedHTMLProps, HTMLAttributes } from 'react';
import { cn } from '~/lib/utils';

type Props = DetailedHTMLProps<
  HTMLAttributes<HTMLHeadingElement>,
  HTMLHeadingElement
> & { variant: 'h1' | 'h2' | 'h3' | 'h4' | 'h5' | 'h6' };

const Heading = ({ className, variant, ...props }: Props) =>
  createElement(variant, {
    className: cn(
      'font-heading text-3xl md:text-5xl tracking-wide font-bold text-foreground',
      className
    ),
    ...props,
  });

export default Heading;
