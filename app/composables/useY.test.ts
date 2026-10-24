import { describe, expect, it } from 'vitest';
import { nextTick, ref } from 'vue';
import * as Y from 'yjs';
import { useYArray, useYArrayFiltered, useYMap, useYObject, useYStructure } from './useY';
import type { YjsSyncOptions } from './useYPrimitives';

function authoritativeProvider(): YjsSyncOptions['provider'] {
    return {
        isAuthoritative: true,
        isSynced: true,
    } as YjsSyncOptions['provider'];
}

function accessMap(id: number, job: string, access: number): Y.Map<unknown> {
    const map = new Y.Map<unknown>();
    map.set('id', id);
    map.set('job', job);
    map.set('access', access);
    return map;
}

describe('useYArrayFiltered', () => {
    it('does not overwrite an existing remote list when the client is authoritative', () => {
        const doc = new Y.Doc();
        const remoteEntries = doc.getArray('access');
        remoteEntries.push([accessMap(10, 'police', 2)]);

        const entries = ref([{ id: 99, job: 'ambulance', access: 3 }]);

        useYArrayFiltered(remoteEntries, entries, { omit: [], nested: false }, { provider: authoritativeProvider() });

        expect(entries.value).toEqual([{ id: 10, job: 'police', access: 2 }]);
        expect(remoteEntries.toJSON()).toEqual([{ id: 10, job: 'police', access: 2 }]);
    });

    it('applies a remote entry update without rewriting the whole list', async () => {
        const doc = new Y.Doc();
        const remoteEntries = doc.getArray('access');
        remoteEntries.push([accessMap(10, 'police', 2), accessMap(11, 'ambulance', 2)]);

        const entries = ref<Record<string, unknown>[]>([]);
        useYArrayFiltered(remoteEntries, entries, { omit: [], nested: false });

        let updates = 0;
        doc.on('update', () => updates++);

        (remoteEntries.get(1) as Y.Map<unknown>).set('access', 3);
        await nextTick();
        await nextTick();

        expect(entries.value).toEqual([
            { id: 10, job: 'police', access: 2 },
            { id: 11, job: 'ambulance', access: 3 },
        ]);
        expect(updates).toBe(1);
        expect(remoteEntries.toJSON()).toEqual([
            { id: 10, job: 'police', access: 2 },
            { id: 11, job: 'ambulance', access: 3 },
        ]);
    });
});

describe('useYArray', () => {
    it('hydrates primitive arrays from Yjs and writes local changes back', async () => {
        const doc = new Y.Doc();
        const values = doc.getArray('values');
        values.push(['remote']);

        const items = ref(['local']);
        const result = useYArray(values, items);

        expect(result).toBe(items);
        expect(items.value).toEqual(['remote']);

        items.value.push('added');
        await nextTick();

        expect(values.toJSON()).toEqual(['remote', 'added']);
    });

    it('seeds an empty authoritative primitive array', () => {
        const doc = new Y.Doc();
        const values = doc.getArray('values');
        const items = ref(['local', 'values']);

        useYArray(values, items, { provider: authoritativeProvider() });

        expect(values.toJSON()).toEqual(['local', 'values']);
    });

    it('hydrates nested object arrays and syncs nested local changes', async () => {
        const doc = new Y.Doc();
        const records = doc.getArray('records');
        const record = new Y.Map<unknown>();
        record.set('name', 'remote');
        records.push([record]);

        const items = ref([{ name: 'local' }]);
        useYArray(records, items);

        expect(items.value).toEqual([{ name: 'remote' }]);

        items.value[0]!.name = 'updated';
        await nextTick();
        await nextTick();

        expect(records.toJSON()).toEqual([{ name: 'updated' }]);
    });
});

describe('useYMap', () => {
    it('syncs only filtered keys in both directions', async () => {
        const doc = new Y.Doc();
        const map = doc.getMap<unknown>('state');
        map.set('name', 'remote');
        map.set('count', 4);

        const state = ref({ name: 'local', count: 1, ignored: true });
        const result = useYMap(map, state, { only: ['name', 'count'] });

        expect(result).toBe(state.value);
        expect(state.value).toEqual({ name: 'remote', count: 4, ignored: true });

        await nextTick();
        state.value.name = 'updated';
        await nextTick();

        expect(map.toJSON()).toEqual({ name: 'updated', count: 4 });
        expect(map.get('ignored')).toBeUndefined();
    });
});

describe('useYStructure and useYObject', () => {
    it('binds nested structures and respects an omitted field', async () => {
        const doc = new Y.Doc();
        const map = doc.getMap<unknown>('state');
        const remoteProfile = new Y.Map<unknown>();
        remoteProfile.set('displayName', 'Remote User');
        map.set('profile', remoteProfile);

        const state = ref({
            profile: { displayName: 'Local User' },
            hidden: 'local-only',
            enabled: false,
        });
        useYStructure(map, state, { omit: ['hidden'] }, { provider: authoritativeProvider() });

        expect(state.value).toEqual({ profile: { displayName: 'Remote User' }, hidden: 'local-only', enabled: false });

        await nextTick();
        state.value.profile.displayName = 'Updated User';
        await nextTick();

        expect((map.get('profile') as Y.Map<unknown>).toJSON()).toEqual({ displayName: 'Updated User' });
        expect(map.toJSON()).toEqual({ profile: { displayName: 'Updated User' }, enabled: false });
    });

    it('returns the same object ref and binds its fields', async () => {
        const doc = new Y.Doc();
        const map = doc.getMap<unknown>('state');
        const state = ref({ title: 'initial' });

        const result = useYObject(map, state, undefined, { provider: authoritativeProvider() });

        expect(result).toBe(state);
        expect(map.toJSON()).toEqual({ title: 'initial' });

        map.set('title', 'remote');
        await nextTick();

        expect(state.value.title).toBe('remote');
    });
});
