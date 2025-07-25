<script lang="ts" setup>
// eslint-disable-next-line @typescript-eslint/consistent-type-imports
import L, { CRS, extend, LatLng, latLngBounds, type PointExpression, Projection, Transformation } from 'leaflet';
import 'leaflet-contextmenu';
import 'leaflet.heat';
import ZoomControls from '~/components/livemap/controls/ZoomControls.vue';
import { simpleGraticule } from '~/composables/leaflet/L.SimpleGraticule';
import { useLivemapStore } from '~/stores/livemap';
import { backgroundColorList, tileLayers } from '~/types/livemap';
import type { ValueOf } from '~/utils/types';
import LayerControls from './controls/LayerControls.vue';
import HeatmapLayer from './HeatmapLayer.vue';

defineProps<{
    /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
    mapOptions?: Record<string, any>;
}>();

const emit = defineEmits<{
    (e: 'mapReady', map: L.Map): void;
    (e: 'overlayadd', event: L.LayersControlEvent): void;
    (e: 'overlayremove', event: L.LayersControlEvent): void;
}>();

const slideover = useSlideover();

const { can } = useAuth();

const settingsStore = useSettingsStore();
const { livemapTileLayer, livemap: livemapSettings } = storeToRefs(settingsStore);

const livemapStore = useLivemapStore();
const { location, selectedMarker, zoom } = storeToRefs(livemapStore);

let map: L.Map | undefined;

function mapResize(): void {
    if (map === undefined) {
        return;
    }

    map.invalidateSize();
}

const mapContainer = useTemplateRef('mapContainer');
const mapResizeDebounced = useDebounceFn(mapResize, 350, { maxWait: 700 });
useResizeObserver(mapContainer, (_) => mapResizeDebounced());

const centerX = 117.3;
const centerY = 172.8;
const scaleX = 0.02072;
const scaleY = 0.0205;

const bounds = latLngBounds([-4_000, -4_000], [8_000, 8_000]);
const maxBounds = latLngBounds([-9_000, -9_000], [11_000, 11_000]);

const customCRS = extend({}, CRS.Simple, {
    projection: Projection.LonLat,
    scale: function (zoom: number): number {
        return Math.pow(2, zoom);
    },
    zoom: function (sc: number): number {
        return Math.log(sc) / 0.6931471805599453;
    },
    distance: function (pos1: L.LatLng, pos2: L.LatLng): number {
        const xDiff = pos2.lng - pos1.lng;
        const yDiff = pos2.lat - pos1.lat;
        return Math.sqrt(xDiff * xDiff + yDiff * yDiff);
    },
    transformation: new Transformation(scaleX, centerX, -scaleY, centerY),
    infinite: true,
});

// eslint-disable-next-line prefer-const
let center: PointExpression = [0, 0];

const mouseLat = ref<number>(0);
const mouseLong = ref<number>(0);

const currentLocationQuery = useRouteQuery<string>('loc', '');

function getZoomOffset(zoom: number): number {
    if (!slideover.isOpen.value) return 0;

    switch (zoom) {
        case 1:
            return 2150;
        case 2:
            return 1650;
        case 3:
            return 1350;
        case 4:
            return 750;
        case 5:
            return 425;
        case 6:
            return 225;
        case 7:
            return 100;
        default:
            return 400;
    }
}

watch(selectedMarker, async () => {
    if (map === undefined || selectedMarker.value === undefined) return;

    map?.panTo([selectedMarker.value.y, selectedMarker.value.x + getZoomOffset(zoom.value)], {
        animate: true,
        duration: 0.75,
    });
});

watch(location, async () => {
    if (map === undefined || location.value === undefined) return;

    map.setView([location.value.y, location.value.x + getZoomOffset(zoom.value)], zoom.value, {
        animate: false,
        duration: 0.75,
    });
});

const isMoving = ref<boolean>(false);

watchDebounced(
    isMoving,
    async () => {
        if (map === undefined || isMoving.value) {
            return;
        }

        const newHash = stringifyHash(map.getZoom(), map.getCenter().lat, map.getCenter().lng);
        if (currentLocationQuery.value !== newHash) {
            currentLocationQuery.value = newHash;
        }
    },
    { debounce: 1000, maxWait: 3000 },
);

const backgroundColor = ref<ValueOf<typeof backgroundColorList>>(backgroundColorList.postal);

