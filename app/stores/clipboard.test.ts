import { createPinia, setActivePinia } from 'pinia';
import { beforeEach, describe, expect, it } from 'vitest';
import { useClipboardStore } from './clipboard';

describe('useClipboardStore account scoping', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    it('restores private data after logging out and back into the same account', () => {
        const store = useClipboardStore();

        store.setAccountScope(123);
        store.addUser({ userId: 1, firstname: 'Ada', lastname: 'Lovelace' } as never);
        store.setAccountScope(null);

        expect(store.users).toHaveLength(0);

        store.setAccountScope(123);

        expect(store.users).toHaveLength(1);
        expect(store.users[0]?.userId).toBe(1);
    });

    it("does not expose one account's data to another account", () => {
        const store = useClipboardStore();

        store.setAccountScope(123);
        store.addUser({ userId: 1, firstname: 'Ada', lastname: 'Lovelace' } as never);
        store.setAccountScope(456);

        expect(store.users).toHaveLength(0);
        expect(Object.keys(store.accountData)).toEqual([]);
    });

    it('does not reset data when activating the already-active account', () => {
        const store = useClipboardStore();

        store.setAccountScope(123);
        store.addUser({ userId: 1, firstname: 'Ada', lastname: 'Lovelace' } as never);
        store.setAccountScope(123);

        expect(store.users).toHaveLength(1);
    });

    it('keeps active selections first when promoting clipboard data', () => {
        const store = useClipboardStore();

        store.setAccountScope(123);
        store.addUser({ userId: 1, firstname: 'Ada', lastname: 'Lovelace' } as never);
        store.addUser({ userId: 2, firstname: 'Grace', lastname: 'Hopper' } as never);
        store.activeStack.users = [store.users.find((user) => user.userId === 1)!];

        store.promoteToActiveStack('citizens');

        expect(store.activeStack.users.map((user) => user.userId)).toEqual([1, 2]);
    });

    it('applies max after active selections when promoting clipboard data', () => {
        const store = useClipboardStore();

        store.setAccountScope(123);
        store.addUser({ userId: 1, firstname: 'Ada', lastname: 'Lovelace' } as never);
        store.addUser({ userId: 2, firstname: 'Grace', lastname: 'Hopper' } as never);
        store.activeStack.users = [store.users.find((user) => user.userId === 1)!];

        store.promoteToActiveStack('citizens', 1);

        expect(store.activeStack.users.map((user) => user.userId)).toEqual([1]);
    });
});
