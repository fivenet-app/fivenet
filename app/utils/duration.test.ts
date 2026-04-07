import { describe, expect, it } from 'vitest';
import { clampDuration, durationToSeconds, fromDuration, secondsToDuration, toDuration } from './duration';

describe('toDuration', () => {
    it('should convert a number to a google.protobuf.Duration object', () => {
        const input = 123.45;
        const result = toDuration(input);

        expect(result).toEqual({ seconds: 123, nanos: 450000000 });
    });

    it('should convert a string to a google.protobuf.Duration object', () => {
        const input = '123.45';
        const result = toDuration(input);

        expect(result).toEqual({ seconds: 123, nanos: 450000000 });
    });

    it('should handle integers correctly', () => {
        const input = 123;
        const result = toDuration(input);

        expect(result).toEqual({ seconds: 123, nanos: 0 });
    });
});

describe('fromDuration', () => {
    it('should convert a google.protobuf.Duration object to a number', () => {
        const input = { seconds: 123, nanos: 450000000 };
        const result = fromDuration(input);

        expect(result).toBe(123.45);
    });

    it('should handle undefined input', () => {
        const result = fromDuration();

        expect(result).toBe(0);
    });

    it('should handle zero nanos correctly', () => {
        const input = { seconds: 123, nanos: 0 };
        const result = fromDuration(input);

        expect(result).toBe(123);
    });
});

describe('durationToSeconds / secondsToDuration', () => {
    it('should preserve precision for nanos', () => {
        const value = { seconds: 3, nanos: 1 };
        const seconds = durationToSeconds(value);

        expect(seconds).toBe(3.000000001);
        expect(secondsToDuration(seconds)).toEqual(value);
    });

    it('should normalize rounded nanos overflow', () => {
        const result = secondsToDuration(1.9999999996);
        expect(result).toEqual({ seconds: 2, nanos: 0 });
    });
});

describe('clampDuration', () => {
    const min = { seconds: 60, nanos: 0 };
    const max = { seconds: 3600, nanos: 0 };

    it('should clamp to minimum', () => {
        const result = clampDuration({ seconds: 30, nanos: 0 }, min, max);
        expect(result).toEqual(min);
    });

    it('should clamp to maximum', () => {
        const result = clampDuration({ seconds: 7200, nanos: 0 }, min, max);
        expect(result).toEqual(max);
    });

    it('should keep in-range values unchanged', () => {
        const value = { seconds: 120, nanos: 500000000 };
        const result = clampDuration(value, min, max);
        expect(result).toEqual(value);
    });

    it('should handle equal bounds', () => {
        const bound = { seconds: 180, nanos: 0 };
        const result = clampDuration({ seconds: 999, nanos: 0 }, bound, bound);
        expect(result).toEqual(bound);
    });
});
