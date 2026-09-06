import { beforeEach, describe, expect, it, vi } from 'vitest';

const mocks = vi.hoisted(() => ({
    authStore: {
        invalidateSession: vi.fn(),
        clearCharacterSession: vi.fn(),
    },
    notificationsStore: {
        add: vi.fn(),
    },
}));

vi.mock('~/stores/auth', () => ({
    useAuthStore: () => mocks.authStore,
}));

vi.mock('~/stores/notifications', () => ({
    useNotificationsStore: () => mocks.notificationsStore,
}));

function rpcError(code: string, message: string = 'server rejected credentials') {
    return {
        name: 'RpcError',
        message,
        code,
        meta: {},
        serviceName: 'test.AuthService',
        methodName: 'Test',
    };
}

describe('handleGRPCError', () => {
    let handleGRPCError: (typeof import('./grpc'))['handleGRPCError'];

    beforeEach(async () => {
        vi.clearAllMocks();
        vi.resetModules();
        ({ handleGRPCError } = await import('./grpc'));
    });

    it('invalidates the account by default for unauthenticated RPC responses', () => {
        handleGRPCError(rpcError('UNAUTHENTICATED'));

        expect(mocks.authStore.invalidateSession).toHaveBeenCalledWith('account-expired');
        expect(mocks.authStore.clearCharacterSession).not.toHaveBeenCalled();
    });

    it('can limit unauthenticated failures to character state', () => {
        handleGRPCError(rpcError('UNAUTHENTICATED'), { authScope: 'character' });

        expect(mocks.authStore.clearCharacterSession).toHaveBeenCalledWith('character-expired');
        expect(mocks.authStore.invalidateSession).not.toHaveBeenCalled();
    });

    it('can report unauthenticated failures without mutating auth state', () => {
        handleGRPCError(rpcError('UNAUTHENTICATED'), { authScope: 'none' });

        expect(mocks.authStore.clearCharacterSession).not.toHaveBeenCalled();
        expect(mocks.authStore.invalidateSession).not.toHaveBeenCalled();
        expect(mocks.notificationsStore.add).toHaveBeenCalledTimes(1);
    });

    it('clears character state for an ErrCharLock response without navigating', () => {
        handleGRPCError(rpcError('PERMISSION_DENIED', 'ErrCharLock'));

        expect(mocks.authStore.clearCharacterSession).toHaveBeenCalledWith('character-expired');
        expect(mocks.notificationsStore.add).toHaveBeenCalledTimes(1);
    });
});
