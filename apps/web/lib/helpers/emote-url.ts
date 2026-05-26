import type { Emote } from '~/lib/api/emotes-service';

const CDN_BASE = 'https://static-cdn.jtvnw.net/emoticons/v2';

const THEME_MODE_PRECEDENCE = ['dark', 'light'] as const;
const FORMAT_PRECEDENCE = ['animated', 'static'] as const;
const SCALE_PRECEDENCE = ['3.0', '2.0', '1.0'] as const;

const pickByPrecedence = <T extends string>(
  availableValues: readonly string[],
  orderedPrecedence: readonly T[]
): T => {
  const preferredValue = orderedPrecedence.find((value) =>
    availableValues.includes(value)
  );

  if (preferredValue) {
    return preferredValue;
  }

  return (availableValues[0] ??
    orderedPrecedence[orderedPrecedence.length - 1]) as T;
};

export const buildEmoteUrl = (emote: Emote): string => {
  const selectedThemeMode = pickByPrecedence(
    emote.theme_mode,
    THEME_MODE_PRECEDENCE
  );
  const selectedFormat = pickByPrecedence(emote.format, FORMAT_PRECEDENCE);
  const selectedScale = pickByPrecedence(emote.scale, SCALE_PRECEDENCE);

  return `${CDN_BASE}/${emote.id}/${selectedFormat}/${selectedThemeMode}/${selectedScale}`;
};
