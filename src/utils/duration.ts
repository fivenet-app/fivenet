import * as googleProtobufDuration from '~~/gen/ts/google/protobuf/duration';

export function toDuration(input: string): googleProtobufDuration.Duration {
    const split = input.split('.');
    console.log(split[1] !== undefined ? parseInt(split[1].replace(/\D/g, '')) * 10000000 : 0);
    return {
        seconds: BigInt(split[0].replace(/\D/g, '')),
        nanos: split[1] !== undefined ? parseInt(split[1].replace(/\D/g, '')) * 10000000 : 0,
    };
}
