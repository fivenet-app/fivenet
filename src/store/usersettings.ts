import { StoreDefinition, defineStore } from 'pinia';

export interface UserSettingsState {
    locale: string;
    livemapMarkerSize: number;
}

export const useUserSettingsStore = defineStore('notificator', {
    state: () =>
    ({
        locale: 'en-US',
        livemapMarkerSize: 26,
    } as UserSettingsState),
    persist: true,
    actions: {
        setLocale(locale: string): void {
            this.locale = locale;
        },
        getLivemapMarkerSize(): number {
            return this.livemapMarkerSize;
        },
        setLivemapMarkerSize(size: number): void {
            this.livemapMarkerSize = size;
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useUserSettingsStore as unknown as StoreDefinition, import.meta.hot));
}
