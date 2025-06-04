import { defineStore } from 'pinia';
import type { Version } from '~/types/history';

const MAX_HISTORY_LENGTH = 20; // Customize this limit as needed

export const useHistoryStore = defineStore('historyStore', () => {
    const history = ref<Version<unknown, unknown>[]>([]);

    // Combine type and id for uniqueness
    const storageKey = (type: string, id: number) => `history_${type}_${id.toString()}`;

    const loadHistory = <TContent, TMeta>(type: string, id: number) => {
        const data = localStorage.getItem(storageKey(type, id));
        history.value = data ? (JSON.parse(data) as Version<TContent, TMeta>[]) : [];
    };

    const saveHistory = <TContent, TMeta>(type: string, id: number) => {
        localStorage.setItem(storageKey(type, id), JSON.stringify(history.value));
    };

    const addVersion = <TContent, TMeta>(type: string, id: number, content: TContent, meta: TMeta, name = '') => {
        const version: Version<TContent, TMeta> = {
            id: new Date().toISOString(),
            type,
            content,
            meta,
            name,
        };
        history.value.push(version as Version<unknown, unknown>);

        // Enforce history length limit
        clearOldVersions(MAX_HISTORY_LENGTH);

        saveHistory(type, id);
    };

    const revertToVersion = <TContent, TMeta>(versionId: string): Version<TContent, TMeta> | undefined => {
        return history.value.find((v) => v.id === versionId) as Version<TContent, TMeta> | undefined;
    };

    const clearOldVersions = (maxVersions: number) => {
        while (history.value.length > maxVersions) {
            history.value.shift();
        }
    };

    const getVersionsByType = <TContent, TMeta>(type: string) =>
        computed(() => history.value.filter((v) => v.type === type) as Version<TContent, TMeta>[]);

    return {
        // State
        history,

        // Actions
        loadHistory,
        saveHistory,
        addVersion,
        revertToVersion,
        clearOldVersions,
        getVersionsByType,
    };
});
