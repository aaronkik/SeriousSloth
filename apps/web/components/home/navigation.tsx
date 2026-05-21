import Link from 'next/link';
import { Card } from '~/components/ui/card';

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
    <ul className='grid grid-cols-1 gap-4 sm:grid-cols-2'>
      {navigationRoutes.map(({ path, title, description }) => (
        <li key={path}>
          <Link href={path} className='group block'>
            <Card className='flex flex-col items-center gap-3 p-5 text-center ring-1 ring-transparent transition-colors group-hover:ring-primary/60'>
              <p className='text-lg font-semibold text-foreground transition-colors group-hover:text-primary'>
                {title}
              </p>
              <p className='text-sm text-muted-foreground'>{description}</p>
            </Card>
          </Link>
        </li>
      ))}
    </ul>
  </nav>
);

export default Navigation;
