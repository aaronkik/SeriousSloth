import { Html, Head, Main, NextScript } from 'next/document';

const Document = () => (
  <Html>
    <Head>
      <link
        href='https://fonts.googleapis.com/css2?family=Roboto:wght@400;500;700&display=swap'
        rel='stylesheet'
      />
    </Head>
    <body className='bg-neutral-900 text-neutral-100'>
      <Main />
      <NextScript />
    </body>
  </Html>
);

export default Document;
