import * as googleProtobufTimestamp from '~~/gen/ts/google/protobuf/timestamp';
import * as resourcesTimestampTimestamp from '~~/gen/ts/resources/timestamp/timestamp';

const secondsPerMinute = 60;
const secondsPerHour = secondsPerMinute * 60;
const secondsPerDay = secondsPerHour * 24;
const secondsPerWeek = secondsPerDay * 7;
const secondsPerYear = secondsPerWeek * 52;

export function fromSecondsToFormattedDuration(seconds: number, options?: { seconds?: boolean; emptyText?: string }): string {
    const { t } = useI18n();

    const years = Math.floor(seconds / secondsPerYear);
    seconds -= years * secondsPerYear;
    const weeks = Math.floor(seconds / secondsPerWeek);
    seconds -= weeks * secondsPerWeek;
    const days = Math.floor(seconds / secondsPerDay);
    seconds -= days * secondsPerDay;
    const hours = Math.floor(seconds / secondsPerHour);
    seconds -= hours * secondsPerHour;
    const minutes = Math.floor(seconds / secondsPerMinute);
    seconds -= minutes * secondsPerMinute;

    const parts: String[] = [];
    if (years > 0) {
        parts.push(`${years} ${t(`common.time_ago.year`, years)}`);
    }
    if (weeks > 0) {
        parts.push(`${weeks} ${t(`common.time_ago.week`, weeks)}`);
    }
    if (days > 0) {
        parts.push(`${days} ${t(`common.time_ago.day`, days)}`);
    }
    if (hours > 0) {
        parts.push(`${hours} ${t(`common.time_ago.hour`, hours)}`);
    }
    if (minutes > 0) {
        parts.push(`${minutes} ${t(`common.time_ago.minute`, minutes)}`);
    }
    if ((!options || options.seconds) && seconds > 0) {
        parts.push(`${seconds} ${t(`common.time_ago.second`, seconds)}`);
    }

    const text = parts.join(', ');
    return text.length > 0 ? text : t(options?.emptyText ?? 'common.unknown');
}

export function toDate(ts: resourcesTimestampTimestamp.Timestamp | undefined): Date {
    if (ts === undefined || ts?.timestamp === undefined) {
        return new Date();
    }

    return googleProtobufTimestamp.Timestamp.toDate(ts.timestamp!);
}

export function toDateLocaleString(ts: resourcesTimestampTimestamp.Timestamp | undefined, d?: Function): string {
    if (ts === undefined || typeof ts === 'undefined') {
        return '-';
    }

    if (d) {
        return d(googleProtobufTimestamp.Timestamp.toDate(ts.timestamp!), 'short');
    }

    return googleProtobufTimestamp.Timestamp.toDate(ts.timestamp!).toLocaleDateString();
}

export function fromString(time: string): Date {
    return new Date(Date.parse(time));
}

export function toTimestamp(date?: Date): resourcesTimestampTimestamp.Timestamp | undefined {
    if (date === undefined) {
        return;
    }

    return {
        timestamp: googleProtobufTimestamp.Timestamp.fromDate(date),
    };
}

export function toDatetimeLocal(date: Date): string {
    return new Date(date.getTime() + date.getTimezoneOffset() * -60 * 1000).toISOString().slice(0, 19);
}
