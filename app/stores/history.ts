import { defineStore } from 'pinia';
import type { Version } from '~/types/history';

const MAX_HISTORY_LENGTH = 20; // Customize this limit as needed

export const useHistoryStore = defineStore(
    'historyStore',
    () => {
        const history = ref<Version<unknown>[]>([]);

        const listHistory = <TContent>(type: string) => {
            return computed(() => {
                return history.value.filter((v) => v.type === type) as Version<TContent>[];
            });
        };

        const addVersion = <TContent>(type: string, id: number, content: TContent, name = '') => {
            const version: Version<TContent> = {
                id: id,
                date: new Date().toISOString(),
                type: type,
                content: content,
                name: name,
            };
            history.value.push(version as Version<unknown>);

            // Enforce history length limit
            clearOldVersions(MAX_HISTORY_LENGTH);
        };

        const revertToVersion = <TContent>(versionDate: string): Version<TContent> | undefined => {
            return history.value.find((v) => v.date === versionDate) as Version<TContent> | undefined;
        };

        const deleteVersion = (versionId: string) => {
            history.value = history.value.filter((v) => v.date !== versionId);
        };

        const clearOldVersions = (maxVersions: number) => {
            const typeCounts: Record<string, number> = {};

            history.value = history.value.filter((version) => {
                typeCounts[version.type] = (typeCounts[version.type] || 0) + 1;

                if (typeCounts[version.type]! > maxVersions) {
                    return false; // Remove excess versions for this type
                }

                return true; // Keep the version
            });
        };

        const getVersionsByType = <TContent>(type: string) =>
            computed(() => history.value.filter((v) => v.type === type) as Version<TContent>[]);

        const handleRefresh = (handleAutoSave: () => void) => {
            // For refresh/close/tab close
            onMounted(() => {
                window.addEventListener('beforeunload', handleAutoSave);
            });
            onBeforeUnmount(() => {
                window.removeEventListener('beforeunload', handleAutoSave);
            });

            // For SPA navigation
            onBeforeRouteLeave((_to, _from, next) => {
                handleAutoSave();
                next();
            });
        };

        return {
            // State
            history,

            // actions
            listHistory,
            addVersion,
            revertToVersion,
            deleteVersion,
            clearOldVersions,
            getVersionsByType,
            handleRefresh,
        };
    },
    {
        persist: true,
    },
);
