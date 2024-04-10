import { defineStore, type StoreDefinition } from 'pinia';

export interface SettingsState {
    version: string;
    updateAvailable: false | string;
    locale: string | null;
    cookiesState: undefined | null | boolean;

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
    startpage: string;
    design: {
        docEditorTheme: 'default' | 'dark';
        ui: {
            primary: string;
            gray: string;
        };
    };
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
            cookiesState: undefined,

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
            startpage: '/overview',
            design: {
                docEditorTheme: 'default',
                ui: {
                    primary: 'sky',
                    gray: 'cool',
                },
            },
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
            'cookiesState',
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
        hasCookiesAccepted(state): boolean {
            return state.cookiesState === true;
        },
        isNUIAvailable(state): boolean {
            return state.nuiEnabled ?? false;
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useSettingsStore as unknown as StoreDefinition, import.meta.hot));
}
