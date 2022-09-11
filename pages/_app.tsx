import 'react-toastify/dist/ReactToastify.css';
import '~/styles/globals.css';
import type { AppProps } from 'next/app';
import { ToastContainer } from 'react-toastify';
import { Container } from '~/components/shared';

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <>
      <Container>
        <Component {...pageProps} />
      </Container>
      <ToastContainer autoClose={3000} position='top-right' limit={3} />
    </>
  );
}

export default MyApp;
