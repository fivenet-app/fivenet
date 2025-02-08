import { defineStore } from 'pinia';
import type { Locale } from 'vue-i18n';

export const logger = useLogger('⚙️ Settings');

export interface SettingsState {
    updateAvailable: false | string;
    version: string;
    locale: Locale | undefined;

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
        documents: {
            listStyle: 'single' | 'double';
        };
        ui: {
            primary: string;
            gray: string;
        };
    };
    audio: {
        notificationsVolume: number;
    };
    calendar: {
        reminderTimes: number[];
    };
    streamerMode: boolean;
    calculatorPosition: 'top' | 'middle' | 'bottom';
    jobsService: {
        cardView: boolean;
    };
}

export const useSettingsStore = defineStore('settings', {
    state: () =>
        ({
            version: APP_VERSION,
            updateAvailable: false,
            locale: undefined,

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
                documents: {
                    listStyle: 'single',
                },
                ui: {
                    primary: 'sky',
                    gray: 'neutral',
                },
            },
            audio: {
                notificationsVolume: 0.15,
            },
            calendar: {
                reminderTimes: [0, 900],
            },
            streamerMode: false,
            calculatorPosition: 'middle',
            jobsService: {
                cardView: true,
            },
        }) as SettingsState,
    persist: {
        omit: ['updateAvailable'],
    },
    actions: {
        setVersion(version: string): void {
            this.version = version;
        },
        async setUpdateAvailable(version: string): Promise<void> {
            this.updateAvailable = version;
        },
        setNuiSettings(enabled: boolean, resourceName: string | undefined): void {
            this.nuiEnabled = enabled;
            this.nuiResourceName = resourceName;
        },
    },
    getters: {
        isNUIEnabled(state): boolean {
            return state.nuiEnabled ?? false;
        },
        getUserLocale(state): Locale {
            return state.locale !== undefined
                ? state.locale
                : useAppConfig().defaultLocale !== ''
                  ? (useAppConfig().defaultLocale as Locale)
                  : 'en';
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useSettingsStore, import.meta.hot));
}
