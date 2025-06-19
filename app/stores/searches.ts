import { defineStore } from 'pinia';

/**
 * A single Pinia store that holds *all* of your search states,
 * keyed by an arbitrary string.
 */
export const useSearchesStore = defineStore('searches', () => {
    // Reactive object so new keys are tracked
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const searches = reactive<Record<string, any>>({});

    function setSearch<S>(key: string, state: S) {
        searches[key] = state;
    }

    function getSearch<S>(key: string): S | undefined {
        return searches[key] as S | undefined;
    }

    return { searches, setSearch, getSearch };
});
