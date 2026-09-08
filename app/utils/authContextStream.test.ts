import { describe, expect, it, vi } from 'vitest';
import { createAuthContextRestartHandler, restartForAuthContext } from './authContextStream';

describe('restartForAuthContext', () => {
    it('aborts and starts an active stream', async () => {
        const stopStream = vi.fn().mockResolvedValue(undefined);
        const startStream = vi.fn().mockResolvedValue(undefined);

        await restartForAuthContext({ abort: new AbortController(), isStopping: () => false }, stopStream, startStream);

        expect(stopStream).toHaveBeenCalledOnce();
        expect(startStream).toHaveBeenCalledOnce();
    });

    it('does not restart an inactive or intentionally stopped stream', async () => {
        const stopStream = vi.fn().mockResolvedValue(undefined);
        const startStream = vi.fn().mockResolvedValue(undefined);

        await restartForAuthContext({ isStopping: () => false }, stopStream, startStream);
        await restartForAuthContext({ abort: new AbortController(), isStopping: () => true }, stopStream, startStream);

        expect(stopStream).not.toHaveBeenCalled();
        expect(startStream).not.toHaveBeenCalled();
    });

    it('does not start again if stopping begins while the old stream is stopping', async () => {
        const state = { abort: new AbortController(), stopping: false };
        const stopStream = vi.fn().mockImplementation(async () => {
            state.stopping = true;
        });
        const startStream = vi.fn().mockResolvedValue(undefined);

        await restartForAuthContext({ abort: state.abort, isStopping: () => state.stopping }, stopStream, startStream);

        expect(stopStream).toHaveBeenCalledOnce();
        expect(startStream).not.toHaveBeenCalled();
    });
});

describe('createAuthContextRestartHandler', () => {
    it('waits for the transition to finish after observing a new key', () => {
        const restart = vi.fn();
        const handle = createAuthContextRestartHandler('key-1', restart);

        handle(['key-2', true]);
        expect(restart).not.toHaveBeenCalled();

        handle(['key-2', false]);
        expect(restart).toHaveBeenCalledOnce();
    });

    it('does not restart repeatedly for the same committed key', () => {
        const restart = vi.fn();
        const handle = createAuthContextRestartHandler('key-1', restart);

        handle(['key-2', false]);
        handle(['key-2', false]);
        handle(['key-2', true]);

        expect(restart).toHaveBeenCalledOnce();
    });
});
