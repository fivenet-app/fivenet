import { beforeEach, describe, expect, it, vi } from 'vitest';
import { authUserTokenKey } from '~/stores/auth_session';
import { getGrpcAuthToken } from './auth';

const activeChar = { value: null as unknown | null };

vi.mock('~/composables/useAuth', () => ({
    useAuth: () => ({
        activeChar,
    }),
}));

describe('getGrpcAuthToken', () => {
    beforeEach(() => {
        activeChar.value = null;
        sessionStorage.clear();
    });

    it('returns null when no character is selected even if a token is stored', () => {
        sessionStorage.setItem(authUserTokenKey, 'stale-char-token');

        expect(getGrpcAuthToken()).toBeNull();
    });

    it('returns the stored token when a character is selected', () => {
        activeChar.value = { userId: 123 } as never;
        sessionStorage.setItem(authUserTokenKey, 'char-token');

        expect(getGrpcAuthToken()).toBe('char-token');
    });
});
