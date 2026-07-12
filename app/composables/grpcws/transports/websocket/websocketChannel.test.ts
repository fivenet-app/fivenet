import type { WebSocketStatus } from '@vueuse/core';
import { describe, expect, it, vi } from 'vitest';
import { ref } from 'vue';
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
            send: vi.fn().mockImplementation(async (payload: ArrayBuffer) => {
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

    it('waits for websocket auth before sending the first stream body frame', async () => {
        const sentFrames: GrpcFrame[] = [];
        const webSocket = {
            data: ref<ArrayBuffer | null>(null),
            status: ref<WebSocketStatus>('OPEN'),
            send: vi.fn().mockImplementation(async (payload: ArrayBuffer) => {
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
