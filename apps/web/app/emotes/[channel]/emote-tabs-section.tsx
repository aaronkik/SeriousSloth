import EmoteTabs from '~/app/emotes/components/emote-tabs';
import { getEmoteData } from '~/app/emotes/[channel]/queries';

const EmoteTabsSection = async ({
  channelParam,
}: {
  channelParam: string;
}) => {
  const { activeEmotes, removedEmotes } = await getEmoteData(channelParam);

  return (
    <EmoteTabs activeEmotes={activeEmotes} removedEmotes={removedEmotes} />
  );
};

export default EmoteTabsSection;
