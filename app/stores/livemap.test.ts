import { createPinia, setActivePinia } from 'pinia';
import { beforeEach, describe, expect, it } from 'vitest';
import type { UserMarker } from '~~/gen/ts/resources/livemap/markers/user_marker';
import { useLivemapStore } from './livemap';

const marker = (overrides: Partial<UserMarker> = {}): UserMarker =>
    ({
        userId: 1,
        x: 10,
        y: 20,
        job: 'police',
        jobLabel: 'Police',
        hidden: false,
        ...overrides,
    }) as UserMarker;

describe('useLivemapStore user markers', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    it('switches directly to a different follow target', () => {
        const store = useLivemapStore();
        const first = marker({ userId: 1 });
        const second = marker({ userId: 2 });

        store.toggleFollowMarker(first);
        store.toggleFollowMarker(second);

        expect(store.followMarker).toBe(true);
        expect(store.selectedMarker).toMatchObject(second);
    });

    it('toggles off when the current follow target is selected again', () => {
        const store = useLivemapStore();
        const selected = marker();

        store.toggleFollowMarker(selected);
        store.toggleFollowMarker(selected);

        expect(store.followMarker).toBe(false);
        expect(store.selectedMarker).toBeUndefined();
    });

    it('keeps selected markers attached to the cached snapshot object', () => {
        const store = useLivemapStore();
        const initial = marker({ userId: 7 });
        const refreshed = marker({ userId: 7, x: 30, y: 40 });

        store.applyUserSnapshot({ markers: [initial] }, null);
        store.toggleFollowMarker(initial);
        store.applyUserSnapshot({ markers: [refreshed] }, null);

        expect(store.selectedMarker).toBe(store.markersUsers.get(7));
        expect(store.selectedMarker).toMatchObject({ x: 30, y: 40 });
        expect(store.followMarker).toBe(true);
    });

    it('refreshes cached user details for an existing marker', () => {
        const store = useLivemapStore();
        const initial = marker({
            user: { firstname: 'Old', lastname: 'Name' } as NonNullable<UserMarker['user']>,
        });
        const updated = marker({
            user: { firstname: 'New', lastname: 'Name' } as NonNullable<UserMarker['user']>,
        });

        store.addOrUpdateUserMarker(initial);
        store.addOrUpdateUserMarker(updated);

        expect(store.markersUsers.get(1)?.user?.firstname).toBe('New');
    });

    it('keeps ownMarker attached to the cached marker across updates', () => {
        const store = useLivemapStore();
        const initial = marker({ userId: 9 });
        const updated = marker({ userId: 9, x: 50, y: 60 });

        store.addOrUpdateUserMarker(initial, 9);
        store.addOrUpdateUserMarker(updated, 9);

        expect(store.ownMarker).toBe(store.markersUsers.get(9));
        expect(store.ownMarker).toMatchObject({ x: 50, y: 60 });
    });
});
