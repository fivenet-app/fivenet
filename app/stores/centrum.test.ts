import { createPinia, setActivePinia } from 'pinia';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches/dispatches';
import { CentrumMode } from '~~/gen/ts/resources/centrum/settings/settings';
import { StatusUnit } from '~~/gen/ts/resources/centrum/units/units';
import { toTimestamp } from '~/utils/time';
import { useCentrumStore } from './centrum';

const mocks = vi.hoisted(() => ({
    notifications: {
        add: vi.fn(),
    },
}));

vi.mock('~/stores/notifications', () => ({
    useNotificationsStore: () => mocks.notifications,
}));

vi.mock('~/composables/useSounds', () => ({
    useSounds: () => ({ play: vi.fn() }),
}));

describe('useCentrumStore', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
        vi.clearAllMocks();
    });

    it('keeps feed events with matching unit and dispatch status IDs', () => {
        const store = useCentrumStore();

        store.addFeedItem({
            id: 17,
            unitId: 3,
            status: StatusUnit.AVAILABLE,
        });
        store.addFeedItem({
            id: 17,
            dispatchId: 4,
            status: StatusDispatch.NEW,
        });

        expect(store.feed).toHaveLength(2);
        expect(store.feed.filter((item) => 'unitId' in item)).toHaveLength(1);
        expect(store.feed.filter((item) => 'dispatchId' in item)).toHaveLength(1);
    });

    it('removes every adjacent pending dispatch absent from the projection', async () => {
        const store = useCentrumStore();
        store.pendingDispatches = [11, 12, 13];

        await store.cleanup();

        expect(store.pendingDispatches).toEqual([]);
    });

    it('clears dispatch revision state when its projection is removed', () => {
        const store = useCentrumStore();
        store.addOrUpdateDispatch({ id: 11, units: [] }, 2);
        store.removeDispatch(11, 3);

        // Removed projections do not retain a tombstone revision indefinitely.
        store.addOrUpdateDispatch({ id: 11, units: [] }, 2);

        expect(store.dispatches.has(11)).toBe(true);
    });

    it('keeps unit and dispatch revision state separate for matching IDs', () => {
        const store = useCentrumStore();
        store.addOrUpdateUnit({ id: 7, users: [] }, 3);
        store.addOrUpdateDispatch({ id: 7, units: [] }, 3);
        store.removeUnit(7, 4);

        expect(store.units.has(7)).toBe(false);
        expect(store.dispatches.has(7)).toBe(true);
    });

    it('keeps a pending dispatch when only another unit assignment expires', async () => {
        const store = useCentrumStore();
        store.ownUnitId = 1;
        store.pendingDispatches = [11];
        store.addOrUpdateDispatch({
            id: 11,
            units: [
                { unitId: 1, expiresAt: toTimestamp(new Date(Date.now() + 60_000)) },
                { unitId: 2, expiresAt: toTimestamp(new Date(Date.now() - 60_000)) },
            ],
        });

        await store.cleanup();

        expect(store.pendingDispatches).toEqual([11]);
    });

    it('clears the previous projection before an authorization-context restart', async () => {
        const store = useCentrumStore();
        store.stopping = true;
        store.settings = { job: 'ambulance' };
        store.acls = { dispatches: { jobs: [{ job: 'ambulance' }] } };
        store.isDispatcher = true;
        store.dispatchers = [{ job: 'ambulance' }];
        store.feed = [{ id: 1, unitId: 2, status: StatusUnit.AVAILABLE }];
        store.units = new Map([[2, { id: 2 }]]);
        store.dispatches = new Map([[3, { id: 3, units: [] }]]);
        store.ownUnitId = 2;
        store.ownDispatches = [3];
        store.pendingDispatches = [3];

        await store.restartForAuthContext();

        expect(store.settings).toBeUndefined();
        expect(store.acls).toBeUndefined();
        expect(store.isDispatcher).toBe(false);
        expect(store.dispatchers).toEqual([]);
        expect(store.feed).toEqual([]);
        expect(store.units.size).toBe(0);
        expect(store.dispatches.size).toBe(0);
        expect(store.ownUnitId).toBeUndefined();
        expect(store.ownDispatches).toEqual([]);
        expect(store.pendingDispatches).toEqual([]);
    });

    it('uses fallback mode until a dispatcher is active', () => {
        const store = useCentrumStore();
        store.settings = {
            mode: CentrumMode.CENTRAL_COMMAND,
            fallbackMode: CentrumMode.SIMPLIFIED,
        };
        store.dispatchers = [{ job: 'ambulance', dispatchers: [] }];

        expect(store.getCurrentMode).toBe(CentrumMode.SIMPLIFIED);

        store.dispatchers[0]!.dispatchers.push({ userId: 7 });

        expect(store.getCurrentMode).toBe(CentrumMode.CENTRAL_COMMAND);
    });
});
