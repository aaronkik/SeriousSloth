import './globals.css';
import type { Metadata } from 'next';
import { Roboto } from 'next/font/google';
import React, { Suspense } from 'react';
import { Container, Footer, Header } from '~/components/shared';
import { Toaster } from '~/components/ui/sonner';
import { NewRelicBrowserTiming } from './newrelic-browser-timing';

const roboto = Roboto({
  display: 'swap',
  preload: true,
});

export const metadata: Metadata = {
  title: 'SeriousSloth',
  description: 'A web app that interacts with the Twitch API',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang='en' className={roboto.className}>
      <head>
        <Suspense fallback={null}>
          <NewRelicBrowserTiming />
        </Suspense>
        <link rel='icon' href='/favicon.ico' sizes='any' />
        <link rel='icon' href='/icon.png' type='image/png' sizes='32x32' />
      </head>
      <body className='flex h-screen flex-col bg-background text-foreground'>
        <Header />
        <main className='flex-1 py-8'>
          <Container>{children}</Container>
        </main>
        <Footer />
        <Toaster position='top-right' />
      </body>
    </html>
  );
}
