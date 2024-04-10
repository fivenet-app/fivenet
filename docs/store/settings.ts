import { defineStore, type StoreDefinition } from 'pinia';

export interface SettingsState {
    version: string;
    locale: string | null;
    cookiesState: undefined | null | boolean;
}

export const useSettingsStore = defineStore('settings', {
    state: () =>
        ({
            version: __APP_VERSION__ as string,
            locale: null,
            cookiesState: undefined,
        }) as SettingsState,
    persist: {
        paths: ['version', 'locale', 'cookiesState'],
    },
    actions: {
        setLocale(locale: string): void {
            this.locale = locale;
        },
    },
    getters: {
        hasCookiesAccepted(state): boolean {
            return state.cookiesState === true;
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useSettingsStore as unknown as StoreDefinition, import.meta.hot));
}
