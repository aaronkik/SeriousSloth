import { render, RenderOptions } from '@testing-library/react';
import { ReactNode } from 'react';
import { ToastContainer } from 'react-toastify';

const Wrapper = ({ children }: { children: ReactNode }) => (
  <>
    {children}
    <ToastContainer />
  </>
);

const customRender = (ui: React.ReactElement, options?: RenderOptions) =>
  render(ui, {
    wrapper: Wrapper,
    ...options,
  });

export * from '@testing-library/react';

export { customRender as render };
