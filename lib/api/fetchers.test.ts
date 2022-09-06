import {
  invalidUsernames,
  longUsername,
  shortUsername,
  validUsername,
} from '~/__mocks__/data/twitch/constants';
import { fetchUsers } from './fetchers';

describe('fetchUsers', () => {
  it('Returns error with an empty username', async () => {
    await expect(fetchUsers({ username: '' })).rejects.toMatchObject({
      message: expect.any(String),
      status: 400,
    });
  });

  it('Returns error with a short username', async () => {
    await expect(fetchUsers({ username: shortUsername })).rejects.toMatchObject(
      {
        message: expect.any(String),
        status: 400,
      }
    );
  });

  it('Returns error with a long username', async () => {
    await expect(fetchUsers({ username: longUsername })).rejects.toMatchObject({
      message: expect.any(String),
      status: 400,
    });
  });

  it('Returns error with an invalid username', async () => {
    for (const username of invalidUsernames) {
      await expect(fetchUsers({ username: username })).rejects.toMatchObject({
        message: expect.any(String),
        status: 400,
      });
    }
  });

  it('Returns array of users with a valid username', async () => {
    const { data } = await fetchUsers({ username: validUsername });

    expect(data[0]).toMatchObject({
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
