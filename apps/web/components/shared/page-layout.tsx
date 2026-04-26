import { ReactNode } from 'react';
import { Container, Footer, Header } from '~/components/shared';

type Props = {
  children: ReactNode;
};

const PageLayout = ({ children }: Props) => (
  <>
    <div className='flex h-screen flex-col'>
      <Header />
      <main className='flex-1 py-8'>
        <Container>{children}</Container>
      </main>
      <Footer />
    </div>
  </>
);

export default PageLayout;
