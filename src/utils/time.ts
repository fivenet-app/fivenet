import TimeAgo from 'javascript-time-ago';
import en from 'javascript-time-ago/locale/en';
import * as resources_timestamp_timestamp_pb from '@fivenet/gen/resources/timestamp/timestamp_pb';

TimeAgo.addDefaultLocale(en);
const timeAgo = new TimeAgo('en-US');

export function fromSecondsToFormattedDuration(seconds: number): string {
    var w = Math.floor(seconds / (7 * (3600 * 24)));
    var d = Math.floor(seconds / (3600 * 24));
    var h = Math.floor((seconds % (3600 * 24)) / 3600);
    var m = Math.floor((seconds % 3600) / 60);
    var s = Math.floor(seconds % 60);

    const parts = new Array<string>();
    if (w > 0) {
        parts.push(w + (w == 1 ? ' week' : ' weeks'));
    }
    if (d > 0) {
        parts.push(d + (d == 1 ? ' day' : ' days'));
    }
    if (h > 0) {
        parts.push(h + (h == 1 ? ' hour' : ' hours'));
    }
    if (m > 0) {
        parts.push(m + (m == 1 ? ' minute' : ' minutes'));
    }
    if (s > 0) {
        parts.push(s + (s == 1 ? ' second' : ' seconds'));
    }
    return parts.join(', ');
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

export function fromString(time: string): undefined | Date {
    return new Date(Date.parse(time));
}
