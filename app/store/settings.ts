import { defineStore } from 'pinia';
import type { Locale } from 'vue-i18n';
import type { Perms } from '~~/gen/ts/perms';

export const logger = useLogger('⚙️ Settings');

export type LivemapLayer = {
    key: string;
    label: string;
    category: string;
    visible?: boolean;
    perm?: Perms;
    attr?: {
        key: string;
        val: string;
    };
};

export type LivemapLayerCategory = {
    key: string;
    label: string;
};

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
        });

        const livemapLayers = ref<LivemapLayer[]>([]);
        const livemapLayerCategories = ref<LivemapLayerCategory[]>([]);

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

        const addOrUpdateLivemapCategory = (category: LivemapLayerCategory): void => {
            const idx = livemapLayerCategories.value.findIndex((l) => l.key === category.key);
            if (idx === -1) {
                livemapLayerCategories.value.push(category);
                return;
            }

            const current = livemapLayerCategories.value[idx]!;
            if (current.label !== category.label) {
                current.label = category.label;
            }
            current.label = category.label;
        };
        const addOrUpdateLivemapLayer = (layer: LivemapLayer): void => {
            const idx = livemapLayers.value.findIndex((l) => l.key === layer.key);
            if (idx === -1) {
                layer.visible = true;
                livemapLayers.value.push(layer);
                return;
            }

            const current = livemapLayers.value[idx]!;
            if (current.label !== layer.label) {
                current.label = layer.label;
            }
            current.category = layer.category;
            if (layer.visible !== undefined) {
                current.visible = layer.visible;
            }
            current.perm = layer.perm;
            current.attr = layer.attr;
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
            livemapLayerCategories,
            livemapLayers,
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
            addOrUpdateLivemapCategory,
            addOrUpdateLivemapLayer,

            getUserLocale,
        };
    },
    {
        persist: {
            omit: ['updateAvailable', 'livemapLayerCategories'],
        },
    },
);

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useSettingsStore, import.meta.hot));
}
