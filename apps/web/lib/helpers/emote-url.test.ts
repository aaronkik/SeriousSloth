import type { Emote } from '~/lib/api/emotes-service';
import { buildEmoteUrl } from './emote-url';

const CDN_BASE = 'https://static-cdn.jtvnw.net/emoticons/v2';

const makeEmote = (overrides: Partial<Emote> = {}): Emote => ({
  id: 'emotesv2_0ca61620b2ab44eb86a64146078631f9',
  name: 'danCrazy',
  images: {
    url_1x: 'https://example.com/1x',
    url_2x: 'https://example.com/2x',
    url_4x: 'https://example.com/4x',
  },
  format: ['static', 'animated'],
  scale: ['1.0', '2.0', '3.0'],
  theme_mode: ['light', 'dark'],
  ...overrides,
});

describe('buildEmoteUrl', () => {
  it('picks animated, dark, 3.0 when all variants available', () => {
    const emote = makeEmote();

    expect(buildEmoteUrl(emote)).toBe(
      `${CDN_BASE}/${emote.id}/animated/dark/3.0`
    );
  });

  it('falls back to static when animated unavailable', () => {
    const emote = makeEmote({ format: ['static'] });

    expect(buildEmoteUrl(emote)).toBe(
      `${CDN_BASE}/${emote.id}/static/dark/3.0`
    );
  });

  it('falls back to light when dark unavailable', () => {
    const emote = makeEmote({ theme_mode: ['light'] });

    expect(buildEmoteUrl(emote)).toBe(
      `${CDN_BASE}/${emote.id}/animated/light/3.0`
    );
  });

  it('falls back to 1.0 when only 1.0 available', () => {
    const emote = makeEmote({ scale: ['1.0'] });

    expect(buildEmoteUrl(emote)).toBe(
      `${CDN_BASE}/${emote.id}/animated/dark/1.0`
    );
  });

  it('uses exact values when only one of each option present', () => {
    const emote = makeEmote({
      format: ['static'],
      theme_mode: ['light'],
      scale: ['2.0'],
    });

    expect(buildEmoteUrl(emote)).toBe(
      `${CDN_BASE}/${emote.id}/static/light/2.0`
    );
  });

  it('prefers 2.0 over 1.0 when 3.0 absent', () => {
    const emote = makeEmote({ scale: ['1.0', '2.0'] });

    expect(buildEmoteUrl(emote)).toBe(
      `${CDN_BASE}/${emote.id}/animated/dark/2.0`
    );
  });

  it('handles empty arrays by falling back to lowest precedence value', () => {
    const emote = makeEmote({
      format: [] as Emote['format'],
      theme_mode: [] as Emote['theme_mode'],
      scale: [] as Emote['scale'],
    });

    expect(buildEmoteUrl(emote)).toBe(
      `${CDN_BASE}/${emote.id}/static/light/1.0`
    );
  });

  it('matches CDN template structure', () => {
    const emote = makeEmote({ id: 'abc123' });

    expect(buildEmoteUrl(emote)).toMatch(
      /^https:\/\/static-cdn\.jtvnw\.net\/emoticons\/v2\/abc123\/(animated|static)\/(dark|light)\/(1\.0|2\.0|3\.0)$/
    );
  });
});
