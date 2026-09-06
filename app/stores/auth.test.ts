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

    it('shares an in-flight character restore instead of resolving concurrent calls early', async () => {
        mocks.authSessionStore.getUserToken.mockReturnValue('char-token');
        mocks.authSessionStore.userInfo.accountId = 123;
        mocks.authSessionStore.userInfo.userId = 123;

        let resolveChooseCharacter: ((value: unknown) => void) | undefined;
        mocks.chooseCharacter.mockReturnValueOnce(
            new Promise((resolve) => {
                resolveChooseCharacter = resolve;
            }),
        );

        const authStore = useAuthStore();
        const first = authStore.chooseCharacter(123, false);
        const second = authStore.chooseCharacter(123, false);

        await vi.waitFor(() => expect(mocks.chooseCharacter).toHaveBeenCalledTimes(1));

        resolveChooseCharacter?.({
            response: {
                username: 'tester',
                token: 'char-token',
                char: { userId: 123 } as never,
                permissions: [],
                attributes: [],
                jobProps: undefined,
            },
        });

        await first;
        await second;
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

    it('clears only in-memory state after a confirmed account-token failure', async () => {
        mocks.refreshAccountSession.mockRejectedValueOnce({ code: 'UNAUTHENTICATED' });

        const authStore = useAuthStore();
        authStore.lastCharID = 123;
        authStore.accountId = 456;
        authStore.setUsername('tester', false);
        authStore.setActiveChar({ userId: 123 } as never);
        authStore.phase = 'anonymous';

        await expect(authStore.ensureAccountSession()).resolves.toEqual({ kind: 'needs-login' });

        expect(authStore.phase).toBe('anonymous');
        expect(authStore.username).toBeNull();
        expect(authStore.activeChar).toBeNull();
        expect(authStore.lastCharID).toBe(123);
        expect(mocks.authSessionStore.setUserToken).toHaveBeenLastCalledWith(null);
    });

    it('keeps a temporary account-session failure out of the login path', async () => {
        mocks.refreshAccountSession.mockRejectedValueOnce(new Error('network unavailable'));

        const authStore = useAuthStore();

        await expect(authStore.ensureAccountSession()).resolves.toMatchObject({ kind: 'temporary-failure' });

        expect(authStore.phase).toBe('bootstrapping');
        expect(authStore.lastFailure?.kind).toBe('temporary');
    });

    it('requires character selection when the account has no remembered character', async () => {
        const authStore = useAuthStore();

        await expect(authStore.ensureCharacterSession()).resolves.toEqual({ kind: 'needs-character' });
        expect(mocks.chooseCharacter).not.toHaveBeenCalled();
    });

    it('invalidates the account when character restoration is rejected', async () => {
        mocks.chooseCharacter.mockRejectedValueOnce({ code: 'UNAUTHENTICATED' });

        const authStore = useAuthStore();
        authStore.lastCharID = 123;

        await expect(authStore.ensureCharacterSession()).resolves.toEqual({ kind: 'needs-login' });

        expect(authStore.phase).toBe('anonymous');
        expect(authStore.lastCharID).toBe(123);
        expect(mocks.authSessionStore.setUserToken).toHaveBeenLastCalledWith(null);
    });

    it('keeps temporary character restoration failures out of the selector redirect path', async () => {
        mocks.chooseCharacter.mockRejectedValueOnce(new Error('network unavailable'));

        const authStore = useAuthStore();
        authStore.lastCharID = 123;

        await expect(authStore.ensureCharacterSession()).resolves.toMatchObject({ kind: 'temporary-failure' });

        expect(authStore.phase).toBe('account-ready');
        expect(authStore.lastFailure?.kind).toBe('temporary');
    });

    it('shares an in-flight account restoration', async () => {
        let resolveRefresh: ((value: unknown) => void) | undefined;
        mocks.refreshAccountSession.mockReturnValueOnce(
            new Promise((resolve) => {
                resolveRefresh = resolve;
            }),
        );

        const authStore = useAuthStore();
        const first = authStore.ensureAccountSession();
        const second = authStore.ensureAccountSession();

        await vi.waitFor(() => expect(mocks.refreshAccountSession).toHaveBeenCalledTimes(1));
        resolveRefresh?.({
            response: {
                accountId: 123,
                canBeConfigAdmin: false,
                username: 'tester',
            },
        });

        await expect(first).resolves.toEqual({ kind: 'ready' });
        await expect(second).resolves.toEqual({ kind: 'ready' });
    });
});
