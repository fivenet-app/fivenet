import { describe, expect, it } from 'vitest';
import { shouldIgnoreClientError } from './client-errors';

describe('shouldIgnoreClientError', () => {
    it('ignores AsyncData requests cancelled by deduplication', () => {
        expect(shouldIgnoreClientError(new Error('AsyncData request cancelled by deduplication'))).toBe(true);
        expect(shouldIgnoreClientError('AsyncData request cancelled by deduplication')).toBe(true);
    });

    it('ignores cancelled gRPC errors', () => {
        expect(shouldIgnoreClientError({ code: 'CANCELLED', message: 'AsyncData request cancelled by `enabled: false`' })).toBe(
            true,
        );
        expect(shouldIgnoreClientError({ code: 'cancelled', message: 'Request cancelled' })).toBe(true);
    });

    it('does not ignore unrelated errors', () => {
        expect(shouldIgnoreClientError(new Error('Request failed'))).toBe(false);
        expect(shouldIgnoreClientError({ code: 'INTERNAL', message: 'Request failed' })).toBe(false);
    });
});
