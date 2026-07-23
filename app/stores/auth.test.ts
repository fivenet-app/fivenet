import { createPinia, setActivePinia } from 'pinia';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useAuthStore } from './auth';

const mocks = vi.hoisted(() => ({
    chooseCharacter: vi.fn(),
    refreshAccountSession: vi.fn(),
    authSessionStore: {
        getUserToken: vi.fn(),
        setUserToken: vi.fn(),
        userInfo: { accountId: null as number | null, userId: null as number | null },
    },
    notificationsStore: {
        add: vi.fn(),
        logger: {
            log: vi.fn(),
            debug: vi.fn(),
            info: vi.fn(),
            warn: vi.fn(),
            error: vi.fn(),
        },
    },
    settingsStore: {
        startpage: '/overview',
    },
    grpcwsTransport: {
        updateUserToken: vi.fn(),
        close: vi.fn(),
    },
    webSocket: {
        status: { value: 'OPEN' },
        open: vi.fn(),
        close: vi.fn(),
        send: vi.fn().mockResolvedValue(true),
        data: { value: null },
    },
}));

vi.mock('~~/gen/ts/clients', () => ({
    getAuthAuthClient: vi.fn(async () => ({
        chooseCharacter: mocks.chooseCharacter,
        refreshAccountSession: mocks.refreshAccountSession,
    })),
}));

vi.mock('~/stores/auth_session', () => ({
    useAuthSessionStore: () => mocks.authSessionStore,
}));

vi.mock('~/stores/notifications', () => ({
    logger: mocks.notificationsStore.logger,
    useNotificationsStore: () => mocks.notificationsStore,
}));

vi.mock('~/stores/settings', () => ({
    useSettingsStore: () => mocks.settingsStore,
}));

vi.mock('~/composables/grpcws', () => ({
    useGRPCWebsocketTransport: () => mocks.grpcwsTransport,
}));

vi.mock('~/composables/grpcws/bridge', () => ({
    webSocket: mocks.webSocket,
}));

describe('useAuthStore', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
        vi.clearAllMocks();
        mocks.webSocket.status.value = 'OPEN';
        mocks.authSessionStore.getUserToken.mockReturnValue(null);
        mocks.authSessionStore.userInfo.accountId = null;
        mocks.authSessionStore.userInfo.userId = null;
        mocks.refreshAccountSession.mockResolvedValue({
            response: {
                accountId: 123,
                canBeConfigAdmin: false,
                username: 'tester',
            },
        });
    });

    it('clears account-level config-admin when selecting a character', async () => {
        mocks.authSessionStore.getUserToken.mockReturnValue('char-token');
        mocks.authSessionStore.userInfo.accountId = 123;
        mocks.authSessionStore.userInfo.userId = 123;
        mocks.chooseCharacter.mockResolvedValueOnce({
            response: {
                username: 'tester',
                token: 'char-token',
                char: { userId: 123 } as never,
                permissions: [],
                attributes: [],
                jobProps: undefined,
            },
        });

        const authStore = useAuthStore();
        authStore.accountId = 123;
        authStore.setAccountCanBeConfigAdmin(true);

        await authStore.chooseCharacter(123, false);

        expect(authStore.accountId).toBe(123);
        expect(authStore.canBeConfigAdmin).toBe(false);
        expect(mocks.chooseCharacter).toHaveBeenCalledWith(
            { charId: 123 },
            {
                meta: {
                    Authorization: 'Bearer char-token',
                },
            },
        );
    });

    it('does not reuse a character token from another account when selecting a character', async () => {
        mocks.authSessionStore.getUserToken.mockReturnValue('stale-char-token');
        mocks.authSessionStore.userInfo.accountId = 123;
        mocks.authSessionStore.userInfo.userId = 999;
        mocks.chooseCharacter.mockResolvedValueOnce({
            response: {
                username: 'tester',
                token: 'fresh-token',
                char: { userId: 123 } as never,
                permissions: [],
                attributes: [],
                jobProps: undefined,
            },
        });

        const authStore = useAuthStore();
        authStore.accountId = 456;

        await authStore.chooseCharacter(123, false);

        expect(authStore.accountId).toBe(123);
        expect(mocks.chooseCharacter).toHaveBeenCalledWith({ charId: 123 }, undefined);
    });

    it('refreshes the account session before restoring a character when accountId is missing', async () => {
        mocks.authSessionStore.getUserToken.mockReturnValue('char-token');
        mocks.authSessionStore.userInfo.accountId = 123;
        mocks.authSessionStore.userInfo.userId = 123;
        mocks.chooseCharacter.mockResolvedValueOnce({
            response: {
                username: 'tester',
                token: 'char-token',
                char: { userId: 123 } as never,
                permissions: [],
                attributes: [],
                jobProps: undefined,
            },
        });

        const authStore = useAuthStore();
        authStore.accountId = null;

        await authStore.chooseCharacter(123, false);

        expect(mocks.refreshAccountSession).toHaveBeenCalled();
        expect(mocks.chooseCharacter).toHaveBeenCalledWith(
            { charId: 123 },
            {
                meta: {
                    Authorization: 'Bearer char-token',
                },
            },
        );
    });

    it('does not open the websocket when refreshing account session metadata', async () => {
        mocks.webSocket.status.value = 'CLOSED';

        const authStore = useAuthStore();
        authStore.setUsername('persisted-user', false);

        await authStore.refreshAccountSession();

        expect(authStore.username).toBe('tester');
        expect(mocks.webSocket.open).not.toHaveBeenCalled();
    });

    it('opens the websocket when restoring an account-only session', async () => {
        mocks.webSocket.status.value = 'CLOSED';

        const authStore = useAuthStore();
        authStore.setUsername('persisted-user', false);

        await authStore.restoreAccountSession();

        expect(authStore.username).toBe('tester');
        expect(mocks.webSocket.open).toHaveBeenCalledTimes(1);
    });
});
