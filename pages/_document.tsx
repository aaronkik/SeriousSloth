import { Html, Head, Main, NextScript } from 'next/document';

const Document = () => (
  <Html lang='en'>
    <Head>
      <link
        href='https://fonts.googleapis.com/css2?family=Roboto:wght@400;500;700&display=swap'
        rel='stylesheet'
      />
      <link rel='icon' href='/assets/favicon.ico' />
      <link
        rel='icon'
        type='image/png'
        sizes='32x32'
        href='/assets/favicon-32x32.png'
      />
      <link
        rel='icon'
        type='image/png'
        sizes='16x16'
        href='/assets/favicon-16x16.png'
      />
    </Head>
    <body className='bg-neutral-900 text-neutral-100'>
      <Main />
      <NextScript />
    </body>
  </Html>
);

export default Document;
