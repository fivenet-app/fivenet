import { beforeEach, describe, expect, it, vi } from 'vitest';
import { authUserTokenKey } from '~/stores/auth_session';
import { GrpcCombinedTransport } from './useGRPCTransport';

const activeChar = { value: null as unknown | null };

type MockUnaryClient = {
    mock: {
        calls: Array<[unknown, unknown, { meta: Record<string, string | undefined> }]>;
    };
};

vi.mock('~/composables/useAuth', () => ({
    useAuth: () => ({
        activeChar,
    }),
}));

function createTransport() {
    const unaryClient = {
        mergeOptions: vi.fn((options) => options),
        unary: vi.fn(),
    };
    const streamClient = {
        mergeOptions: vi.fn((options) => options),
        serverStreaming: vi.fn(),
        clientStreaming: vi.fn(),
        duplex: vi.fn(),
    };

    return {
        transport: new GrpcCombinedTransport(unaryClient as never, streamClient as never),
        unaryClient,
        streamClient,
    };
}

describe('GrpcCombinedTransport auth headers', () => {
    beforeEach(() => {
        activeChar.value = null;
        sessionStorage.clear();
    });

    it('keeps account-only unary calls tokenless when no character is active', () => {
        sessionStorage.setItem(authUserTokenKey, 'stale-char-token');
        const { transport, unaryClient } = createTransport();

        transport.unary(
            { name: 'GetAppConfig', service: { typeName: 'services.settings.ConfigService' } } as never,
            {} as never,
            {},
        );

        expect(unaryClient.unary).toHaveBeenCalled();
        expect(unaryClient.mergeOptions).toHaveBeenCalledWith({ meta: {} });
        const firstCall = (unaryClient.unary as MockUnaryClient).mock.calls[0];
        expect(firstCall[2].meta.Authorization).toBeUndefined();
    });

    it('sends the stored token for choose-character restore even when no character is active', () => {
        sessionStorage.setItem(authUserTokenKey, 'char-token');
        const { transport, unaryClient } = createTransport();

        transport.unary(
            { name: 'ChooseCharacter', service: { typeName: 'services.auth.AuthService' } } as never,
            {} as never,
            {},
        );

        const firstCall = (unaryClient.unary as MockUnaryClient).mock.calls[0];
        expect(firstCall[2].meta.Authorization).toBe('Bearer char-token');
    });

    it('sends the stored token for unary calls when a character is active', () => {
        activeChar.value = { userId: 123 } as never;
        sessionStorage.setItem(authUserTokenKey, 'char-token');
        const { transport, unaryClient } = createTransport();

        transport.unary(
            { name: 'GetCharacters', service: { typeName: 'services.auth.AuthService' } } as never,
            {} as never,
            {},
        );

        const firstCall = (unaryClient.unary as MockUnaryClient).mock.calls[0];
        expect(firstCall[2].meta.Authorization).toBe('Bearer char-token');
    });
});
