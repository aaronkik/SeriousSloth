import DynamicLastUpdated from '~/app/emotes/components/dynamic-last-updated';
import { getEmoteData } from '~/app/emotes/[channel]/queries';

const LastUpdatedSection = async ({
  channelParam,
}: {
  channelParam: string;
}) => {
  const { updatedAt } = await getEmoteData(channelParam);

  return <DynamicLastUpdated lastUpdated={updatedAt} />;
};

export default LastUpdatedSection;
