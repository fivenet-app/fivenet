import { StoreDefinition, defineStore } from 'pinia';

export interface SettingsState {
    version: string;
    locale: string;
    livemap: {
        markerSize: number;
        centerSelectedMarker: boolean;
    };
}

export const useSettingsStore = defineStore('settings', {
    state: () =>
        ({
            version: __APP_VERSION__ as string,
            locale: 'de',
            livemap: {
                markerSize: 22,
                centerSelectedMarker: false,
            },
        }) as SettingsState,
    persist: true,
    actions: {
        setVersion(version: string): void {
            this.version = version;
        },
        setLocale(locale: string): void {
            this.locale = locale;
        },
    },
    getters: {
        getVersion: (state) => state.version,
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useSettingsStore as unknown as StoreDefinition, import.meta.hot));
}
