import { defineStore, type StoreDefinition } from 'pinia';

export interface CookiesState {
    cookiesState: null | boolean;
}

export const useCookiesStore = defineStore('cookies', {
    state: () =>
        ({
            cookiesState: null,
        }) as CookiesState,
    persist: true,
    getters: {
        hasCookiesAccepted(state): boolean {
            return state.cookiesState === true;
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useCookiesStore as unknown as StoreDefinition, import.meta.hot));
}
