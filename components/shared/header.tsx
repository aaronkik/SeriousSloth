import Link from 'next/link';
import { ReactNode } from 'react';
import SlothLogo from '~/components/shared/sloth-logo';

type HeaderLinkProps = {
  href: string;
  children: ReactNode;
};
const HeaderLink = ({ href, children }: HeaderLinkProps) => (
  <Link href={href}>
    <a className='rounded-md bg-neutral-800 px-4 py-2 text-lg font-medium transition-colors duration-150 hover:text-purple-400'>
      {children}
    </a>
  </Link>
);

const Header = () => (
  <header className='flex flex-row items-center gap-4 py-4'>
    <Link href='/' passHref>
      <a>
        <SlothLogo className='rounded-full' width={40} height={40} />
      </a>
    </Link>
    <nav>
      <HeaderLink href='/'>Home</HeaderLink>
    </nav>
  </header>
);

export default Header;
