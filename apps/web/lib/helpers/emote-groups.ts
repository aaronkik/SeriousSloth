import type {
  Emote,
  EmoteListEntry,
  EmotesByDate,
} from '~/lib/api/emotes-service';
import { buildEmoteUrl } from './emote-url';

const toEntry = (emote: Emote): EmoteListEntry => ({
  id: emote.id,
  name: emote.name,
  emoteUrl: buildEmoteUrl(emote),
});

const toUtcDate = (timestamp: string) =>
  new Date(timestamp).toISOString().split('T')[0];

const uniqueById = (emotes: Emote[]) => [
  ...new Map(emotes.map((emote) => [emote.id, emote])).values(),
];

export const sortEmotesByName = (emotes: Emote[]): EmoteListEntry[] =>
  emotes.map(toEntry).toSorted((a, b) => a.name.localeCompare(b.name));

export const countEmotes = (groups: EmotesByDate) =>
  Object.values(groups).reduce((total, emotes) => total + emotes.length, 0);

export const groupEmotesByDate = <K extends string>(
  items: Array<{ emote: Emote } & Record<K, string>>,
  timestampKey: K,
): EmotesByDate => {
  const groups = Object.groupBy(items, (item) => toUtcDate(item[timestampKey]));

  return Object.fromEntries(
    Object.entries(groups)
      .toSorted(([a], [b]) => b.localeCompare(a))
      .map(([date, group = []]) => [
        date,
        sortEmotesByName(uniqueById(group.map(({ emote }) => emote))),
      ]),
  );
};
