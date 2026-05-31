import './globals.css';
import React from 'react';
import newrelic from 'newrelic';
import type { Metadata } from 'next';
import { Roboto } from 'next/font/google';
import Script from 'next/script';
import { Container, Footer, Header } from '~/components/shared';
import { Toaster } from '~/components/ui/sonner';

const roboto = Roboto({
  display: 'swap',
  preload: true,
});

export const metadata: Metadata = {
  title: 'SeriousSloth',
  description: 'A web app that interacts with the Twitch API',
};

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  // @ts-ignore
  if (newrelic.agent.collector.isConnected() === false) {
    await new Promise((resolve) => {
      // @ts-ignore
      newrelic.agent.on('connected', resolve);
    });
  }

  const browserTimingHeader = newrelic.getBrowserTimingHeader({
    hasToRemoveScriptWrapper: true,
    allowTransactionlessInjection: true,
  });

  return (
    <html lang='en' className={roboto.className}>
      <head>
        <Script
          id='nr-browser-agent'
          strategy='beforeInteractive'
          dangerouslySetInnerHTML={{ __html: browserTimingHeader }}
        />
        <link rel='icon' href='/favicon.ico' sizes='any' />
        <link rel='icon' href='/icon.png' type='image/png' sizes='32x32' />
      </head>
      <body className='bg-background text-foreground flex h-screen flex-col'>
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
