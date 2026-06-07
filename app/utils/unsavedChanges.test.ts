import { describe, expect, it } from 'vitest';
import { runSharedUnsavedChangesConfirmation } from './unsavedChanges';

function createDeferred<T>() {
    let resolve!: (value: T) => void;
    const promise = new Promise<T>((res) => {
        resolve = res;
    });

    return { promise, resolve };
}

describe('runSharedUnsavedChangesConfirmation', () => {
    it('reuses an in-flight confirmation for the same group', async () => {
        const confirmationGroup = Symbol('shared-group');
        const deferred = createDeferred<boolean>();
        let confirmCalls = 0;

        const firstConfirmation = runSharedUnsavedChangesConfirmation(confirmationGroup, async () => {
            confirmCalls++;
            return deferred.promise;
        });
        const secondConfirmation = runSharedUnsavedChangesConfirmation(confirmationGroup, async () => {
            confirmCalls++;
            return true;
        });

        expect(secondConfirmation).toBe(firstConfirmation);

        deferred.resolve(true);

        await expect(firstConfirmation).resolves.toBe(true);
        await expect(secondConfirmation).resolves.toBe(true);
        expect(confirmCalls).toBe(1);
    });

    it('allows a new confirmation after the previous one settled', async () => {
        const confirmationGroup = Symbol('shared-group');
        let confirmCalls = 0;

        const firstConfirmation = runSharedUnsavedChangesConfirmation(confirmationGroup, async () => {
            confirmCalls++;
            return true;
        });

        await expect(firstConfirmation).resolves.toBe(true);

        const secondConfirmation = runSharedUnsavedChangesConfirmation(confirmationGroup, async () => {
            confirmCalls++;
            return false;
        });

        await expect(secondConfirmation).resolves.toBe(false);
        expect(confirmCalls).toBe(2);
    });
});
