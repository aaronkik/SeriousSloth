import { formatEmoteCDNUrl, EmoteCDNValues } from './helpers';

describe('formatEmoteCDNUrl', () => {
  it('replaces template holders with emote values', () => {
    const templateUrl =
      'https://static-cdn.jtvnw.net/emoticons/v2/{{id}}/{{format}}/{{theme_mode}}/{{scale}}';

    const emoteValues = {
      id: '123',
      format: 'default',
      theme_mode: 'light',
      scale: '1.0',
    } as EmoteCDNValues;

    const { id, format, theme_mode, scale } = emoteValues;
    const expectedUrl = `https://static-cdn.jtvnw.net/emoticons/v2/${id}/${format}/${theme_mode}/${scale}`;

    expect(formatEmoteCDNUrl(templateUrl, emoteValues)).toBe(expectedUrl);
  });
});
