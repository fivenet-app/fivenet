import * as resources_timestamp_timestamp_pb from '@arpanet/gen/resources/timestamp/timestamp_pb';

export function getSecondsFormattedAsDuration(seconds: number): string {
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

export function getDate(ts: resources_timestamp_timestamp_pb.Timestamp | undefined): undefined | Date {
    if (typeof ts === undefined) {
        return new Date();
    }
    return ts?.getTimestamp()?.toDate();
}

export function getDateLocaleString(ts: resources_timestamp_timestamp_pb.Timestamp | undefined): undefined | string {
    if (typeof ts === undefined) {
        return "-";
    }
    return ts?.getTimestamp()?.toDate().toLocaleString();
}
