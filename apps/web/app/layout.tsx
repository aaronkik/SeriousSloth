import { Roboto } from 'next/font/google';
import { Container, Footer, Header } from '~/components/shared';
import './globals.css';
import React from 'react';
import type { Metadata } from 'next';

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
        <link rel='icon' href='/favicon.ico' sizes='any' />
        <link rel='icon' href='/icon.png' type='image/png' sizes='32x32' />
      </head>
      <body className='bg-background text-foreground flex h-screen flex-col'>
        <Header />
        <main className='flex-1 py-8'>
          <Container>{children}</Container>
        </main>
        <Footer />
      </body>
    </html>
  );
}
