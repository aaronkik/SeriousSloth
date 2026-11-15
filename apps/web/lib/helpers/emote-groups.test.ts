import { describe, expect, it } from 'vitest';
import type { Emote } from '~/lib/api/emotes-service';
import { groupEmotesByDate, sortEmotesByName } from './emote-groups';

const makeEmote = (name: string): Emote => ({
  id: `id-${name}`,
  name,
  images: { url_1x: '', url_2x: '', url_4x: '' },
  format: ['static'],
  scale: ['1.0'],
  theme_mode: ['dark'],
});

const names = (emotes: { name: string }[]) => emotes.map(({ name }) => name);

describe('sortEmotesByName', () => {
  it('orders emotes alphabetically', () => {
    const emotes = ['zeta', 'Alpha', 'beta'].map(makeEmote);

    expect(names(sortEmotesByName(emotes))).toEqual(['Alpha', 'beta', 'zeta']);
  });
});

describe('groupEmotesByDate', () => {
  const added = [
    { emote: makeEmote('older'), addedAt: '2026-01-01T23:59:59Z' },
    { emote: makeEmote('newerB'), addedAt: '2026-03-10T08:00:00Z' },
    { emote: makeEmote('newerA'), addedAt: '2026-03-10T20:00:00Z' },
  ];

  it('groups by UTC day with the most recent day first', () => {
    const groups = groupEmotesByDate(added, 'addedAt');

    expect(Object.keys(groups)).toEqual(['2026-03-10', '2026-01-01']);
  });

  it('orders emotes alphabetically within a day', () => {
    const groups = groupEmotesByDate(added, 'addedAt');

    expect(names(groups['2026-03-10'])).toEqual(['newerA', 'newerB']);
  });

  it('lists an emote once per day when re-added on the same day', () => {
    const reAdded = [
      ...added,
      { emote: makeEmote('newerA'), addedAt: '2026-03-10T22:00:00Z' },
    ];

    const groups = groupEmotesByDate(reAdded, 'addedAt');

    expect(names(groups['2026-03-10'])).toEqual(['newerA', 'newerB']);
  });

  it('keeps re-adds on different days', () => {
    const reAdded = [
      ...added,
      { emote: makeEmote('older'), addedAt: '2026-03-10T22:00:00Z' },
    ];

    const groups = groupEmotesByDate(reAdded, 'addedAt');

    expect(names(groups['2026-01-01'])).toEqual(['older']);
    expect(names(groups['2026-03-10'])).toEqual(['newerA', 'newerB', 'older']);
  });
});