async function updateBackground(layer: string): Promise<void> {
    switch (layer) {
        case 'satelite':
            backgroundColor.value = backgroundColorList.satelite;
            break;

        case 'postal':
        default:
            backgroundColor.value = backgroundColorList.postal;
            break;
    }
}

watch(livemapTileLayer, async (layer) => updateBackground(layer));

function stringifyHash(currZoom: number, centerLat: number, centerLong: number): string {
    const precision = Math.max(0, Math.ceil(Math.log(zoom.value) / Math.LN2));

    const hash = [currZoom, centerLat.toFixed(precision), centerLong.toFixed(precision)].join('/');
    return hash;
}

function parseLocationQuery(query: string): { latlng: L.LatLng; zoom: number } | undefined {
    const args = query.split('/');

    const zoom = args[0] ? parseInt(args[0]) : 2;
    const lat = args[1] ? parseFloat(args[1]) : 0;
    const lng = args[2] ? parseFloat(args[2]) : 0;

    if (isNaN(zoom) || isNaN(lat) || isNaN(lng)) return;

    return {
        latlng: new LatLng(lat, lng),
        zoom,
    };
}

const graticuleLayer = simpleGraticule({
    interval: 200,
    showOriginLabel: true,
    redraw: 'moveend',
    zoomIntervals: [
        { start: 1, end: 1, interval: 1000 },
        { start: 2, end: 3, interval: 500 },
        { start: 4, end: 5, interval: 250 },
        { start: 6, end: 7, interval: 100 },
    ],
});

const heat = ref<L.HeatLayer | undefined>(undefined);

async function onMapReady(m: L.Map): Promise<void> {
    updateBackground(livemapTileLayer.value);

    map = m;
    map.invalidateSize();

    const startPos = parseLocationQuery(currentLocationQuery.value as string);
    if (startPos) {
        map.setView(startPos.latlng, startPos.zoom);
    }

    map.on('baselayerchange', async (event: L.LayersControlEvent) => updateBackground(event.name));

    map.on('overlayadd', (event) => emit('overlayadd', event));
    map.on('overlayremove', (event) => emit('overlayremove', event));

    map.addEventListener('mousemove', async (event: L.LeafletMouseEvent) => {
        if (!event.latlng) return;

        mouseLat.value = Math.round(event.latlng.lat * 100000) / 100000;
        mouseLong.value = Math.round(event.latlng.lng * 100000) / 100000;
    });

    map.on('movestart', async () => {
        isMoving.value = true;
    });

    map.on('moveend', async () => {
        isMoving.value = false;
    });

    emit('mapReady', map);

    heat.value = await useLHeat({
        leafletObject: map,
        heatPoints: [],
        radius: 10,
    });

    if (livemapSettings.value.showGrid) {
        graticuleLayer.addTo(map);
    }
}

provide('heat', heat);

watch(
    () => livemapSettings.value.showGrid,
    (newVal) => {
        if (!graticuleLayer) return;

        if (newVal && map) {
            graticuleLayer.addTo(map);
        } else {
            graticuleLayer.remove();
        }
    },
);

onBeforeUnmount(() => {
    map = undefined;
});
</script>

<template>
    <div ref="mapContainer" class="mapContainer flex h-full flex-row" :style="{ backgroundColor }">
        <LMap
            v-model:zoom="zoom"
            v-model:center="center"
            :bounds="bounds"
            :max-bounds="maxBounds"
            :crs="customCRS"
            :min-zoom="1"
            :max-zoom="7"
            :inertia="false"
            :style="{ backgroundColor: 'rgba(0,0,0,0.0)' }"
            :use-global-leaflet="true"
            :options="mapOptions"
            @click="selectedMarker = undefined"
            @ready="onMapReady($event)"
        >
            <LTileLayer
                v-for="layer in tileLayers"
                :key="layer.key"
                :url="layer.url"
                layer-type="base"
                :name="$t(layer.label)"
                :no-wrap="true"
                :tms="true"
                :visible="livemapTileLayer === layer.key"
                :min-zoom="1"
                :max-zoom="layer.options?.maxZoom || 7"
                :attribution="layer.options?.attribution || ''"
            />

            <ZoomControls />

            <LayerControls>
                <div v-if="can('centrum.CentrumService/TakeControl').value">
                    <div class="mt-1 inline-flex gap-1 overflow-y-hidden px-1">
                        <UToggle v-model="livemapSettings.showHeatmap" />
                        <span class="truncate hover:line-clamp-2">{{ $t('common.heatmap') }}</span>
                    </div>
                </div>

                <slot name="layerControls" />
            </LayerControls>

            <!-- eslint-disable-next-line tailwindcss/no-custom-classname -->
            <LControl class="leaflet-control-attribution" position="bottomleft">
                <span class="font-semibold">{{ $t('common.longitude') }}:</span> {{ mouseLat.toFixed(3) }} |
                <span class="font-semibold">{{ $t('common.latitude') }}:</span> {{ mouseLong.toFixed(3) }}
            </LControl>

            <slot />

            <HeatmapLayer :show="livemapSettings.showHeatmap" />
        </LMap>

        <slot name="afterMap" />
    </div>
