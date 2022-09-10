import dayjs from 'dayjs';
import localizedFormat from 'dayjs/plugin/localizedFormat';
import relativeTime from 'dayjs/plugin/relativeTime';

dayjs.extend(localizedFormat);
dayjs.extend(relativeTime);

/**
 * @see https://day.js.org/docs/en/display/format
 * @param time
 * @param format
 * @returns Formatted date based on format tokens
 */
export const formatDate = (time: dayjs.ConfigType, format: string) =>
  dayjs(time).format(format);

export const timeFromNow = (time: dayjs.ConfigType) => dayjs(time).fromNow();
