import { defineStore } from 'pinia';
import type { Version } from '~/types/history';

/**
 * Maximum number of history entries allowed per type.
 */
const MAX_HISTORY_LENGTH = 20;

export const useHistoryStore = defineStore(
    'historyStore',
    () => {
        /**
         * State: Stores the history of versions for different types.
         * Each entry contains metadata such as type, date, and content.
         */
        const history = ref<Version<unknown>[]>([]);

        /**
         * Lists history entries filtered by a specific type.
         * @template TContent - The type of the content stored in the version.
         * @param {string} type - The type of history to filter by.
         * @returns {ComputedRef<Version<TContent>[]>} - A computed reference to the filtered history.
         */
        const listHistory = <TContent>(type: string): ComputedRef<Version<TContent>[]> => {
            return computed(() => {
                return history.value.filter((v) => v.type === type) as Version<TContent>[];
            });
        };

        /**
         * Adds a new version to the history.
         * @template TContent - The type of the content stored in the version.
         * @param {string} type - The type of the version.
         * @param {number} id - The unique identifier for the version.
         * @param {TContent} content - The content of the version.
         * @param {string} [name=''] - An optional name for the version.
         */
        const addVersion = <TContent>(type: string, id: number, content: TContent, name: string = '') => {
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

        /**
         * Reverts to a specific version by its date.
         * @template TContent - The type of the content stored in the version.
         * @param {string} versionDate - The date of the version to revert to.
         * @returns {Version<TContent> | undefined} - The version to revert to, or undefined if not found.
         */
        const revertToVersion = <TContent>(versionDate: string): Version<TContent> | undefined => {
            return history.value.find((v) => v.date === versionDate) as Version<TContent> | undefined;
        };

        /**
         * Deletes a version by its unique identifier.
         * @param {string} versionId - The unique identifier of the version to delete.
         */
        const deleteVersion = (versionId: string) => {
            history.value = history.value.filter((v) => v.date !== versionId);
        };

        /**
         * Clears old versions exceeding the maximum allowed per type.
         * @param {number} maxVersions - The maximum number of versions allowed per type.
         */
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

        /**
         * Retrieves versions filtered by a specific type.
         * @template TContent - The type of the content stored in the version.
         * @param {string} type - The type of history to filter by.
         * @returns {ComputedRef<Version<TContent>[]>} - A computed reference to the filtered history.
         */
        const getVersionsByType = <TContent>(type: string): ComputedRef<Version<TContent>[]> =>
            computed(() => history.value.filter((v) => v.type === type) as Version<TContent>[]);

        /**
         * Sets up watchers for handling refresh events, such as page reloads or SPA navigation.
         * @param {() => void} handleAutoSave - A callback function to handle auto-saving before refresh events.
         */
        const handleRefresh = (handleAutoSave: () => void) => {
            // For refresh/close/tab close
            onBeforeMount(() => {
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

            // Actions
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

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useHistoryStore, import.meta.hot));
}
