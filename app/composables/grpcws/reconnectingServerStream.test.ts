import { describe, expect, it, vi } from 'vitest';
import { reconnectingServerStream } from './reconnectingServerStream';

async function* responses<T>(value: T): AsyncGenerator<T> {
    yield value;
}

function throwingResponses(error: unknown): AsyncIterable<never> {
    return {
        [Symbol.asyncIterator]: () => ({
            async next(): Promise<IteratorResult<never>> {
                throw error;
            },
        }),
    };
}

describe('reconnectingServerStream', () => {
    it('retries with capped exponential backoff', async () => {
        vi.useFakeTimers();

        try {
            let attempts = 0;
            const retryDelays: number[] = [];
            const controller = new AbortController();
            const stream = reconnectingServerStream({
                signal: controller.signal,
                initialDelayMs: 10,
                maxDelayMs: 25,
                jitterRatio: 0,
                create: () => {
                    attempts++;
                    return {
                        responses: attempts < 4 ? throwingResponses(new Error('temporary failure')) : responses('connected'),
                    };
                },
                onRetry: (_, delayMs) => retryDelays.push(delayMs),
            });

            const result = stream.next();
            await vi.advanceTimersByTimeAsync(10);
            await vi.advanceTimersByTimeAsync(20);
            await vi.advanceTimersByTimeAsync(25);

            await expect(result).resolves.toEqual({ value: 'connected', done: false });
            expect(attempts).toBe(4);
            expect(retryDelays).toEqual([10, 20, 25]);

            controller.abort();
            await stream.return(undefined);
        } finally {
            vi.useRealTimers();
        }
    });

    it('resets the backoff after receiving a response', async () => {
        vi.useFakeTimers();

        try {
            const retryDelays: number[] = [];
            const controller = new AbortController();
            let attempts = 0;
            const stream = reconnectingServerStream({
                signal: controller.signal,
                initialDelayMs: 10,
                maxDelayMs: 100,
                jitterRatio: 0,
                create: () => {
                    attempts++;
                    return {
                        responses: attempts === 1 ? throwingResponses(new Error('temporary failure')) : responses('connected'),
                    };
                },
                onRetry: (_, delayMs) => retryDelays.push(delayMs),
            });

            const first = stream.next();
            await vi.advanceTimersByTimeAsync(10);
            await expect(first).resolves.toEqual({ value: 'connected', done: false });

            const second = stream.next();
            await vi.advanceTimersByTimeAsync(10);
            await expect(second).resolves.toEqual({ value: 'connected', done: false });
            expect(retryDelays).toEqual([10, 10]);

            controller.abort();
            await stream.return(undefined);
        } finally {
            vi.useRealTimers();
        }
    });

    it('stops cleanly when aborted during backoff', async () => {
        vi.useFakeTimers();

        try {
            const controller = new AbortController();
            const stream = reconnectingServerStream({
                signal: controller.signal,
                initialDelayMs: 100,
                jitterRatio: 0,
                create: () => ({
                    responses: throwingResponses(new Error('temporary failure')),
                }),
            });

            const result = stream.next();
            await Promise.resolve();
            controller.abort();

            await expect(result).resolves.toEqual({ value: undefined, done: true });
        } finally {
            vi.useRealTimers();
        }
    });

    it('propagates errors marked as non-retryable', async () => {
        const error = new Error('permanent failure');
        const controller = new AbortController();
        const stream = reconnectingServerStream({
            signal: controller.signal,
            create: () => ({
                responses: throwingResponses(error),
            }),
            shouldRetry: () => false,
        });

        await expect(stream.next()).rejects.toBe(error);
    });
});
