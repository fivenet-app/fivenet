import { beforeEach, describe, expect, it, vi } from 'vitest';
import { getGrpcCharacterAuthToken, getGrpcWebsocketAuthToken } from './auth';

const mocks = vi.hoisted(() => ({
    activeChar: { value: null as unknown | null },
    authUserTokenKey: 'fivenet:user_token_v1',
}));

vi.mock('../useAuth', () => ({
    useAuth: () => ({
        activeChar: mocks.activeChar,
    }),
}));

describe('grpc auth token accessors', () => {
    beforeEach(() => {
        mocks.activeChar.value = null;
        sessionStorage.clear();
    });

    it('returns the stored token for character-scoped auth even when no character is selected', () => {
        sessionStorage.setItem(mocks.authUserTokenKey, 'stale-char-token');

        expect(getGrpcCharacterAuthToken()).toBe('stale-char-token');
    });

    it('keeps websocket control auth tokenless in account-only mode', () => {
        sessionStorage.setItem(mocks.authUserTokenKey, 'stale-char-token');

        expect(getGrpcWebsocketAuthToken()).toBeNull();
    });

    it('returns the stored websocket token when a character is selected', () => {
        mocks.activeChar.value = { userId: 123 } as never;
        sessionStorage.setItem(mocks.authUserTokenKey, 'char-token');

        expect(getGrpcWebsocketAuthToken()).toBe('char-token');
    });
});
