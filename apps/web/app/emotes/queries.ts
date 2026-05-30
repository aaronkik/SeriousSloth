import { cacheLife } from 'next/cache';
import { getChannels, type Channel } from '~/lib/api/emotes-service';

export const getCachedChannels = async (): Promise<Channel[]> => {
  'use cache';
  cacheLife({ revalidate: 300 });

  return getChannels();
};
