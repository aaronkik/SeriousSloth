import Link from 'next/link';
import { ReactNode } from 'react';

type Props = {
  href: string;
  children: ReactNode;
};

const HeaderLink = ({ href, children }: Props) => (
  <Link
    href={href}
    className='rounded-md bg-card px-4 py-2 text-lg font-medium transition-colors duration-150 hover:text-primary'
  >
    {children}
  </Link>
);

export default HeaderLink;
