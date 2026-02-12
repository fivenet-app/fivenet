import { describe, expect, it } from 'vitest';
import { parseJWTPayload } from './jwt';

describe('parseJWTPayload', () => {
    it('should parse a valid JWT payload', async () => {
        const token = `${btoa('header')}.${btoa('{"key":"value"}')}.signature`;
        const result = parseJWTPayload<{ key: string }>(token);
        expect(result).toEqual({ key: 'value' });
    });

    it('should throw an error for an invalid JWT token', async () => {
        const token = 'invalid.token';
        expect(() => parseJWTPayload(token)).toThrow('Invalid JWT token');
    });

    it('should throw an error for a malformed payload', async () => {
        const token = `${btoa('header')}.invalidPayload.signature`;
        expect(() => parseJWTPayload(token)).toThrow();
    });
});
