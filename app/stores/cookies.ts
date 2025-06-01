import { defineStore } from 'pinia';

export const useCookiesStore = defineStore(
    'cookies',
    () => {
        // State
        const cookiesState = ref<null | boolean>(null);
        const isConsentModalOpen = ref(false);

        watch(cookiesState, (newValue) => {
            if (newValue === true) {
                isConsentModalOpen.value = false;
            } else if (newValue === null) {
                isConsentModalOpen.value = true;
            }
        });

        // Getters
        const hasCookiesAccepted = computed(() => cookiesState.value === true);

        return {
            cookiesState,
            isConsentModalOpen,
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
