import { describe, expect, it } from 'vitest';
import { shouldIgnoreClientError } from './client-errors';

describe('shouldIgnoreClientError', () => {
    it('ignores AsyncData requests cancelled by deduplication', () => {
        expect(shouldIgnoreClientError(new Error('AsyncData request cancelled by deduplication'))).toBe(true);
        expect(shouldIgnoreClientError('AsyncData request cancelled by deduplication')).toBe(true);
    });

    it('does not ignore unrelated errors', () => {
        expect(shouldIgnoreClientError(new Error('Request failed'))).toBe(false);
    });
});
