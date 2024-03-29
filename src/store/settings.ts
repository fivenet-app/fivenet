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
    updateAvailable: false | string;
    locale: string | null;
    nuiEnabled: boolean;
    nuiResourceName: string | undefined;
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
    featureGates: {
        desktop: boolean;
    };
}

export const useSettingsStore = defineStore('settings', {
    state: () =>
        ({
            version: __APP_VERSION__ as string,
            updateAvailable: false,
            locale: null,
            nuiEnabled: false,
            nuiResourceName: undefined,

            livemap: {
                markerSize: 22,
                centerSelectedMarker: false,
                showUnitNames: true,
                showUnitStatus: true,
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
            featureGates: {
                desktop: false,
            },
        }) as SettingsState,
    persist: {
        paths: [
            'version',
            'locale',
            'nuiEnabled',
            'nuiResourceName',
            'livemap',
            'documents',
            'startpage',
            'theme',
            'audio',
            'streamerMode',
            'featureGates',
        ],
    },
    actions: {
        setVersion(version: string): void {
            this.version = version;
        },
        async setUpdateAvailable(version: string): Promise<void> {
            this.updateAvailable = version;
        },
        setNuiDetails(enabled: boolean, resourceName: string | undefined): void {
            this.nuiEnabled = enabled;
            this.nuiResourceName = resourceName;
        },
        setLocale(locale: string): void {
            this.locale = locale;
        },
    },
    getters: {
        isNUIAvailable(state): boolean {
            return state.nuiEnabled ?? false;
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useSettingsStore as unknown as StoreDefinition, import.meta.hot));
}
