import { StoreDefinition, defineStore } from 'pinia';

export interface UserSettingsState {
    version: string;
    locale: string;
    livemapMarkerSize: number;
    livemapCenterSelectedMarker: boolean;
}

export const useUserSettingsStore = defineStore('userSettings', {
    state: () =>
        ({
            version: __APP_VERSION__ as string,
            locale: 'de',
            livemapMarkerSize: 22,
            livemapCenterSelectedMarker: false,
        }) as UserSettingsState,
    persist: true,
    actions: {
        setVersion(version: string): void {
            this.version = version;
        },
        setLocale(locale: string): void {
            this.locale = locale;
        },
        setLivemapMarkerSize(size: number): void {
            this.livemapMarkerSize = size;
        },
        setLivemapCenterSelectedMarker(value: boolean): void {
            this.livemapCenterSelectedMarker = value;
        },
    },
    getters: {
        getVersion: (state) => state.version,
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useUserSettingsStore as unknown as StoreDefinition, import.meta.hot));
}
