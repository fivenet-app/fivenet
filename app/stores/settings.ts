import type { RoutePathSchema } from '@typed-router';
import { defineStore } from 'pinia';
import type { Locale } from 'vue-i18n';
import type { Perms } from '~~/gen/ts/perms';

const logger = useLogger('⚙️ Settings');

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
    disabled?: boolean;
    order?: number;
};

export type LivemapLayerCategory = {
    key: string;
    label: string;
    order?: number;
};

export type LivemapSettings = {
    markerSize: number;
    centerSelectedMarker: boolean;
    showUnitNames: boolean;
    showUnitStatus: boolean;
    showAllDispatches: boolean;
    showGrid: boolean;
    showHeatmap: boolean;
    useUnitColor: boolean;
};

export type SignatureSettings = {
    penColor: string;
    minStrokeWidth: number;
    maxStrokeWidth: number;
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

        const eventsDisabled = ref<boolean>(false);

        const livemap = ref<LivemapSettings>({
            markerSize: 22,
            centerSelectedMarker: false,
            showUnitNames: true,
            showUnitStatus: true,
            showAllDispatches: false,
            showGrid: false,
            showHeatmap: false,
            useUnitColor: true,
        });

        const centrum = ref({
            dispatchListCardStyle: false,
        });

        const livemapTileLayer = ref<string>('postal');
        const livemapLayers = ref<LivemapLayer[]>([]);
        const livemapLayerCategories = ref<LivemapLayerCategory[]>([]);

        const startpage = ref<RoutePathSchema>('/overview');
        const design = ref({
            documents: {
                listStyle: 'single' as 'single' | 'double',
                viewCollapsedTitle: false,
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
        const jobsService = ref({
            cardView: true,
        });

        const editor = ref<{ showInvisibleCharacters: boolean }>({
            showInvisibleCharacters: false,
        });

        const signature = ref<SignatureSettings>({
            penColor: 'rgb(0, 0, 0)',
            minStrokeWidth: 2,
            maxStrokeWidth: 6,
        });

        // Actions
        /**
         * Set the application version.
         *
         * @param {string} newVersion - The new version string to set.
         */
        const setVersion = (newVersion: string): void => {
            version.value = newVersion;
        };

        /**
         * Set the update availability status.
         *
         * @param {string} newVersion - The version string indicating the available update.
         */
        const setUpdateAvailable = async (newVersion: string): Promise<void> => {
            updateAvailable.value = newVersion;
        };

        /**
         * Configure NUI settings.
         *
         * @param {boolean} enabled - Whether NUI is enabled.
         * @param {string | undefined} [resourceName] - The name of the NUI resource, if applicable.
         */
        const setNuiSettings = (enabled: boolean, resourceName?: string | undefined): void => {
            nuiEnabled.value = enabled;
            nuiResourceName.value = resourceName;
        };

        /**
         * Add or update a livemap layer category.
         *
         * @param {LivemapLayerCategory} category - The category object containing the data to add or update.
         */
        const addOrUpdateLivemapCategory = (category: LivemapLayerCategory): void => {
            const idx = livemapLayerCategories.value.findIndex((l) => l.key === category.key);
            if (idx === -1) {
                livemapLayerCategories.value.push(category);
                return;
            }

            const current = livemapLayerCategories.value[idx]!;
            if (category.label && current.label !== category.label) {
                current.label = category.label;
            }
            current.order = category.order;
        };

        /**
         * Add or update a livemap layer.
         *
         * @param {LivemapLayer} layer - The layer object containing the data to add or update.
         */
        const addOrUpdateLivemapLayer = (layer: LivemapLayer): void => {
            const idx = livemapLayers.value.findIndex((l) => l.key === layer.key);
            if (idx === -1) {
                layer.visible = true;
                livemapLayers.value.push(layer);
                return;
            }

            const current = livemapLayers.value[idx]!;
            if (layer.label && current.label !== layer.label) {
                current.label = layer.label;
            }
            current.category = layer.category;
            if (layer.visible !== undefined) {
                current.visible = layer.visible;
            }
            current.perm = layer.perm;
            current.attr = layer.attr;
            current.disabled = layer.disabled;
            current.order = layer.order;
        };

        /**
         * Remove a livemap layer by its key.
         *
         * @param {string} key - The key of the layer to remove.
         */
        const removeLivemapLayer = (key: string): void => {
            const idx = livemapLayers.value.findIndex((l) => l.key === key);
            if (idx === -1) return;
            livemapLayers.value.splice(idx, 1);
        };

        // Getters
        /**
         * Get the user's locale.
         *
         * @returns {Locale} - The user's locale, falling back to the default locale if not set.
         */
        const getUserLocale = computed<Locale>(() => {
            if (locale.value !== undefined) {
                return locale.value;
            }
            if (useAppConfig().defaultLocale !== '') {
                return useAppConfig().defaultLocale as Locale;
            }
            return 'en';
        });

        const getLogger = (): ILogger => logger;

        return {
            // State
            version,
            updateAvailable,
            locale,

            nuiEnabled,
            nuiResourceName,

            eventsDisabled,

            livemap,
            livemapLayerCategories,
            livemapLayers,
            livemapTileLayer,
            centrum,
            startpage,
            design,
            audio,
            calendar,
            streamerMode,
            jobsService,
            editor,
            signature,

            // Actions
            getLogger,
            setVersion,
            setUpdateAvailable,
            setNuiSettings,
            addOrUpdateLivemapCategory,
            addOrUpdateLivemapLayer,
            removeLivemapLayer,

            // Getters
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
