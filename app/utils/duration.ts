import type * as googleProtobufDuration from '~~/gen/ts/google/protobuf/duration';

const NANOSECONDS_PER_SECOND = 1_000_000_000;

function normalizeDuration(seconds: number, nanos: number): googleProtobufDuration.Duration {
    if (nanos >= NANOSECONDS_PER_SECOND) {
        seconds += Math.floor(nanos / NANOSECONDS_PER_SECOND);
        nanos %= NANOSECONDS_PER_SECOND;
    } else if (nanos <= -NANOSECONDS_PER_SECOND) {
        seconds += Math.ceil(nanos / NANOSECONDS_PER_SECOND);
        nanos %= NANOSECONDS_PER_SECOND;
    }

    if (seconds > 0 && nanos < 0) {
        seconds -= 1;
        nanos += NANOSECONDS_PER_SECOND;
    } else if (seconds < 0 && nanos > 0) {
        seconds += 1;
        nanos -= NANOSECONDS_PER_SECOND;
    }

    return {
        seconds,
        nanos,
    };
}

export function durationToSeconds(input?: googleProtobufDuration.Duration): number {
    if (!input) {
        return 0;
    }

    return input.seconds + input.nanos / NANOSECONDS_PER_SECOND;
}

export function secondsToDuration(seconds: number): googleProtobufDuration.Duration {
    if (!Number.isFinite(seconds)) {
        return { seconds: 0, nanos: 0 };
    }

    const wholeSeconds = Math.trunc(seconds);
    const nanos = Math.round((seconds - wholeSeconds) * NANOSECONDS_PER_SECOND);

    return normalizeDuration(wholeSeconds, nanos);
}

export function clampDuration(
    duration: googleProtobufDuration.Duration,
    min?: googleProtobufDuration.Duration,
    max?: googleProtobufDuration.Duration,
): googleProtobufDuration.Duration {
    const durationSeconds = durationToSeconds(duration);
    const minSeconds = min ? durationToSeconds(min) : undefined;
    const maxSeconds = max ? durationToSeconds(max) : undefined;

    if (minSeconds !== undefined && maxSeconds !== undefined && minSeconds > maxSeconds) {
        const clamped = Math.min(Math.max(durationSeconds, maxSeconds), minSeconds);
        return secondsToDuration(clamped);
    }

    let clamped = durationSeconds;
    if (minSeconds !== undefined) {
        clamped = Math.max(clamped, minSeconds);
    }
    if (maxSeconds !== undefined) {
        clamped = Math.min(clamped, maxSeconds);
    }

    return secondsToDuration(clamped);
}

export function toDuration(seconds: string | number): googleProtobufDuration.Duration {
    const parsed = typeof seconds === 'number' ? seconds : parseFloat(seconds);
    if (!Number.isFinite(parsed)) {
        return { seconds: 0, nanos: 0 };
    }

    return secondsToDuration(parsed);
}

export function fromDuration(input?: googleProtobufDuration.Duration): number {
    return durationToSeconds(input);
}
