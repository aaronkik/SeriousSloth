import EmoteTabs from '~/app/emotes/components/emote-tabs';
import { getEmoteData } from '~/app/emotes/[channel]/queries';

const EmoteTabsSection = async ({ channelParam }: { channelParam: string }) => {
  const emoteData = await getEmoteData(channelParam);

  return <EmoteTabs {...emoteData} />;
};

export default EmoteTabsSection;
