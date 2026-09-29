import { describe, expect, it } from 'vitest';
import { nextTick, ref } from 'vue';
import * as Y from 'yjs';
import { useYArrayFiltered } from './useY';
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
