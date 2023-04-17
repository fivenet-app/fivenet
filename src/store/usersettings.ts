import { defineStore } from 'pinia';

export interface UserSettingsState {
    locale: string;
}

export const useUserSettingsStore = defineStore('notificator', {
    state: () =>
    ({
        locale: 'en-US',
    } as UserSettingsState),
    persist: true,
    actions: {
        setLocale(locale: string): void {
            this.locale = locale;
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useUserSettingsStore, import.meta.hot));
}
