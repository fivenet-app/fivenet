import type { WebSocketStatus } from '@vueuse/core';
import { describe, expect, it, vi } from 'vitest';
import { ref, watch } from 'vue';
import { Metadata } from '~/composables/grpcws/metadata';
import type { ILogger } from '~/utils/logger';
import { GrpcFrame } from '~~/gen/ts/resources/grpcws/grpcws';
import type { TransportOptions } from '../transport';
import { WebsocketChannelImpl, type GrpcStream } from './websocketChannel';

function createLogger(): ILogger {
    return {
        log: vi.fn(),
        debug: vi.fn(),
        info: vi.fn(),
        warn: vi.fn(),
        error: vi.fn(),
    };
}

function createWebSocketStub() {
    return {
        data: ref<ArrayBuffer | null>(null),
        status: ref<WebSocketStatus>('OPEN'),
        send: vi.fn().mockResolvedValue(true),
        open: vi.fn(),
    };
}

function createGrpcStreamMock(): GrpcStream {
    return {
        streamId: 1,
        service: 'test.Service',
        method: 'TestMethod',
        isStream: false,
        closed: false,
        start: vi.fn(),
        sendMessage: vi.fn().mockResolvedValue(undefined),
        finishSend: vi.fn().mockResolvedValue(undefined),
        cancel: vi.fn().mockResolvedValue(undefined),
    };
}

function createAuthOkBuffer(): ArrayBuffer {
    return GrpcFrame.toBinary(
        GrpcFrame.create({
            streamId: 0,
            payload: {
                oneofKind: 'header',
                header: {
                    operation: 'auth_ok',
                    headers: {},
                    status: 200,
                },
            },
        }),
    ).buffer.slice(0) as ArrayBuffer;
}

function createFailureBuffer(streamId: number, message: string): ArrayBuffer {
    return GrpcFrame.toBinary(
        GrpcFrame.create({
            streamId,
            payload: {
                oneofKind: 'failure',
                failure: {
                    errorStatus: '1',
                    errorMessage: message,
                    headers: {},
                },
            },
        }),
    ).buffer.slice(0) as ArrayBuffer;
}

function expectFrame(sentFrames: GrpcFrame[], index: number): GrpcFrame {
    const frame = sentFrames.at(index);
    expect(frame).toBeDefined();
    return frame!;
}

