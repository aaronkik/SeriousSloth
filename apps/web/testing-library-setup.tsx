import { render, RenderOptions } from '@testing-library/react';
import { ReactNode } from 'react';
import { Toaster } from '~/components/ui/sonner';

const Wrapper = ({ children }: { children: ReactNode }) => (
  <>
    {children}
    <Toaster />
  </>
);

const customRender = (ui: React.ReactElement, options?: RenderOptions) =>
  render(ui, {
    wrapper: Wrapper,
    ...options,
  });

export * from '@testing-library/react';

export { customRender as render };
