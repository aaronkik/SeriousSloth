import { FaGithub } from 'react-icons/fa';
import { Link } from '~/components/shared';

const SourceCodeLink = () => (
  <Link
    className='font-medium text-neutral-100 hover:text-purple-500'
    href='https://github.com/aaronkik/SeriousSloth'
    rel='noopener noreferrer'
    target='_blank'
  >
    Source code
    <FaGithub className='ml-2 h-4 w-4' />
  </Link>
);

export default SourceCodeLink;
