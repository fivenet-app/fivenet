import { StoreDefinition, defineStore } from 'pinia';

export interface UserSettingsState {
    locale: string;
    livemapMarkerSize: number;
    livemapCenterSelectedMarker: boolean;
}

export const useUserSettingsStore = defineStore('userSettings', {
    state: () =>
    ({
        locale: 'de',
        livemapMarkerSize: 22,
        livemapCenterSelectedMarker: false,
    } as UserSettingsState),
    persist: true,
    actions: {
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
        getLocale: (state): string => state.locale,
        getLivemapMarkerSize: (state): number => state.livemapMarkerSize,
        getLivemapCenterSelectedMarker: (state): boolean => state.livemapCenterSelectedMarker,
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useUserSettingsStore as unknown as StoreDefinition, import.meta.hot));
}
