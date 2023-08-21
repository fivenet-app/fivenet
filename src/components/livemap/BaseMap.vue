<script lang="ts" setup>
import { LControl, LControlLayers, LMap, LTileLayer } from '@vue-leaflet/vue-leaflet';
import { useResizeObserver, watchDebounced } from '@vueuse/core';
import L from 'leaflet';
import 'leaflet-contextmenu';
import 'leaflet-contextmenu/dist/leaflet.contextmenu.min.css';
import 'leaflet/dist/leaflet.css';
import { useLivemapStore } from '~/store/livemap';
import { ValueOf } from '~/utils/types';

const { $loading } = useNuxtApp();
const route = useRoute();

defineProps<{
    mapOptions?: Record<string, any>;
}>();

const emits = defineEmits<{
    (e: 'mapReady', map: L.Map): void;
}>();

const livemapStore = useLivemapStore();
const { location, zoom } = storeToRefs(livemapStore);

let map: L.Map | undefined = undefined;

const mapContainer = ref<HTMLDivElement | null>(null);
useResizeObserver(mapContainer, (_) => {
    if (map === undefined) return;

    map.invalidateSize();
});

const centerX = 117.3;
const centerY = 172.8;
const scaleX = 0.02072;
const scaleY = 0.0205;

const customCRS = L.extend({}, L.CRS.Simple, {
    projection: L.Projection.LonLat,
    scale: function (zoom: number): number {
        return Math.pow(2, zoom);
    },
    zoom: function (sc: number): number {
        return Math.log(sc) / 0.6931471805599453;
    },
    distance: function (pos1: L.LatLng, pos2: L.LatLng): number {
        var x_difference = pos2.lng - pos1.lng;
        var y_difference = pos2.lat - pos1.lat;
        return Math.sqrt(x_difference * x_difference + y_difference * y_difference);
    },
    transformation: new L.Transformation(scaleX, centerX, -scaleY, centerY),
    infinite: true,
});

let center: L.PointExpression = [0, 0];
const attribution = '<a href="http://www.rockstargames.com/V/">Grand Theft Auto V</a>';
const selectedMarker = ref<bigint>();

const mouseLat = ref<string>((0).toFixed(3));
const mouseLong = ref<string>((0).toFixed(3));

const currentHash = ref<string>('');

watch(currentHash, () => window.location.replace(currentHash.value));

watch(location, () => {
    if (location.value === undefined) return;
    if (map === undefined) return;

    map?.setZoom(5, {
        animate: false,
    });
    map?.panTo([location.value?.y!, location.value?.x!], {
        animate: true,
        duration: 1.0,
    });
});

const isMoving = ref<boolean>(false);

watchDebounced(
    isMoving,
    () => {
        if (!map || isMoving.value) return;

        const newHash = stringifyHash(map.getZoom(), map.getCenter().lat, map.getCenter().lng);
        if (currentHash.value !== newHash) currentHash.value = newHash;
    },
    { debounce: 1000, maxWait: 3000 },
);

const backgroundColorList = {
    Satelite: '#143d6b',
    Postal: '#74aace',
} as const;
const backgroundColor = ref<ValueOf<typeof backgroundColorList>>(backgroundColorList.Postal);

async function updateBackground(layer: string): Promise<void> {
    switch (layer) {
        case 'Satelite':
            backgroundColor.value = backgroundColorList.Satelite;
            return;
        case 'Postal':
            backgroundColor.value = backgroundColorList.Postal;
            return;
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
    if (args.length !== 3) return;

    const zoom = parseInt(args[0], 10);
    const lat = parseFloat(args[1]);
    const lng = parseFloat(args[2]);

    if (isNaN(zoom) || isNaN(lat) || isNaN(lng)) return;

    return {
        latlng: new L.LatLng(lat, lng),
        zoom,
    };
}

async function onMapReady($event: any): Promise<void> {
    map = $event as L.Map;

    map.invalidateSize();

    const startingHash = route.hash;
    const startPos = parseHash(startingHash);
    if (startPos) map.setView(startPos.latlng, startPos.zoom);

    map.on('baselayerchange', async (event: L.LayersControlEvent) => {
        updateBackground(event.name);
    });

    map.addEventListener('mousemove', async (event: L.LeafletMouseEvent) => {
        if (!event.latlng) return;
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

    emits('mapReady', map);
}

onBeforeMount(() => {
    $loading.start();
});

onBeforeUnmount(() => {
    map = undefined;
});
</script>

<style>
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

<template>
    <div ref="mapContainer" class="h-full flex flex-row">
        <LMap
            class="z-0"
            v-model:zoom="zoom"
            v-model:center="center"
            :crs="customCRS"
            :min-zoom="1"
            :max-zoom="6"
            @click="selectedMarker = undefined"
            :inertia="false"
            :style="{ backgroundColor }"
            @ready="onMapReady($event)"
            :use-global-leaflet="true"
            :options="mapOptions"
        >
            <LTileLayer
                url="/images/livemap/tiles/postal/{z}/{x}/{y}.png"
                layer-type="base"
                name="Postal"
                :no-wrap="true"
                :tms="true"
                :visible="true"
                :attribution="attribution"
            />
            <LTileLayer
                url="/images/livemap/tiles/satelite/{z}/{x}/{y}.png"
                layer-type="base"
                name="Satelite"
                :no-wrap="true"
                :tms="true"
                :visible="false"
                :attribution="attribution"
            />

            <LControlLayers />

            <LControl position="bottomleft" class="leaflet-control-attribution mouseposition text-xs">
                {{ $t('common.longitude') }}: {{ mouseLat }} | {{ $t('common.latitude') }}: {{ mouseLong }}
            </LControl>

            <slot />
        </LMap>

        <slot name="afterMap" />
    </div>
</template>
