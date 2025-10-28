import { defineStore } from 'pinia';
import { computed, ref, watch } from 'vue';

/**
 * Pinia store for managing cookie consent state and modal visibility.
 */
export const useCookiesStore = defineStore(
    'cookies',
    () => {
        // State
        /**
         * Cookie consent state: null (undecided), true (accepted), false (rejected)
         */
        const cookiesState = ref<null | boolean>(null);
        /**
         * Controls visibility of the consent modal
         */
        const isConsentModalOpen = ref(false);

        // Watchers
        /**
         * Watch cookiesState and update modal visibility accordingly
         */
        watch(cookiesState, (newValue) => {
            if (newValue === true) {
                isConsentModalOpen.value = false;
            } else if (newValue === null) {
                isConsentModalOpen.value = true;
            }
        });

        // Getters
        /**
         * Returns true if cookies have been accepted
         */
        const hasCookiesAccepted = computed(() => cookiesState.value === true);

        // Methods
        /**
         * Accept cookies and close modal
         */
        const acceptCookies = (): void => {
            cookiesState.value = true;
            isConsentModalOpen.value = false;
        };

        /**
         * Reject cookies and close modal
         */
        const rejectCookies = (): void => {
            cookiesState.value = false;
            isConsentModalOpen.value = false;
        };

        /**
         * Open the consent modal
         */
        const openConsentModal = (): void => {
            isConsentModalOpen.value = true;
        };

        /**
         * Close the consent modal
         */
        const closeConsentModal = (): void => {
            isConsentModalOpen.value = false;
        };

        return {
            // State
            cookiesState,
            isConsentModalOpen,

            // Getters
            hasCookiesAccepted,

            // Methods
            acceptCookies,
            rejectCookies,
            openConsentModal,
            closeConsentModal,
        };
    },
    {
        persist: true,
    },
);

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useCookiesStore, import.meta.hot));
}
