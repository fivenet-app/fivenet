import { defineStore } from 'pinia';
import type { Locale } from 'vue-i18n';

export const logger = useLogger('⚙️ Settings');

export const useSettingsStore = defineStore(
    'settings',
    () => {
        // State
        const version = ref<string>(APP_VERSION);
        const updateAvailable = ref<false | string>(false);
        const locale = ref<Locale | undefined>(undefined);

        const nuiEnabled = ref<boolean>(false);
        const nuiResourceName = ref<string | undefined>(undefined);

        const livemap = ref({
            markerSize: 22,
            centerSelectedMarker: false,
            showUnitNames: true,
            showUnitStatus: true,
            showAllDispatches: false,
            activeLayers: [] as string[],
        });

        const startpage = ref<string>('/overview');
        const design = ref({
            documents: {
                listStyle: 'single' as 'single' | 'double',
            },
            ui: {
                primary: 'sky',
                gray: 'neutral',
            },
        });

        const audio = ref({
            notificationsVolume: 0.15,
        });

        const calendar = ref({
            reminderTimes: [0, 900],
        });

        const streamerMode = ref<boolean>(false);
        const calculatorPosition = ref<'top' | 'middle' | 'bottom'>('middle');
        const jobsService = ref({ cardView: true });

        // Actions
        const setVersion = (newVersion: string): void => {
            version.value = newVersion;
        };

        const setUpdateAvailable = async (newVersion: string): Promise<void> => {
            updateAvailable.value = newVersion;
        };

        const setNuiSettings = (enabled: boolean, resourceName: string | undefined): void => {
            nuiEnabled.value = enabled;
            nuiResourceName.value = resourceName;
        };

        // Getters
        const getUserLocale = computed<Locale>(() => {
            if (locale.value !== undefined) {
                return locale.value;
            }
            if (useAppConfig().defaultLocale !== '') {
                return useAppConfig().defaultLocale as Locale;
            }
            return 'en';
        });

        return {
            version,
            updateAvailable,
            locale,
            nuiEnabled,
            nuiResourceName,
            livemap,
            startpage,
            design,
            audio,
            calendar,
            streamerMode,
            calculatorPosition,
            jobsService,

            setVersion,
            setUpdateAvailable,
            setNuiSettings,

            getUserLocale,
        };
    },
    {
        persist: {
            omit: ['updateAvailable'],
        },
    },
);

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useSettingsStore, import.meta.hot));
}
