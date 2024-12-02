import type * as googleProtobufDuration from '~~/gen/ts/google/protobuf/duration';

export function toDuration(input: string | number): googleProtobufDuration.Duration {
    if (typeof input === 'number') {
        input = input.toFixed(2);
    }

    const split = input.split('.');
    return {
        seconds: split[0] !== undefined ? parseInt(split[0].replace(/\D/g, '')) : 1,
        nanos: split[1] !== undefined ? parseInt(split[1].replace(/\D/g, '')) * 1_000_000 : 0,
    };
}

export function fromDuration(input?: googleProtobufDuration.Duration): number {
    return parseFloat((input ? input?.seconds.toString() : '0') + '.' + (input?.nanos ?? 0) / 1000000);
}
