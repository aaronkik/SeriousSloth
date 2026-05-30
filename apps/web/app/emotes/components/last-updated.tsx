import { MutedText } from '~/components/shared';
import { timeFromNow } from '~/lib/helpers';

type Props = {
  lastUpdated: number;
};

const LastUpdated = ({ lastUpdated }: Props) => (
  <MutedText className='text-sm'>
    Last updated: {timeFromNow(lastUpdated)}
  </MutedText>
);

export default LastUpdated;
