import type { WebSocketStatus } from '@vueuse/core';
import { describe, expect, it, vi } from 'vitest';
import { ref } from 'vue';
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
});
