export interface GlobalChannel {
  type: 'global';
  id: 'global';
  displayName: string;
  icon: string;
}

export interface TwitchChannel {
  type: 'twitch';
  id: string;
  twitchId: string;
  displayName: string;
  imageUrl: string;
}

export type Channel = GlobalChannel | TwitchChannel;

export const channelSlug = (channel: Channel): string =>
  channel.type === 'global' ? 'global' : channel.twitchId;

export const GLOBAL_CHANNEL: GlobalChannel = {
  type: 'global',
  id: 'global',
  displayName: 'Global',
  icon: '🌐',
};
