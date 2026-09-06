import { createPinia, setActivePinia } from 'pinia';
import { beforeEach, describe, expect, it } from 'vitest';
import { useHistoryStore } from './history';

describe('useHistoryStore account scoping', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    it('restores history after logging out and back into the same account', () => {
        const store = useHistoryStore();

        store.setAccountScope(123);
        store.addVersion('document', 1, { title: 'Private' });
        store.setAccountScope(null);

        expect(store.history).toHaveLength(0);

        store.setAccountScope(123);

        expect(store.history).toHaveLength(1);
        expect(store.history[0]?.content).toEqual({ title: 'Private' });
    });

    it("does not expose one account's history to another account", () => {
        const store = useHistoryStore();

        store.setAccountScope(123);
        store.addVersion('document', 1, { title: 'Private' });
        store.setAccountScope(456);

        expect(store.history).toHaveLength(0);
        expect(Object.keys(store.accountData)).toEqual([]);
    });
});
