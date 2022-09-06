import { faker } from '@faker-js/faker';
import { rest } from 'msw';
import {
  globalEmotesEndpoint,
  oauth2TokenEndpoint,
  usersEndpoint,
} from '~/constants/twitch';
import {
  GetUsersResponse,
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
  rest.get<undefined, {}, GetUsersResponse>(
    usersEndpoint,
    async (req, res, ctx) => {
      const bearerToken = req.headers.get('Authorization');
      const clientId = req.headers.get('Client-Id');
      const loginNames = req.url.searchParams.getAll('login');

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

      if (loginNames.some((loginName) => loginName === '')) {
        return res(
          ctx.status(400),
          ctx.json({
            status: 400,
            error: 'Empty login name',
            message: 'Login name cannot be empty',
          })
        );
      }

      return res(
        ctx.status(200),
        ctx.json({
          data: [
            {
              id: '12826',
              login: 'twitch',
              display_name: 'Twitch',
              type: '',
              broadcaster_type: 'partner',
              description:
                'TwitchCon San Diego 2022: October 7 - 9, 2022. TwitchCon is baaaaack! San Diego, get ready to squad up.',
              profile_image_url:
                'https://static-cdn.jtvnw.net/jtv_user_pictures/26e28203-495c-4bb3-9c6e-c0f4c2a87a9d-profile_image-300x300.png',
              offline_image_url:
                'https://static-cdn.jtvnw.net/jtv_user_pictures/bdc19970-3a3b-4516-9f23-4203d59f0a5d-channel_offline_image-1920x1080.png',
              view_count: 336705975,
              created_at: '2007-05-22T10:39:54Z',
            },
            {
              id: '1005',
              login: 'abc',
              display_name: 'ABC',
              type: 'admin',
              broadcaster_type: 'affiliate',
              description: '',
              profile_image_url: '',
              offline_image_url: '',
              view_count: 9,
              created_at: '2015-08-25T19:58:00Z',
            },
            {
              id: '100534',
              login: 'abcdef',
              display_name: 'ABCDEF',
              type: 'staff',
              broadcaster_type: '',
              description: '',
              profile_image_url: '',
              offline_image_url: '',
              view_count: 23,
              created_at: '2013-04-25T13:43:00Z',
            },
            {
              id: '43993',
              login: '123456',
              display_name: '123456',
              type: 'global_mod',
              broadcaster_type: '',
              description: '',
              profile_image_url: '',
              offline_image_url: '',
              view_count: 343,
              created_at: '2017-10-01T09:36:00Z',
            },
          ],
        })
      );
    }
  ),
];

export default handlers;
