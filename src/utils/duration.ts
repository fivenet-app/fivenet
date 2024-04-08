import * as googleProtobufDuration from '~~/gen/ts/google/protobuf/duration';

export function toDuration(input: string): googleProtobufDuration.Duration {
    const split = input.split('.');
    return {
        seconds: BigInt(split[0].replace(/\D/g, '')),
        nanos: split[1] !== undefined ? parseInt(split[1].replace(/\D/g, '')) * 10000000 : 0,
    };
}

export function fromDuration(input: googleProtobufDuration.Duration): string {
    return parseFloat(input.seconds.toString() + '.' + (input.nanos ?? 0) / 1000000).toString() + 's';
}
