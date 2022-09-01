export type EmoteCDNValues = {
  id: string;
  format: 'default' | 'static' | 'animated';
  theme_mode: 'light' | 'dark';
  scale: '1.0' | '2.0' | '3.0';
};

/**
 * @see https://dev.twitch.tv/docs/irc/emotes#cdn-template
 */
export const formatEmoteCDNUrl = (
  templateUrl: string,
  emoteValues: EmoteCDNValues
) => {
  const { id, format, theme_mode, scale } = emoteValues;
  return templateUrl
    .replace('{{id}}', id)
    .replace('{{format}}', format)
    .replace('{{theme_mode}}', theme_mode)
    .replace('{{scale}}', scale);
};
