import { beforeEach, describe, expect, it, vi } from 'vitest';
import { getGrpcCharacterAuthToken, getGrpcWebsocketAuthToken } from './auth';

const mocks = vi.hoisted(() => ({
    token: null as string | null,
}));

vi.mock('~/stores/auth_session', () => ({
    useAuthSessionStore: () => ({
        getUserToken: () => mocks.token,
    }),
}));

describe('grpc auth token accessors', () => {
    beforeEach(() => {
        mocks.token = null;
    });

    it('returns the in-memory token for character-scoped auth', () => {
        mocks.token = 'char-token';

        expect(getGrpcCharacterAuthToken()).toBe('char-token');
    });

    it('keeps websocket auth tokenless without a character session', () => {
        expect(getGrpcWebsocketAuthToken()).toBeNull();
    });

    it('returns the in-memory websocket token when available', () => {
        mocks.token = 'char-token';

        expect(getGrpcWebsocketAuthToken()).toBe('char-token');
    });
});
