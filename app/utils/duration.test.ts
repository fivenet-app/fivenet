import { describe, expect, it } from 'vitest';
import { fromDuration, toDuration } from './duration';

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
