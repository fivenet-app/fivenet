import { defineStore } from 'pinia';

/**
 * Pinia store for managing UI state.
 */
export const useUIStateStore = defineStore(
    'uiState',
    () => {
        // State
        /**
         * Tracks window focus state (true if focused, false otherwise)
         */
        const windowFocus = useWindowFocus();

        return {
            // State
            windowFocus,
        };
    },
    {
        persist: false,
    },
);

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useUIStateStore, import.meta.hot));
}
