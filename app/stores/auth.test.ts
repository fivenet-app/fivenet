import { createPinia, setActivePinia } from 'pinia';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useAuthStore } from './auth';

const mocks = vi.hoisted(() => ({
    chooseCharacter: vi.fn(),
    authSessionStore: {
        getUserToken: vi.fn(),
        setUserToken: vi.fn(),
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
        mocks.authSessionStore.getUserToken.mockReturnValue(null);
    });

    it('clears account-level config-admin when selecting a character', async () => {
        mocks.authSessionStore.getUserToken.mockReturnValue('char-token');
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
        authStore.setAccountCanBeConfigAdmin(true);

        await authStore.chooseCharacter(123, false);

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
});
