import { http, HttpResponse } from 'msw';
import { server } from '~/__mocks__/msw-server';
import { get } from './fetch';

const url = 'https://emotes-service.test/v1/emotes/unknown';

describe('get', () => {
  it('rejects with the status, url and body of a failed response', async () => {
    server.use(
      http.get(url, () =>
        HttpResponse.json(
          { message: 'Channel not found', data: { channel: 'unknown' } },
          { status: 404 },
        ),
      ),
    );

    await expect(get(url)).rejects.toMatchObject({
      name: 'HttpError',
      message: 'Channel not found',
      statusCode: 404,
      url,
      method: 'GET',
      data: { channel: 'unknown' },
    });
  });
});
