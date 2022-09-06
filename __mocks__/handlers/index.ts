import apiHandlers from './api-handlers';
import twitchHandlers from './twitch-handlers';

const handlers = [...apiHandlers, ...twitchHandlers];

export default handlers;
