import TimeAgo from 'javascript-time-ago';
import en from 'javascript-time-ago/locale/en';
import * as resources_timestamp_timestamp_pb from '@arpanet/gen/resources/timestamp/timestamp_pb';

TimeAgo.addDefaultLocale(en);
const timeAgo = new TimeAgo('en-US');

export function fromSecondsToFormattedDuration(seconds: number): string {
    var w = Math.floor(seconds / (7 * (3600 * 24)));
    var d = Math.floor(seconds / (3600 * 24));
    var h = Math.floor((seconds % (3600 * 24)) / 3600);
    var m = Math.floor((seconds % 3600) / 60);
    var s = Math.floor(seconds % 60);

    var dWeeks = w > 0 ? w + (w == 1 ? ' week, ' : ' weeks, ') : '';
    var dDisplay = d > 0 ? d + (d == 1 ? ' day, ' : ' days, ') : '';
    var hDisplay = h > 0 ? h + (h == 1 ? ' hour, ' : ' hours, ') : '';
    var mDisplay = m > 0 ? m + (m == 1 ? ' minute, ' : ' minutes, ') : '';
    var sDisplay = s > 0 ? s + (s == 1 ? ' second' : ' seconds') : '';
    return dWeeks + dDisplay + hDisplay + mDisplay + sDisplay;
}

export function toDate(ts: resources_timestamp_timestamp_pb.Timestamp | undefined): undefined | Date {
    if (typeof ts === undefined) {
        return new Date();
    }
    return ts?.getTimestamp()?.toDate();
}

export function toDateLocaleString(ts: resources_timestamp_timestamp_pb.Timestamp | undefined): undefined | string {
    if (typeof ts === undefined) {
        return '-';
    }
    return ts?.getTimestamp()?.toDate().toLocaleString('de-DE');
}

export function toDateRelativeString(ts: resources_timestamp_timestamp_pb.Timestamp | undefined): undefined | string {
    if (typeof ts === undefined) {
        return '-';
    }

    const date = ts?.getTimestamp()?.toDate();
    if (!date) return;

    return timeAgo.format(date, 'round');
}

export function fromString(time: string): undefined | Date{
    return new Date(Date.parse(time));
}
