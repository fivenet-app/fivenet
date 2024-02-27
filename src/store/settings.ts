import { defineStore, type StoreDefinition } from 'pinia';

export const JOB_THEME_KEY = '__job_theme__';

export const availableThemes = [
    { name: 'Default', key: 'defaultTheme' },
    { name: 'Baddie Orange', key: 'themeBaddieOrange' },
    { name: 'Baddie Pink', key: 'themeBaddiePink' },
    { name: 'Baddie Yellow', key: 'themeBaddieYellow' },
    { name: 'Da Medic', key: 'themeDaMedic' },
    { name: 'Purple', key: 'themePurple' },
];

export interface SettingsState {
    version: string;
    locale: string;
    livemap: {
        markerSize: number;
        centerSelectedMarker: boolean;
        showUnitNames: boolean;
        showUnitStatus: boolean;
        showAllDispatches: boolean;
        activeLayers: string[];
    };
    documents: {
        editorTheme: 'default' | 'dark';
    };
    startpage: string;
    theme: string;
    audio: {
        notificationsVolume: number;
    };
    streamerMode: boolean;
    toggles: {};
}

export const useSettingsStore = defineStore('settings', {
    state: () =>
        ({
            version: __APP_VERSION__ as string,
            locale: 'de',
            livemap: {
                markerSize: 22,
                centerSelectedMarker: false,
                showUnitNames: false,
                showUnitStatus: false,
                showAllDispatches: false,
                activeLayers: [],
            },
            documents: {
                editorTheme: 'default',
            },
            startpage: '/overview',
            theme: JOB_THEME_KEY,
            audio: {
                notificationsVolume: 0.15,
            },
            streamerMode: false,
            toggles: {},
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
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useSettingsStore as unknown as StoreDefinition, import.meta.hot));
}
