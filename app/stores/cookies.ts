import { defineStore } from 'pinia';

export const useCookiesStore = defineStore(
    'cookies',
    () => {
        // State
        const cookiesState = ref<null | boolean>(null);

        // Getters
        const hasCookiesAccepted = computed(() => cookiesState.value === true);

        return {
            cookiesState,
            hasCookiesAccepted,
        };
    },
    {
        persist: true,
    },
);

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useCookiesStore, import.meta.hot));
}
