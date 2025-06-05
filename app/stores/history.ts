import { defineStore } from 'pinia';
import type { Version } from '~/types/history';

const MAX_HISTORY_LENGTH = 20; // Customize this limit as needed

export const useHistoryStore = defineStore(
    'historyStore',
    () => {
        const history = ref<Version<unknown>[]>([]);

        // Combine type and id for uniqueness
        const storageKey = (type: string, id: number) => `history_${type}_${id.toString()}`;

        const loadHistory = <TContent>(type: string, id: number) => {
            const data = localStorage.getItem(storageKey(type, id));
            history.value = data ? (JSON.parse(data) as Version<TContent>[]) : [];
        };

        const saveHistory = <_TContent>(type: string, id: number) => {
            localStorage.setItem(storageKey(type, id), JSON.stringify(history.value));
        };

        const addVersion = <TContent>(type: string, id: number, content: TContent, name = '') => {
            const version: Version<TContent> = {
                id: new Date().toISOString(),
                type,
                content,
                name,
            };
            history.value.push(version as Version<unknown>);

            // Enforce history length limit
            clearOldVersions(MAX_HISTORY_LENGTH);

            saveHistory(type, id);
        };

        const revertToVersion = <TContent>(versionId: string): Version<TContent> | undefined => {
            return history.value.find((v) => v.id === versionId) as Version<TContent> | undefined;
        };

        const clearOldVersions = (maxVersions: number) => {
            while (history.value.length > maxVersions) {
                history.value.shift();
            }
        };

        const getVersionsByType = <TContent>(type: string) =>
            computed(() => history.value.filter((v) => v.type === type) as Version<TContent>[]);

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
    },
    {
        persist: true,
    },
);
