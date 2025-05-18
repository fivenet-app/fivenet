import { defineStore } from 'pinia';

export const useUIState = defineStore(
    'uiState',
    () => {
        // State
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
    import.meta.hot.accept(acceptHMRUpdate(useSettingsStore, import.meta.hot));
}
