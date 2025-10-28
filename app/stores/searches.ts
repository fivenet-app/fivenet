import { defineStore } from 'pinia';

/**
 * A Pinia store for managing and persisting search masks.
 * This store provides methods to set, retrieve, and clear search states.
 */
export const useSearchesStore = defineStore('searches', () => {
    // State
    /**
     * Reactive object to hold all search states, keyed by a unique string.
     */
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const searches = reactive<Record<string, any>>({});

    // Actions
    /**
     * Sets the state for a specific search key.
     * @param key - The unique identifier for the search state.
     * @param state - The state object to associate with the key.
     */
    const setSearch = <S>(key: string, state: S): void => {
        searches[key] = state;
    };

    /**
     * Retrieves the state for a specific search key.
     * @param key - The unique identifier for the search state.
     * @returns The state object associated with the key, or undefined if not found.
     */
    const getSearch = <S>(key: string): S | undefined => {
        return searches[key] as S | undefined;
    };

    /**
     * Clears all search states.
     */
    const clear = (): void => {
        Object.keys(searches).forEach((key) => {
            // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
            delete searches[key];
        });
    };

    return { searches, setSearch, getSearch, clear };
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useSearchesStore, import.meta.hot));
}
