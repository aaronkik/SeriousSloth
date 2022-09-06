import { validUsername } from '~/__mocks__/data/twitch/constants';
import {
  fetchClientCredentials,
  fetchGlobalEmotes,
  fetchUsers,
} from './fetchers';

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

describe('fetchUsers', () => {
  it('Returns users from API', async () => {
    const { data } = await fetchUsers('accessToken', validUsername);

    data.forEach((user) => {
      expect(user).toMatchObject({
        broadcaster_type: expect.stringMatching(/^partner$|^affiliate$|^$/),
        created_at: expect.any(String),
        description: expect.any(String),
        display_name: expect.any(String),
        id: expect.any(String),
        login: expect.any(String),
        offline_image_url: expect.any(String),
        profile_image_url: expect.any(String),
        type: expect.stringMatching(/^staff$|^admin$|^global_mod$|^$/),
      });
    });
  });
});
