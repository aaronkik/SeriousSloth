import { faker } from '@faker-js/faker';
import { rest } from 'msw';
import { globalEmotesEndpoint, oauth2TokenEndpoint } from '~/constants/twitch';
import {
  GlobalEmotesResponse,
  OAuthClientCredentialsResponse,
} from '~/types/twitch';

const handlers = [
  rest.post<string, {}, OAuthClientCredentialsResponse>(
    oauth2TokenEndpoint,
    async (req, res, ctx) => {
      const reqBody = await req.text();
      const urlParams = new URLSearchParams(reqBody);

      if (!urlParams.get('client_id')) {
        return res(
          ctx.status(401),
          ctx.json({ status: 401, message: 'Missing client_id' })
        );
      }

      if (!urlParams.get('client_secret')) {
        return res(
          ctx.status(401),
          ctx.json({ status: 401, message: 'Missing client_secret' })
        );
      }

      if (urlParams.get('grant_type') !== 'client_credentials') {
        return res(
          ctx.status(401),
          ctx.json({
            status: 401,
            message: 'Grant type not set to client_credentials',
          })
        );
      }

      return res(
        ctx.status(200),
        ctx.json({
          access_token: faker.random.alphaNumeric(10),
          expires_in: faker.datatype.number({ min: 1, max: 100000 }),
          token_type: 'bearer',
        })
      );
    }
  ),
  rest.get<undefined, {}, GlobalEmotesResponse>(
    globalEmotesEndpoint,
    async (req, res, ctx) => {
      const bearerToken = req.headers.get('Authorization');
      const clientId = req.headers.get('Client-Id');

      if (!bearerToken) {
        return res(
          ctx.status(401),
          ctx.json({
            status: 401,
            error: 'Missing access token',
            message: 'Missing access token',
          })
        );
      }
      if (!clientId) {
        return res(
          ctx.status(401),
          ctx.json({
            status: 401,
            error: 'Missing client id',
            message: 'Missing client id',
          })
        );
      }

      return res(
        ctx.status(200),
        ctx.json({
          data: [
            {
              id: '3',
              name: ':D',
              images: {
                url_1x:
                  'https://static-cdn.jtvnw.net/emoticons/v2/3/static/light/1.0',
                url_2x:
                  'https://static-cdn.jtvnw.net/emoticons/v2/3/static/light/2.0',
                url_4x:
                  'https://static-cdn.jtvnw.net/emoticons/v2/3/static/light/3.0',
              },
              format: ['animated'],
              scale: ['1.0', '2.0'],
              theme_mode: ['light'],
            },
            {
              id: '2',
              name: ':(',
              images: {
                url_1x:
                  'https://static-cdn.jtvnw.net/emoticons/v2/2/static/light/1.0',
                url_2x:
                  'https://static-cdn.jtvnw.net/emoticons/v2/2/static/light/2.0',
                url_4x:
                  'https://static-cdn.jtvnw.net/emoticons/v2/2/static/light/3.0',
              },
              format: ['static'],
              scale: ['1.0', '2.0', '3.0'],
              theme_mode: ['light', 'dark'],
            },
            {
              id: '1',
              name: ':)',
              images: {
                url_1x:
                  'https://static-cdn.jtvnw.net/emoticons/v2/1/static/light/1.0',
                url_2x:
                  'https://static-cdn.jtvnw.net/emoticons/v2/1/static/light/2.0',
                url_4x:
                  'https://static-cdn.jtvnw.net/emoticons/v2/1/static/light/3.0',
              },
              format: ['static'],
              scale: ['1.0', '2.0', '3.0'],
              theme_mode: ['light', 'dark'],
            },
          ],
          template:
            'https://static-cdn.jtvnw.net/emoticons/v2/{{id}}/{{format}}/{{theme_mode}}/{{scale}}',
        })
      );
    }
  ),
];

export default handlers;
