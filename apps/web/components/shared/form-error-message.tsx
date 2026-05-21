import { ComponentProps } from 'react';
import { cn } from '~/lib/utils';

const FormErrorMessage = ({ className, ...props }: ComponentProps<'p'>) => (
  <p
    className={cn('text-sm text-destructive', className)}
    role='alert'
    {...props}
  />
);

export default FormErrorMessage;
