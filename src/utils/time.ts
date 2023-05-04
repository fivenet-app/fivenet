import * as resources_timestamp_timestamp_pb from '@fivenet/gen/resources/timestamp/timestamp_pb';

export function fromSecondsToFormattedDuration(seconds: number): string {
    const { t } = useI18n();

    const w = Math.floor(seconds / (7 * (3600 * 24)));
    const d = Math.floor(seconds / (3600 * 24));
    const h = Math.floor((seconds % (3600 * 24)) / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = Math.floor(seconds % 60);

    const parts = new Array<string>();
    if (w > 0) {
        parts.push(`${w} ${t(`common.time_ago.week`, w)}`);
    }
    if (d > 0) {
        parts.push(`${d} ${t(`common.time_ago.day`, d)}`);
    }
    if (h > 0) {
        parts.push(`${h} ${t(`common.time_ago.hour`, h)}`);
    }
    if (m > 0) {
        parts.push(`${m} ${t(`common.time_ago.minute`, m)}`);
    }
    if (s > 0) {
        parts.push(`${s} ${t(`common.time_ago.second`, s)}`);
    }
    return parts.join(', ');
}

export function toDate(ts: resources_timestamp_timestamp_pb.Timestamp | undefined): undefined | Date {
    if (typeof ts === undefined) {
        return new Date();
    }
    return ts?.getTimestamp()?.toDate();
}

export function toDateLocaleString(ts: resources_timestamp_timestamp_pb.Timestamp | undefined, d?: Function): undefined | string {
    if (typeof ts === undefined) {
        return '-';
    }

    if (d) {
        return d(ts?.getTimestamp()?.toDate()!, 'short');
    }
    return ts?.getTimestamp()?.toDate().toLocaleDateString();
}

export function fromString(time: string): undefined | Date {
    return new Date(Date.parse(time));
}
