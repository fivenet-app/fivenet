type StreamLike<T> = {
    responses: AsyncIterable<T>;
};

export type ReconnectingStreamState = 'idle' | 'connecting' | 'connected' | 'reconnecting' | 'stopped' | 'error';

export type ReconnectingServerStreamOptions<T> = {
    signal: AbortSignal;
    create: (signal: AbortSignal) => StreamLike<T>;
    initialDelayMs?: number;
    maxDelayMs?: number;
    jitterRatio?: number;
    shouldRetry?: (error: unknown) => boolean;
    onState?: (state: ReconnectingStreamState) => void;
    onReady?: () => void;
    onRetry?: (error: unknown, delayMs: number) => void;
};

/** Waits between automatic stream attempts, unless the owner is stopping. */
function waitForRetry(delayMs: number, signal: AbortSignal): Promise<boolean> {
    if (signal.aborted) return Promise.resolve(false);

    return new Promise<boolean>((resolve) => {
        const onAbort = () => {
            clearTimeout(timer);
            resolve(false);
        };
        const timer = setTimeout(() => {
            signal.removeEventListener('abort', onAbort);
            resolve(true);
        }, delayMs);

        signal.addEventListener('abort', onAbort, { once: true });
    });
}

/** Recreates a server stream after transport failures until the caller aborts. */
export async function* reconnectingServerStream<T>(
    options: ReconnectingServerStreamOptions<T>,
): AsyncGenerator<T, void, undefined> {
    const initialDelayMs = options.initialDelayMs ?? 1_000;
    const maxDelayMs = options.maxDelayMs ?? 15_000;
    const jitterRatio = options.jitterRatio ?? 0.2;
    let delayMs = initialDelayMs;
    let attempted = false;

    /** Applies retry policy and waits before the next stream attempt. */
    const scheduleRetry = async (error: unknown): Promise<boolean> => {
        if (error !== undefined && options.shouldRetry && !options.shouldRetry(error)) {
            options.onState?.('error');
            throw error;
        }

        const jitter = delayMs * jitterRatio * (Math.random() * 2 - 1);
        const retryDelayMs = Math.max(0, Math.round(delayMs + jitter));
        options.onState?.('reconnecting');
        options.onRetry?.(error, retryDelayMs);
        const shouldContinue = await waitForRetry(retryDelayMs, options.signal);
        if (shouldContinue) delayMs = Math.min(delayMs * 2, maxDelayMs);
        return shouldContinue;
    };

    while (!options.signal.aborted) {
        options.onState?.(attempted ? 'reconnecting' : 'connecting');
        attempted = true;
        let receivedResponse = false;
        const streamAbort = new AbortController();
        const abortStream = () => streamAbort.abort();
        options.signal.addEventListener('abort', abortStream, { once: true });
        let iterator: AsyncIterator<T> | undefined;

        try {
            try {
                const stream = options.create(streamAbort.signal);
                iterator = stream.responses[Symbol.asyncIterator]();
            } catch (error) {
                if (options.signal.aborted) return;
                if (!(await scheduleRetry(error))) return;
                continue;
            }

            while (!options.signal.aborted) {
                let result: IteratorResult<T>;
                try {
                    result = await iterator.next();
                } catch (error) {
                    if (options.signal.aborted) return;
                    if (!(await scheduleRetry(error))) return;
                    break;
                }

                if (result.done) {
                    if (options.signal.aborted) return;
                    if (!(await scheduleRetry(undefined))) return;
                    break;
                }

                const response = result.value;
                if (!receivedResponse) {
                    receivedResponse = true;
                    delayMs = initialDelayMs;
                    options.onState?.('connected');
                    options.onReady?.();
                }

                yield response;
            }
        } finally {
            options.signal.removeEventListener('abort', abortStream);
            streamAbort.abort();
            await iterator?.return?.();
        }
    }
}
