import * as googleProtobufTimestamp from '~~/gen/ts/google/protobuf/timestamp';
import type { Timestamp as resourcesTimestampTimestamp } from '~~/gen/ts/resources/timestamp/timestamp';

export function toDate(ts: resourcesTimestampTimestamp | undefined, correction?: number): Date {
    if (ts === undefined || ts?.timestamp === undefined) return new Date();

    if (correction === undefined) return googleProtobufTimestamp.Timestamp.toDate(ts.timestamp!);

    return new Date(googleProtobufTimestamp.Timestamp.toDate(ts.timestamp!).getTime() - -correction);
}

export function stringToDate(time: string): Date {
    return new Date(Date.parse(time));
}

export function toTimestamp(date?: Date): resourcesTimestampTimestamp | undefined {
    if (date === undefined) return;

    return {
        timestamp: googleProtobufTimestamp.Timestamp.fromDate(date),
    };
}

export function toDatetimeLocal(date: Date): string {
    return new Date(date.getTime() + date.getTimezoneOffset() * -60 * 1000).toISOString().slice(0, 16);
}

export function dateToDateString(date: Date): string {
    const d = new Date(date);
    let month = '' + (d.getMonth() + 1);
    let day = '' + d.getDate();
    const year = d.getFullYear();

    if (month.length < 2) {
        month = '0' + month;
    }
    if (day.length < 2) {
        day = '0' + day;
    }

    return [year, month, day].join('-');
}

export function getWeekNumber(date: Date): number {
    const d = new Date(date);
    const dayNum = d.getUTCDay() || 7;
    d.setUTCDate(d.getUTCDate() + 4 - dayNum);
    const yearStart = new Date(Date.UTC(d.getUTCFullYear(), 0, 1));
    return Math.ceil(((d.getTime() - yearStart.getTime()) / 86400000 + 1) / 7);
}