describe('WebsocketChannelImpl', () => {
    it('reuses stream ids after a stream completes', async () => {
        const channel = new WebsocketChannelImpl(createLogger(), createWebSocketStub(), () => 'token');

        const allocated = Array.from({ length: 7 }, () => channel.getNextStreamId());
        expect(allocated).toEqual([1, 2, 3, 4, 5, 6, 7]);
        expect(() => channel.getNextStreamId()).toThrow('No available websocket stream ids');

        const onEnd = vi.fn();
        const stream = createGrpcStreamMock();
        channel.activeStreams.set(1, [
            {
                debug: false,
                onEnd,
            } as unknown as TransportOptions,
            stream,
        ]);

        const completeFrame = GrpcFrame.toBinary(
            GrpcFrame.create({
                streamId: 1,
                payload: {
                    oneofKind: 'complete',
                    complete: {},
                },
            }),
        );
        const completeBuffer = completeFrame.buffer.slice(
            completeFrame.byteOffset,
            completeFrame.byteOffset + completeFrame.byteLength,
        ) as ArrayBuffer;

        await channel.onMessage(completeBuffer);

        expect(onEnd).toHaveBeenCalledTimes(1);
        expect(channel.activeStreams.has(1)).toBe(false);
        expect(channel.getNextStreamId()).toBe(1);
    });

    it('authenticates with account-only auth first and upgrades when the token changes', async () => {
        const sentFrames: GrpcFrame[] = [];
        const webSocket = {
            data: ref<ArrayBuffer | null>(null),
            status: ref<WebSocketStatus>('OPEN'),
            send: vi.fn().mockImplementation((payload: ArrayBuffer) => {
                sentFrames.push(GrpcFrame.fromBinary(new Uint8Array(payload)));
                return true;
            }),
            open: vi.fn(),
        };

        let token: string | null = null;
        const channel = new WebsocketChannelImpl(createLogger(), webSocket, () => token);

        const accountAuth = channel.ensureAuthenticated();
        expect(sentFrames).toHaveLength(1);
        const authFrame = expectFrame(sentFrames, 0);
        expect(authFrame.streamId).toBe(0);
        expect(authFrame.payload.oneofKind).toBe('header');
        if (authFrame.payload.oneofKind !== 'header') throw new Error('Expected auth header frame');
        expect(authFrame.payload.header.operation).toBe('auth');
        expect(authFrame.payload.header.headers.Authorization).toBeUndefined();

        await channel.onMessage(createAuthOkBuffer());
        await accountAuth;

        token = 'char-token';
        const charAuth = channel.ensureAuthenticated();
        expect(sentFrames).toHaveLength(2);
        const reauthFrame = expectFrame(sentFrames, 1);
        expect(reauthFrame.streamId).toBe(0);
        expect(reauthFrame.payload.oneofKind).toBe('header');
        if (reauthFrame.payload.oneofKind !== 'header') throw new Error('Expected reauth header frame');
        expect(reauthFrame.payload.header.operation).toBe('reauth');
        const authorization = reauthFrame.payload.header.headers.Authorization;
        expect(authorization).toBeDefined();
        expect(authorization?.value).toEqual(['Bearer char-token']);

        await channel.onMessage(createAuthOkBuffer());
        await charAuth;
    });

    it('rejects a superseded pending auth handshake before starting a new one', async () => {
        const sentFrames: GrpcFrame[] = [];
        const webSocket = {
            data: ref<ArrayBuffer | null>(null),
            status: ref<WebSocketStatus>('OPEN'),
            send: vi.fn().mockImplementation((payload: ArrayBuffer) => {
                sentFrames.push(GrpcFrame.fromBinary(new Uint8Array(payload)));
                return true;
            }),
            open: vi.fn(),
        };

        let token: string | null = null;
        const channel = new WebsocketChannelImpl(createLogger(), webSocket, () => token);
        const onEnd = vi.fn();
        const stream = channel.getStream({
            debug: false,
            methodDefinition: {
                service: { typeName: 'test.Service' },
                name: 'TestMethod',
                serverStreaming: true,
                clientStreaming: false,
            } as never,
            url: '',
            onChunk: vi.fn(),
            onEnd,
            onHeaders: vi.fn(),
        });

        const startPromise = stream.start(new Metadata());
        expect(sentFrames).toHaveLength(1);
        const authFrame = expectFrame(sentFrames, 0);
        expect(authFrame.streamId).toBe(0);
        expect(authFrame.payload.oneofKind).toBe('header');
        if (authFrame.payload.oneofKind !== 'header') throw new Error('Expected auth header frame');
        expect(authFrame.payload.header.operation).toBe('auth');

        token = 'char-token';
        const authPromise = channel.ensureAuthenticated();
        expect(sentFrames).toHaveLength(2);
        const reauthFrame = expectFrame(sentFrames, 1);
        expect(reauthFrame.streamId).toBe(0);
        expect(reauthFrame.payload.oneofKind).toBe('header');
        if (reauthFrame.payload.oneofKind !== 'header') throw new Error('Expected reauth header frame');
        expect(reauthFrame.payload.header.operation).toBe('reauth');

        await channel.onMessage(createAuthOkBuffer());
        await Promise.resolve();

        await startPromise;
        expect(onEnd).toHaveBeenCalledTimes(1);
        expect(Array.from({ length: 6 }, () => channel.getNextStreamId())).toEqual([2, 3, 4, 5, 6, 7]);
        expect(channel.getNextStreamId()).toBe(1);

        await channel.onMessage(createAuthOkBuffer());
        await authPromise;
    });

    it('rejects a pending auth handshake when the websocket closes', async () => {
        const sentFrames: GrpcFrame[] = [];
        const webSocket = {
            data: ref<ArrayBuffer | null>(null),
            status: ref<WebSocketStatus>('OPEN'),
            send: vi.fn().mockImplementation((payload: ArrayBuffer) => {
                sentFrames.push(GrpcFrame.fromBinary(new Uint8Array(payload)));
                return true;
            }),
            open: vi.fn(),
        };

        const channel = new WebsocketChannelImpl(createLogger(), webSocket, () => null);
        const stream = channel.getStream({
            debug: false,
            methodDefinition: {
                service: { typeName: 'test.Service' },
                name: 'TestMethod',
                serverStreaming: true,
                clientStreaming: false,
            } as never,
            url: '',
            onChunk: vi.fn(),
            onEnd: vi.fn(),
            onHeaders: vi.fn(),
        });

        const startPromise = stream.start(new Metadata());
        expect(sentFrames).toHaveLength(1);
        expect(sentFrames[0]?.payload.oneofKind).toBe('header');

        webSocket.status.value = 'CLOSED';
        await new Promise((resolve) => setTimeout(resolve, 300));

        await expect(startPromise).resolves.toBeUndefined();
        expect(channel.getNextStreamId()).toBe(1);
    });

    it('clears ignored auth replies immediately when the websocket reconnects', async () => {
        const sentFrames: GrpcFrame[] = [];
        const webSocket = {
            data: ref<ArrayBuffer | null>(null),
            status: ref<WebSocketStatus>('OPEN'),
            send: vi.fn().mockImplementation((payload: ArrayBuffer) => {
                sentFrames.push(GrpcFrame.fromBinary(new Uint8Array(payload)));
                return true;
            }),
            open: vi.fn(),
        };

        let token: string | null = null;
        const channel = new WebsocketChannelImpl(createLogger(), webSocket, () => token);

        const firstAuth = channel.ensureAuthenticated();
        firstAuth.catch(() => undefined);
        expect(sentFrames).toHaveLength(1);

        token = 'char-token';
        const supersedingAuth = channel.ensureAuthenticated();
        supersedingAuth.catch(() => undefined);
        expect(sentFrames).toHaveLength(2);

        webSocket.status.value = 'CLOSED';
        await Promise.resolve();
        webSocket.status.value = 'OPEN';

        const reconnectAuth = channel.ensureAuthenticated();
        expect(sentFrames).toHaveLength(3);

        await channel.onMessage(createAuthOkBuffer());
        await expect(reconnectAuth).resolves.toBeUndefined();

        await expect(firstAuth).rejects.toThrow('Superseded websocket auth handshake');
        await expect(supersedingAuth).rejects.toThrow('WebSocket closed');
    });

    it('marks a stream closed before auth completes so it does not send its RPC header', async () => {
        const sentFrames: GrpcFrame[] = [];
        const webSocket = {
            data: ref<ArrayBuffer | null>(null),
            status: ref<WebSocketStatus>('OPEN'),
            send: vi.fn().mockImplementation((payload: ArrayBuffer) => {
                sentFrames.push(GrpcFrame.fromBinary(new Uint8Array(payload)));
                return true;
            }),
            open: vi.fn(),
        };

        const channel = new WebsocketChannelImpl(createLogger(), webSocket, () => null);
        const onEnd = vi.fn();
        const stream = channel.getStream({
            debug: false,
            methodDefinition: {
                service: { typeName: 'test.Service' },
                name: 'TestMethod',
                serverStreaming: true,
                clientStreaming: false,
            } as never,
            url: '',
            onChunk: vi.fn(),
            onEnd,
            onHeaders: vi.fn(),
        });

        const startPromise = stream.start(new Metadata());
        expect(sentFrames).toHaveLength(1);
        expect(sentFrames[0]?.streamId).toBe(0);

        const cancelPromise = stream.cancel();
        expect(onEnd).toHaveBeenCalledTimes(1);

        await channel.onMessage(createAuthOkBuffer());
        await expect(startPromise).resolves.toBeUndefined();
        await expect(cancelPromise).resolves.toBeUndefined();

        expect(sentFrames).toHaveLength(1);
        expect(channel.activeStreams.has(1)).toBe(false);
    });

    it('releases a cancelled stream before its failure frame arrives', async () => {
        const sentFrames: GrpcFrame[] = [];
        const webSocket = {
            data: ref<ArrayBuffer | null>(null),
            status: ref<WebSocketStatus>('OPEN'),
            send: vi.fn().mockImplementation((payload: ArrayBuffer) => {
                sentFrames.push(GrpcFrame.fromBinary(new Uint8Array(payload)));
                return true;
            }),
            open: vi.fn(),
        };

        const channel = new WebsocketChannelImpl(createLogger(), webSocket, () => null);
        const onEnd = vi.fn();
        const stream = channel.getStream({
            debug: false,
            methodDefinition: {
                service: { typeName: 'test.Service' },
                name: 'TestMethod',
                serverStreaming: true,
                clientStreaming: false,
            } as never,
            url: '',
            onChunk: vi.fn(),
            onEnd,
            onHeaders: vi.fn(),
        });

        const startPromise = stream.start(new Metadata());
        await channel.onMessage(createAuthOkBuffer());
        await startPromise;

        await stream.cancel();
        expect(onEnd).toHaveBeenCalledTimes(1);
        expect(channel.activeStreams.has(1)).toBe(true);

        await channel.onMessage(createFailureBuffer(1, 'context canceled'));
        expect(onEnd).toHaveBeenCalledTimes(1);
        expect(channel.activeStreams.has(1)).toBe(false);
        expect(Array.from({ length: 6 }, () => channel.getNextStreamId())).toEqual([2, 3, 4, 5, 6, 7]);
        expect(channel.getNextStreamId()).toBe(1);
        expect(sentFrames.some((frame) => frame.payload.oneofKind === 'cancel')).toBe(true);
    });

    it('does not buffer stream frames across a websocket reconnect', async () => {
        const sentFrames: GrpcFrame[] = [];
        const bufferedFrames: ArrayBuffer[] = [];
        const status = ref<WebSocketStatus>('OPEN');
        const webSocket = {
            data: ref<ArrayBuffer | null>(null),
            status,
            send: vi.fn().mockImplementation((payload: ArrayBuffer, useBuffer = true) => {
                if (status.value !== 'OPEN') {
                    if (useBuffer) bufferedFrames.push(payload);
                    return false;
                }

                while (bufferedFrames.length > 0) {
                    const buffered = bufferedFrames.shift();
                    if (buffered) {
                        sentFrames.push(GrpcFrame.fromBinary(new Uint8Array(buffered)));
                    }
                }

                sentFrames.push(GrpcFrame.fromBinary(new Uint8Array(payload)));
                return true;
            }),
            open: vi.fn(),
        };

        watch(status, (value) => {
            if (value !== 'OPEN') return;

            while (bufferedFrames.length > 0) {
                const buffered = bufferedFrames.shift();
                if (buffered) {
                    sentFrames.push(GrpcFrame.fromBinary(new Uint8Array(buffered)));
                }
            }
        });

        const token: string | null = null;
        const channel = new WebsocketChannelImpl(createLogger(), webSocket, () => token);
        const stream = channel.getStream({
            debug: false,
            methodDefinition: {
                service: { typeName: 'test.Service' },
                name: 'TestMethod',
                serverStreaming: true,
                clientStreaming: false,
            } as never,
            url: '',
            onChunk: vi.fn(),
            onEnd: vi.fn(),
            onHeaders: vi.fn(),
        });

        const startPromise = stream.start(new Metadata());
        expect(sentFrames).toHaveLength(1);

        await channel.onMessage(createAuthOkBuffer());
        await startPromise;
        expect(sentFrames).toHaveLength(2);

        status.value = 'CLOSED';

        await expect(stream.sendMessage(new Uint8Array([1, 2, 3]), true)).rejects.toThrow('WebSocket not open');
        expect(bufferedFrames).toHaveLength(0);

        status.value = 'OPEN';
        expect(sentFrames).toHaveLength(2);
    });

    it('keeps the tracked token when an auth_ok arrives outside a pending handshake', async () => {
        const sentFrames: GrpcFrame[] = [];
        const webSocket = {
            data: ref<ArrayBuffer | null>(null),
            status: ref<WebSocketStatus>('OPEN'),
            send: vi.fn().mockImplementation((payload: ArrayBuffer) => {
                sentFrames.push(GrpcFrame.fromBinary(new Uint8Array(payload)));
                return true;
            }),
            open: vi.fn(),
        };

        let token: string | null = 'char-token';
        const channel = new WebsocketChannelImpl(createLogger(), webSocket, () => token);

        const authPromise = channel.ensureAuthenticated();
        expect(sentFrames).toHaveLength(1);
        await channel.onMessage(createAuthOkBuffer());
        await authPromise;

        await channel.onMessage(createAuthOkBuffer());

        // The tracked token should remain the character token, not get clobbered to null.
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        expect((channel as any).authState).toMatchObject({ kind: 'ok', token: 'char-token' });

        token = null;
        const accountAuth = channel.ensureAuthenticated();
        expect(sentFrames).toHaveLength(2);
        expect(sentFrames[1]?.payload.oneofKind).toBe('header');
        if (sentFrames[1]?.payload.oneofKind !== 'header') throw new Error('Expected auth header frame');
        expect(sentFrames[1].payload.header.headers.Authorization).toBeUndefined();

        await channel.onMessage(createAuthOkBuffer());
        await accountAuth;
    });

    it('clears ignored auth replies after reconnecting from a superseded handshake', async () => {
        const sentFrames: GrpcFrame[] = [];
        const webSocket = {
            data: ref<ArrayBuffer | null>(null),
            status: ref<WebSocketStatus>('OPEN'),
            send: vi.fn().mockImplementation((payload: ArrayBuffer) => {
                sentFrames.push(GrpcFrame.fromBinary(new Uint8Array(payload)));
                return true;
            }),
            open: vi.fn(),
        };

        let token: string | null = null;
        const channel = new WebsocketChannelImpl(createLogger(), webSocket, () => token);

        const firstAuth = channel.ensureAuthenticated();
        firstAuth.catch(() => undefined);
        expect(sentFrames).toHaveLength(1);

        token = 'char-token';
        const supersedingAuth = channel.ensureAuthenticated();
        supersedingAuth.catch(() => undefined);
        expect(sentFrames).toHaveLength(2);

        webSocket.status.value = 'CLOSED';
        await new Promise((resolve) => setTimeout(resolve, 300));

        await expect(firstAuth).rejects.toThrow('Superseded websocket auth handshake');
        await expect(supersedingAuth).rejects.toThrow('WebSocket closed');

        webSocket.status.value = 'OPEN';
        const reconnectAuth = channel.ensureAuthenticated();
        expect(sentFrames).toHaveLength(3);
        const reconnectFrame = expectFrame(sentFrames, 2);
        expect(reconnectFrame.streamId).toBe(0);
        expect(reconnectFrame.payload.oneofKind).toBe('header');
        if (reconnectFrame.payload.oneofKind !== 'header') throw new Error('Expected reconnect auth header frame');
        expect(reconnectFrame.payload.header.operation).toBe('auth');
        expect(reconnectFrame.payload.header.headers.Authorization?.value).toEqual(['Bearer char-token']);

        await channel.onMessage(createAuthOkBuffer());
        await expect(reconnectAuth).resolves.toBeUndefined();
    });

    it('waits for websocket auth before sending the first stream body frame', async () => {
        const sentFrames: GrpcFrame[] = [];
        const webSocket = {
            data: ref<ArrayBuffer | null>(null),
            status: ref<WebSocketStatus>('OPEN'),
            send: vi.fn().mockImplementation((payload: ArrayBuffer) => {
                sentFrames.push(GrpcFrame.fromBinary(new Uint8Array(payload)));
                return true;
            }),
            open: vi.fn(),
        };

        const channel = new WebsocketChannelImpl(createLogger(), webSocket, () => null);
        const stream = channel.getStream({
            debug: false,
            methodDefinition: {
                service: { typeName: 'test.Service' },
                name: 'TestMethod',
                serverStreaming: true,
                clientStreaming: false,
            } as never,
            url: '',
            onChunk: vi.fn(),
            onEnd: vi.fn(),
            onHeaders: vi.fn(),
        });

        const startPromise = stream.start(new Metadata());
        expect(sentFrames).toHaveLength(1);
        const authFrame = expectFrame(sentFrames, 0);
        expect(authFrame.streamId).toBe(0);
        expect(authFrame.payload.oneofKind).toBe('header');
        if (authFrame.payload.oneofKind !== 'header') throw new Error('Expected auth header frame');
        expect(authFrame.payload.header.operation).toBe('auth');

        const sendPromise = stream.sendMessage(new Uint8Array([1, 2, 3]), true);
        expect(sentFrames).toHaveLength(1);

        await channel.onMessage(createAuthOkBuffer());
        await startPromise;
        await sendPromise;

        expect(sentFrames).toHaveLength(3);
        const streamHeaderFrame = expectFrame(sentFrames, 1);
        expect(streamHeaderFrame.streamId).toBe(1);
        expect(streamHeaderFrame.payload.oneofKind).toBe('header');
        if (streamHeaderFrame.payload.oneofKind !== 'header') throw new Error('Expected stream header frame');
        expect(streamHeaderFrame.payload.header.operation).toBe('test.Service/TestMethod');

        const bodyFrame = expectFrame(sentFrames, 2);
        expect(bodyFrame.streamId).toBe(1);
        expect(bodyFrame.payload.oneofKind).toBe('body');
    });
});
