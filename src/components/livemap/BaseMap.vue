<script lang="ts" setup>
import { LControl, LControlLayers, LMap, LTileLayer } from '@vue-leaflet/vue-leaflet';
import { useDebounceFn, useResizeObserver, watchDebounced } from '@vueuse/core';
import L, { extend, latLngBounds, CRS, LatLng, Projection, Transformation, type PointExpression } from 'leaflet';
import 'leaflet/dist/leaflet.css';
import 'leaflet-contextmenu';
import 'leaflet-contextmenu/dist/leaflet.contextmenu.min.css';
import { useLivemapStore } from '~/store/livemap';
import { type ValueOf } from '~/utils/types';

const { $loading } = useNuxtApp();
const route = useRoute();

defineProps<{
    mapOptions?: Record<string, any>;
}>();

const emit = defineEmits<{
    (e: 'mapReady', map: L.Map): void;
}>();

const livemapStore = useLivemapStore();
const { location, zoom } = storeToRefs(livemapStore);

let map: L.Map | undefined;

function mapResize(): void {
    if (map === undefined) {
        return;
    }

    map.invalidateSize();
}

const mapContainer = ref<HTMLElement | null>(null);
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
const attribution = '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>';
const selectedMarker = ref<string>();

const mouseLat = ref<string>((0).toFixed(3));
const mouseLong = ref<string>((0).toFixed(3));

const currentHash = ref<string>('');

watch(currentHash, () => window.location.replace(currentHash.value));

watch(location, () => {
    if (location.value === undefined || map === undefined) {
        return;
    }

    map?.panTo([location.value.y!, location.value.x!], {
        animate: true,
        duration: 0.85,
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
        if (currentHash.value !== newHash) currentHash.value = newHash;
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
    if (hash.indexOf('#') === 0) hash = hash.substring(1);

    const args = hash.split('/');
    if (args.length !== 3) {
        return;
    }

    const zoom = parseInt(args[0], 10);
    const lat = parseFloat(args[1]);
    const lng = parseFloat(args[2]);

    if (isNaN(zoom) || isNaN(lat) || isNaN(lng)) {
        return;
    }

    return {
        latlng: new LatLng(lat, lng),
        zoom,
    };
}

async function onMapReady($event: any): Promise<void> {
    map = $event as L.Map;

    map.invalidateSize();

    const startingHash = route.hash;
    const startPos = parseHash(startingHash);
    if (startPos) {
        map.setView(startPos.latlng, startPos.zoom);
    }

    map.on('baselayerchange', async (event: L.LayersControlEvent) => {
        updateBackground(event.name);
    });

    map.addEventListener('mousemove', async (event: L.LeafletMouseEvent) => {
        if (!event.latlng) {
            return;
        }

        mouseLat.value = (Math.round(event.latlng.lat * 100000) / 100000).toFixed(3);
        mouseLong.value = (Math.round(event.latlng.lng * 100000) / 100000).toFixed(3);
    });

    map.on('movestart', async () => {
        isMoving.value = true;
    });

    map.on('moveend', async () => {
        isMoving.value = false;
    });

    setTimeout(() => {
        $loading.finish();
    }, 500);

    emit('mapReady', map);
}

onBeforeMount(() => {
    $loading.start();
});

onBeforeUnmount(() => {
    map = undefined;
});
</script>

<template>
    <div ref="mapContainer" class="h-full flex flex-row" :style="{ backgroundColor }">
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

            <LControlLayers :hide-single-base="true" />

            <LControl position="bottomleft" class="leaflet-control-attribution mouseposition text-xs">
                <span class="font-semibold">{{ $t('common.longitude') }}</span
                >: {{ mouseLat }} | <span class="font-semibold">{{ $t('common.latitude') }}</span
                >: {{ mouseLong }}
            </LControl>

            <slot />
        </LMap>

        <slot name="afterMap" />
    </div>
</template>

<style lang="scss">
.leaflet-container {
    font-family: var(--font-sans);
}

.leaflet-container a {
    color: #7161ef; // text-primary-500
}
.leaflet-container a:hover {
    color: #8d81f2; // text-primary-400
}

.leaflet-div-icon {
    background: none;
    border: none;
}

.leaflet-div-icon svg path {
    stroke: #000;
    stroke-width: 0.75px;
    stroke-linejoin: round;
}

.leaflet-marker-icon {
    transition: transform 1s ease;
}

.leaflet-popup-content-wrapper {
    background-color: #16171a;
    color: #fff;
}
.leaflet-popup-content p {
    margin: 0.25em 0;
}
.leaflet-popup-tip {
    background-color: #16171a;
}
</style>
