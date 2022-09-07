import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';

dayjs.extend(relativeTime);

export const timeFromNow = (time: dayjs.ConfigType) => dayjs(time).fromNow();
