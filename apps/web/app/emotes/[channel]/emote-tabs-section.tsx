import EmoteTabs from '~/app/emotes/components/emote-tabs';
import { getEmoteData } from '~/app/emotes/[channel]/queries';

const EmoteTabsSection = async ({ channelParam }: { channelParam: string }) => {
  const { activeEmotes, activeEmotesCount, removedEmotes, removedEmotesCount } =
    await getEmoteData(channelParam);

  return (
    <EmoteTabs
      activeEmotes={activeEmotes}
      activeEmotesCount={activeEmotesCount}
      removedEmotes={removedEmotes}
      removedEmotesCount={removedEmotesCount}
    />
  );
};

export default EmoteTabsSection;