</template>

<style scoped>
.mapContainer:deep(.leaflet-container) {
    font-family: var(--font-sans);

    .leaflet-container a {
        color: rgb(var(--color-primary-500));
    }
    .leaflet-container a:hover {
        color: rgb(var(--color-primary-400));
    }

    .leaflet-map-pane {
        z-index: 0;
    }
    .leaflet-overlay-pane {
        z-index: 400;
    }

    .leaflet-div-icon {
        background: none !important;
        border: none !important;
    }

    .leaflet-div-icon svg path {
        stroke: #000000;
        stroke-width: 0.75px;
        stroke-linejoin: round;
    }

    .leaflet-marker-icon svg path {
        stroke: #000000;
        stroke-width: 0.75px;
        stroke-linejoin: round;
    }

    .leaflet-marker-icon {
        transition: transform 1s ease;
        background: none;
        border: none;
    }

    .leaflet-popup-content-wrapper {
        background-color: rgb(var(--ui-background)) !important;
        color: #ffffff;
    }
    .leaflet-popup-content p {
        margin: 0.25em 0;
    }
    .leaflet-popup-tip {
        background-color: rgb(var(--ui-background)) !important;
    }

    .leaflet-control-layers-toggle {
        background-color: rgb(var(--color-primary-500)) !important;
    }
    .leaflet-control-layers {
        color: rgb(var(--color-primary-500));
        background-color: rgb(var(--ui-background)) !important;
    }

    .leaflet-control-attribution {
        color: rgb(var(--color-primary-500));
        background-color: rgb(var(--ui-background)) !important;
    }

    .leaflet-control-attribution a {
        color: rgb(var(--color-primary-500));
    }
    .leaflet-control-attribution a:hover {
        color: rgb(var(--color-primary-400));
    }

    /* Leaflet Contextmenu */
    .leaflet-contextmenu {
        display: none;
        box-shadow: 0 1px 7px rgba(0, 0, 0, 0.4);
        -webkit-border-radius: 2px;
        border-radius: 2px;
        padding: 4px 0;
        background-color: rgb(var(--ui-background));
        cursor: default;
        -webkit-user-select: none;
        -moz-user-select: none;
        user-select: none;
    }

    .leaflet-contextmenu a.leaflet-contextmenu-item {
        display: block;
        color: rgb(var(--color-primary-500));
        font-size: 12px;
        line-height: 20px;
        text-decoration: none;
        padding: 0 12px;
        cursor: default;
        outline: none;
    }

    .leaflet-contextmenu a.leaflet-contextmenu-item-disabled {
        opacity: 0.5;
    }

    .leaflet-contextmenu a.leaflet-contextmenu-item.over {
        background-color: rgb(var(--color-primary-100));
    }

    .leaflet-contextmenu a.leaflet-contextmenu-item-disabled.over {
        background-color: inherit;
        border-top: 1px solid transparent;
        border-bottom: 1px solid transparent;
    }

    .leaflet-contextmenu-icon {
        margin: 2px 8px 0 0;
        width: 16px;
        height: 16px;
        float: left;
        border: 0;
    }

    .leaflet-contextmenu-separator {
        border-bottom: 1px solid #ccc;
        margin: 5px 0;
    }

    /* Graticle */
    .leaflet-grid-label .gridlabel-vert {
        margin-left: 8px;
        -webkit-transform: rotate(90deg);
        transform: rotate(90deg);
    }

    .leaflet-grid-label .gridlabel-vert,
    .leaflet-grid-label .gridlabel-horiz {
        padding-left: 2px;
        text-shadow:
            -1px 0 #000000,
            0 1px #000000,
            1px 0 #000000,
            0 -1px #000000;
    }
}
</style>
