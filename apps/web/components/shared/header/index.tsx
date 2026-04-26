import Link from 'next/link';
import { Container, SlothLogo } from '~/components/shared';
import HeaderLink from './header-link';

const Header = () => (
  <header>
    <Container className='flex flex-row items-center gap-4 py-4'>
      <Link href='/' passHref>
        <SlothLogo className='rounded-full' width={40} height={40} />
      </Link>
      <nav>
        <HeaderLink href='/'>Home</HeaderLink>
      </nav>
    </Container>
  </header>
);

export default Header;
