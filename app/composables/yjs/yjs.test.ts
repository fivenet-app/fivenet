import type { DuplexStreamingCall, RpcOptions } from '@protobuf-ts/runtime-rpc';
import { describe, expect, it, vi } from 'vitest';
import * as Y from 'yjs';
import type { ServerPacket } from '~~/gen/ts/resources/collab/collab';
import GrpcProvider, { type StreamConnectFn } from './yjs';

class TestResponses {
    private readonly messageListeners: Array<(message: ServerPacket) => void> = [];
    private readonly errorListeners: Array<(error: Error) => void> = [];

    onMessage(listener: (message: ServerPacket) => void) {
        this.messageListeners.push(listener);
        return () => undefined;
    }

    onError(listener: (error: Error) => void) {
        this.errorListeners.push(listener);
        return () => undefined;
    }

    emitMessage(message: ServerPacket) {
        this.messageListeners.forEach((listener) => listener(message));
    }
}

function createCall(responses: TestResponses): DuplexStreamingCall<object, ServerPacket> {
    return {
        requests: {
            send: vi.fn().mockResolvedValue(undefined),
            complete: vi.fn().mockResolvedValue(undefined),
        },
        responses,
    } as unknown as DuplexStreamingCall<object, ServerPacket>;
}

describe('GrpcProvider', () => {
    it('abandons a collab stream that never handshakes and retries it', async () => {
        vi.useFakeTimers();

        try {
            const options: RpcOptions[] = [];
            const connect = vi.fn((opts?: RpcOptions) => {
                options.push(opts ?? {});
                return createCall(new TestResponses());
            }) as StreamConnectFn;
            const provider = new GrpcProvider(new Y.Doc(), connect, { targetId: 42 });

            provider.connect();
            expect(connect).toHaveBeenCalledTimes(1);

            await vi.advanceTimersByTimeAsync(15_000);
            expect(options[0]?.abort?.aborted).toBe(true);

            await vi.advanceTimersByTimeAsync(750);
            expect(connect).toHaveBeenCalledTimes(2);

            provider.destroy();
        } finally {
            vi.useRealTimers();
        }
    });

    it('keeps a successfully handshaken collab stream connected', async () => {
        vi.useFakeTimers();

        try {
            const responses = new TestResponses();
            const connect = vi.fn(() => createCall(responses)) as StreamConnectFn;
            const provider = new GrpcProvider(new Y.Doc(), connect, { targetId: 42 });

            provider.connect();
            responses.emitMessage({
                senderId: 1n,
                msg: { oneofKind: 'handshake', handshake: { clientId: 1n } },
            });

            await vi.advanceTimersByTimeAsync(15_000);
            expect(provider.isSynced).toBe(false);
            expect(connect).toHaveBeenCalledTimes(1);

            provider.destroy();
        } finally {
            vi.useRealTimers();
        }
    });
});
