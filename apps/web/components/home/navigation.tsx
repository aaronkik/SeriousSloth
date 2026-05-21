import Link from 'next/link';
import { Card, Heading } from '../shared';

const navigationRoutes = [
  {
    path: '/emotes',
    title: 'Emotes',
    description: 'Browse Twitch emotes by channel',
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
          <Link href={path}>
            <Card className='items-center p-4 transition-all duration-150 hover:shadow-md hover:shadow-primary/10'>
              <Heading
                className='text-xl text-primary md:text-2xl'
                variant='h2'
              >
                {title}
              </Heading>
              <p>{description}</p>
            </Card>
          </Link>
        </li>
      ))}
    </ul>
  </nav>
);

export default Navigation;
