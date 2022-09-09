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
    <ul className='grid grid-cols-1 sm:grid-cols-2 gap-6 grid-rows-1'>
      {navigationRoutes.map(({ path, title, description }) => (
        <li key={path}>
          <Link href={path} passHref>
            <a>
              <Card className='p-4 flex flex-col items-center'>
                <Heading
                  className='text-xl md:text-2xl text-purple-400'
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
