<script lang="ts" setup>
import { LControl, LControlLayers, LMap, LTileLayer } from '@vue-leaflet/vue-leaflet';
import type L from 'leaflet';
import { CRS, extend, LatLng, latLngBounds, Projection, Transformation, type PointExpression } from 'leaflet';
import 'leaflet-contextmenu';
import 'leaflet/dist/leaflet.css';
import ZoomControls from '~/components/livemap/controls/ZoomControls.vue';
import { useLivemapStore } from '~/store/livemap';
import type { ValueOf } from '~/utils/types';

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
const mapResizeDebounced = useDebounceFn(mapResize, 350, { maxWait: 750 });
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
const attribution = '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>' as const;

const mouseLat = ref<number>(0);
const mouseLong = ref<number>(0);

const currentHash = useRouteHash('');

function getZoomOffset(zoom: number): number {
    if (!slideover.isOpen.value) {
        return 0;
    }

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
    if (map === undefined || selectedMarker.value === undefined) {
        return;
    }

    map?.panTo([selectedMarker.value.info!.y!, selectedMarker.value.info!.x! + getZoomOffset(zoom.value)], {
        animate: true,
        duration: 0.75,
    });
});

watch(location, async () => {
    if (map === undefined || location.value === undefined) {
        return;
    }

    map?.panTo([location.value.y!, location.value.x! + getZoomOffset(zoom.value)], {
        animate: true,
        duration: 0.75,
    });
});

const isMoving = ref<boolean>(false);

watchDebounced(
    isMoving,
    async () => {
        if (!map || isMoving.value) {
            return;
        }

        const newHash = stringifyHash(map.getZoom(), map.getCenter().lat, map.getCenter().lng);
        if ((currentHash.value as string) !== newHash) {
            currentHash.value = newHash;
        }
    },
    { debounce: 1000, maxWait: 3000 },
);

const backgroundColorList = {
    Postal: '#74aace',
} as const;
const backgroundColor = ref<ValueOf<typeof backgroundColorList>>(backgroundColorList.Postal);

async function updateBackground(layer: string): Promise<void> {
    switch (layer) {
        case 'Postal':
        default:
            backgroundColor.value = backgroundColorList.Postal;
            break;
    }
}

function stringifyHash(currZoom: number, centerLat: number, centerLong: number): string {
    const precision = Math.max(0, Math.ceil(Math.log(zoom.value) / Math.LN2));

    const hash = '#' + [currZoom, centerLat.toFixed(precision), centerLong.toFixed(precision)].join('/');
    return hash;
}

function parseHash(hash: string): { latlng: L.LatLng; zoom: number } | undefined {
    if (hash.indexOf('#') === 0) {
        hash = hash.substring(1);
    }

    const args = hash.split('/');

    const zoom = args[0] ? parseInt(args[0]) : 2;
    const lat = args[1] ? parseFloat(args[1]) : 0;
    const lng = args[2] ? parseFloat(args[2]) : 0;

    if (isNaN(zoom) || isNaN(lat) || isNaN(lng)) {
        return;
    }

    return {
        latlng: new LatLng(lat, lng),
        zoom,
    };
}

async function onMapReady(map: L.Map): Promise<void> {
    map.invalidateSize();

    const startPos = parseHash(currentHash.value as string);
    if (startPos) {
        map.setView(startPos.latlng, startPos.zoom);
    }

    map.on('baselayerchange', async (event: L.LayersControlEvent) => {
        updateBackground(event.name);
    });

    map.on('overlayadd', (event) => {
        emit('overlayadd', event);
    });
    map.on('overlayremove', (event) => {
        emit('overlayremove', event);
    });

    map.addEventListener('mousemove', async (event: L.LeafletMouseEvent) => {
        if (!event.latlng) {
            return;
        }

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
}

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
            class="z-0"
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
            <ZoomControls />

            <LTileLayer
                url="/images/livemap/tiles/postal/{z}/{x}/{y}.png"
                layer-type="base"
                name="Postal"
                :no-wrap="true"
                :tms="true"
                :visible="true"
                :attribution="attribution"
                :min-zoom="1"
                :max-zoom="7"
            />

            <LControlLayers />

            <!-- eslint-disable-next-line tailwindcss/no-custom-classname -->
            <LControl position="bottomleft" class="leaflet-control-attribution">
                <span class="font-semibold">{{ $t('common.longitude') }}:</span> {{ mouseLat.toFixed(3) }} |
                <span class="font-semibold">{{ $t('common.latitude') }}:</span> {{ mouseLong.toFixed(3) }}
            </LControl>

            <slot />
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
        background-color: #16171a;
        color: #ffffff;
    }
    .leaflet-popup-content p {
        margin: 0.25em 0;
    }
    .leaflet-popup-tip {
        background-color: #16171a;
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
}
</style>
