import 'react-toastify/dist/ReactToastify.css';
import '~/styles/globals.css';
import type { AppProps } from 'next/app';
import { ToastContainer } from 'react-toastify';
import { PageLayout } from '~/components/shared';

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <>
      <PageLayout>
        <Component {...pageProps} />
      </PageLayout>
      <ToastContainer autoClose={3000} position='top-right' limit={3} />
    </>
  );
}

export default MyApp;
