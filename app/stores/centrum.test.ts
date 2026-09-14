import { createPinia, setActivePinia } from 'pinia';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { toTimestamp } from '~/utils/time';
import { type Dispatchers } from '~~/gen/ts/resources/centrum/dispatchers/dispatchers';
import {
    type Dispatch,
    type DispatchAssignment,
    type DispatchStatus,
    StatusDispatch,
} from '~~/gen/ts/resources/centrum/dispatches/dispatches';
import { CentrumMode, type EffectiveAccess, type Settings } from '~~/gen/ts/resources/centrum/settings/settings';
import { StatusUnit, type Unit, type UnitStatus } from '~~/gen/ts/resources/centrum/units/units';
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

const dispatch = (overrides: Partial<Dispatch> = {}): Dispatch => ({
    id: 0,
    job: 'ambulance',
    message: '',
    x: 0,
    y: 0,
    anon: false,
    units: [],
    ...overrides,
});

const unit = (overrides: Partial<Unit> = {}): Unit => ({
    id: 0,
    job: 'ambulance',
    sortOrder: 0,
    name: '',
    initials: '',
    color: '',
    users: [],
    ...overrides,
});

const settings = (overrides: Partial<Settings> = {}): Settings => ({
    job: 'ambulance',
    enabled: true,
    type: 0,
    public: false,
    mode: CentrumMode.UNSPECIFIED,
    fallbackMode: CentrumMode.UNSPECIFIED,
    ...overrides,
});

const access = (overrides: Partial<EffectiveAccess> = {}): EffectiveAccess => ({
    dispatches: { jobs: [] },
    ...overrides,
});

