import { beforeEach, describe, expect, it, vi } from 'vitest';
import { authUserTokenKey } from '~/stores/auth_session';
import { getGrpcRpcAuthToken, getGrpcWebsocketAuthToken } from './auth';

const activeChar = { value: null as unknown | null };

vi.mock('~/composables/useAuth', () => ({
    useAuth: () => ({
        activeChar,
    }),
}));

describe('grpc auth token accessors', () => {
    beforeEach(() => {
        activeChar.value = null;
        sessionStorage.clear();
    });

    it('returns the stored token for RPC auth even when no character is selected', () => {
        sessionStorage.setItem(authUserTokenKey, 'stale-char-token');

        expect(getGrpcRpcAuthToken()).toBe('stale-char-token');
    });

    it('keeps websocket control auth tokenless in account-only mode', () => {
        sessionStorage.setItem(authUserTokenKey, 'stale-char-token');

        expect(getGrpcWebsocketAuthToken()).toBeNull();
    });

    it('returns the stored websocket token when a character is selected', () => {
        activeChar.value = { userId: 123 } as never;
        sessionStorage.setItem(authUserTokenKey, 'char-token');

        expect(getGrpcWebsocketAuthToken()).toBe('char-token');
        expect(getGrpcRpcAuthToken()).toBe('char-token');
    });
});
