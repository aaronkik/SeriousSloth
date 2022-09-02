import { fetchClientCredentials, fetchGlobalEmotes } from './fetchers';

describe('fetchClientCredentials', () => {
  it('Payload includes expected values', async () => {
    const payload = await fetchClientCredentials();
    const expectedPayload = {
      access_token: expect.any(String),
      expires_in: expect.any(Number),
      token_type: expect.any(String),
    };
    expect(payload).toMatchObject(expectedPayload);
  });
});

describe('fetchGlobalEmotes', () => {
  it('Payload includes expected values', async () => {
    const { data, template } = await fetchGlobalEmotes('accessToken');

    expect(template).toEqual(expect.any(String));

    data.forEach((emote) => {
      const { id, name, images, format, scale, theme_mode } = emote;
      expect(id).toEqual(expect.any(String));
      expect(name).toEqual(expect.any(String));
      expect(images).toMatchObject({
        url_1x: expect.any(String),
        url_2x: expect.any(String),
        url_4x: expect.any(String),
      });
      expect(format).toEqual(
        expect.arrayContaining([expect.stringMatching(/^(static|animated)$/)])
      );
      expect(scale).toEqual(
        expect.arrayContaining([expect.stringMatching(/^(1\.0|2\.0|3\.0)$/)])
      );
      expect(theme_mode).toEqual(
        expect.arrayContaining([expect.stringMatching(/^(light|dark)$/)])
      );
    });
  });
});