const dispatchers = (overrides: Partial<Dispatchers> = {}): Dispatchers => ({
    job: 'ambulance',
    dispatchers: [],
    ...overrides,
});

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
        store.addOrUpdateDispatch(dispatch({ id: 11 }), 2);
        store.removeDispatch(11, 3);

        // Removed projections do not retain a tombstone revision indefinitely.
        store.addOrUpdateDispatch(dispatch({ id: 11 }), 2);

        expect(store.dispatches.has(11)).toBe(true);
    });

    it('keeps unit and dispatch revision state separate for matching IDs', () => {
        const store = useCentrumStore();
        store.addOrUpdateUnit(unit({ id: 7 }), 3);
        store.addOrUpdateDispatch(dispatch({ id: 7 }), 3);
        store.removeUnit(7, 4);

        expect(store.units.has(7)).toBe(false);
        expect(store.dispatches.has(7)).toBe(true);
    });

    it('keeps a pending dispatch when only another unit assignment expires', async () => {
        const store = useCentrumStore();
        store.ownUnitId = 1;
        store.pendingDispatches = [11];
        store.addOrUpdateDispatch(
            dispatch({
                id: 11,
                units: [
                    { dispatchId: 11, unitId: 1, expiresAt: toTimestamp(new Date(Date.now() + 60_000)) } as DispatchAssignment,
                    { dispatchId: 11, unitId: 2, expiresAt: toTimestamp(new Date(Date.now() - 60_000)) } as DispatchAssignment,
                ],
            }),
        );

        await store.cleanup();

        expect(store.pendingDispatches).toEqual([11]);
    });

    it('clears the previous projection before an authorization-context restart', async () => {
        const store = useCentrumStore();
        store.stopping = true;
        store.settings = settings();
        store.acls = access({ dispatches: { jobs: [{ job: 'ambulance', access: 0 }] } });
        store.isDispatcher = true;
        store.dispatchers = [dispatchers()];
        store.feed = [{ id: 1, unitId: 2, status: StatusUnit.AVAILABLE }];
        store.units = new Map([[2, unit({ id: 2 })]]);
        store.dispatches = new Map([[3, dispatch({ id: 3 })]]);
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
        store.settings = settings({
            mode: CentrumMode.CENTRAL_COMMAND,
            fallbackMode: CentrumMode.SIMPLIFIED,
        });
        store.dispatchers = [dispatchers()];

        expect(store.getCurrentMode).toBe(CentrumMode.SIMPLIFIED);

        store.dispatchers[0]!.dispatchers.push({
            userId: 7,
            job: 'ambulance',
            jobGrade: 0,
            firstname: '',
            lastname: '',
            dateofbirth: '',
        });

        expect(store.getCurrentMode).toBe(CentrumMode.CENTRAL_COMMAND);
    });

    it('moves an assigned dispatch from pending to own after acceptance', () => {
        const store = useCentrumStore();
        store.ownUnitId = 7;

        store.addOrUpdateDispatch(
            dispatch({
                id: 21,
                units: [{ dispatchId: 21, unitId: 7, expiresAt: toTimestamp(new Date(Date.now() + 60_000)) }],
            }),
        );
        expect(store.pendingDispatches).toEqual([21]);
        expect(store.ownDispatches).toEqual([]);

        store.addOrUpdateDispatch(
            dispatch({
                id: 21,
                units: [{ dispatchId: 21, unitId: 7 }],
            }),
        );

        expect(store.pendingDispatches).toEqual([]);
        expect(store.ownDispatches).toEqual([21]);
    });

    it('notifies once when a dispatch first becomes pending', () => {
        const store = useCentrumStore();
        store.ownUnitId = 7;
        const pending = dispatch({
            id: 21,
            units: [{ dispatchId: 21, unitId: 7, expiresAt: toTimestamp(new Date(Date.now() + 60_000)) }],
        });

        store.addOrUpdateDispatch(pending);
        store.addOrUpdateDispatch(dispatch({ id: 21, units: [...pending.units] }));

        expect(store.pendingDispatches).toEqual([21]);
        expect(mocks.notifications.add).toHaveBeenCalledTimes(1);
    });

    it('does not replay an older dispatch status over a newer one', () => {
        const store = useCentrumStore();
        store.addOrUpdateDispatch(dispatch({ id: 21 }));

        store.updateDispatchStatus({ id: 9, dispatchId: 21, status: StatusDispatch.NEW } as DispatchStatus);
        store.updateDispatchStatus({ id: 8, dispatchId: 21, status: StatusDispatch.COMPLETED } as DispatchStatus);

        expect(store.dispatches.get(21)?.status?.status).toBe(StatusDispatch.NEW);
    });

    it('removes only the unassigned unit from a dispatch', () => {
        const store = useCentrumStore();
        store.addOrUpdateDispatch(
            dispatch({
                id: 21,
                units: [
                    { dispatchId: 21, unitId: 7 },
                    { dispatchId: 21, unitId: 8 },
                ],
            }),
        );

        store.updateDispatchStatus({
            id: 3,
            dispatchId: 21,
            unitId: 7,
            status: StatusDispatch.UNIT_UNASSIGNED,
        } as DispatchStatus);

        expect(store.dispatches.get(21)?.units.map((assignment) => assignment.unitId)).toEqual([8]);
    });

    it('does not replay an older unit status over a newer one', () => {
        const store = useCentrumStore();
        store.addOrUpdateUnit(unit({ id: 7 }));

        store.updateUnitStatus({ id: 9, unitId: 7, status: StatusUnit.BUSY } as UnitStatus);
        store.updateUnitStatus({ id: 8, unitId: 7, status: StatusUnit.AVAILABLE } as UnitStatus);

        expect(store.units.get(7)?.status?.status).toBe(StatusUnit.BUSY);
    });

    it('keeps the feed bounded while retaining the newest entries', () => {
        const store = useCentrumStore();

        for (let id = 1; id <= 101; id++) {
            store.addFeedItem({ id, unitId: 7, status: StatusUnit.AVAILABLE });
        }

        expect(store.feed).toHaveLength(100);
        expect(store.feed[0]?.id).toBe(101);
        expect(store.feed.at(-1)?.id).toBe(2);
    });

    it('removes completed dispatches during cleanup after the retention interval', async () => {
        const store = useCentrumStore();
        const completedAt = new Date(Date.now() - 60_000);
        store.addOrUpdateDispatch(
            dispatch({
                id: 21,
                createdAt: toTimestamp(completedAt),
                status: {
                    id: 4,
                    dispatchId: 21,
                    createdAt: toTimestamp(completedAt),
                    status: StatusDispatch.COMPLETED,
                },
            }),
        );

        await store.cleanup();

        expect(store.dispatches.has(21)).toBe(false);
    });
});
