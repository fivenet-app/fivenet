import type * as googleProtobufDuration from '~~/gen/ts/google/protobuf/duration';

export function toDuration(seconds: string | number): googleProtobufDuration.Duration {
    if (typeof seconds === 'number') {
        seconds = seconds.toFixed(2);
    }

    const split = seconds.split('.');
    return {
        seconds: split[0] !== undefined ? parseInt(split[0].replace(/\D/g, ''), 10) : 1,
        nanos: split[1] !== undefined ? Math.floor(parseFloat('0.' + split[1]) * 1_000_000_000) : 0,
    };
}

export function fromDuration(input?: googleProtobufDuration.Duration): number {
    return parseFloat((input ? input?.seconds.toString() : '0') + '.' + (input?.nanos ?? 0) / 1_000_000);
}
