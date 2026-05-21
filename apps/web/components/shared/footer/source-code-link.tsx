import { ExternalLink } from 'lucide-react';
import { Link } from '~/components/shared';

const SourceCodeLink = () => (
  <Link
    className='gap-2 font-medium'
    href='https://github.com/aaronkik/SeriousSloth'
    rel='noopener noreferrer'
    target='_blank'
  >
    Source code
    <ExternalLink className='size-4' />
  </Link>
);

export default SourceCodeLink;
