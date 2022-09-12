import Link from 'next/link';
import { Card, Heading } from '../shared';

const navigationRoutes = [
  {
    path: '/global-emotes',
    title: 'Global Emotes',
    description: 'See current global emotes on Twitch',
  },
  {
    path: '/user-search',
    title: 'User Search',
    description:
      'Search for users on Twitch, see account creation dates and more',
  },
];

const Navigation = () => (
  <nav>
    <ul className='grid grid-cols-1 grid-rows-1 gap-6 sm:grid-cols-2'>
      {navigationRoutes.map(({ path, title, description }) => (
        <li key={path}>
          <Link href={path} passHref>
            <a>
              <Card className='flex flex-col items-center p-4 transition-all duration-150 hover:shadow-md hover:shadow-purple-500/10'>
                <Heading
                  className='text-xl text-purple-500 md:text-2xl'
                  variant='h2'
                >
                  {title}
                </Heading>
                <p>{description}</p>
              </Card>
            </a>
          </Link>
        </li>
      ))}
    </ul>
  </nav>
);

export default Navigation;
