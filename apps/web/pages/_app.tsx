import '~/styles/globals.css';
import type { AppProps } from 'next/app';
import dynamic from 'next/dynamic';
import { PageLayout } from '~/components/shared';

const Toaster = dynamic(
  () => import('~/components/ui/sonner').then((m) => m.Toaster),
  { ssr: false }
);

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <>
      <PageLayout>
        <Component {...pageProps} />
      </PageLayout>
      <Toaster position='top-right' />
    </>
  );
}

export default MyApp;
