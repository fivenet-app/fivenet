<script lang="ts" setup>
import type { LatLngExpression, PointExpression } from 'leaflet';
import { MapMarkerDownIcon } from 'mdi-vue3';
import {
    customMapCRS,
    getMapBackgroundColor,
    mapBounds,
    mapMaxBounds,
    mapTileLayers,
} from '~/composables/livemap/useMapProjection';
import { overlayCayoPericoBounds, tileLayers } from '~/types/livemap';
import MapZoomControls from './controls/MapZoomControls.vue';
import TileLayerSelector from './controls/TileLayerSelector.vue';
import MapCayoPerico from './MapCayoPerico.vue';

const props = withDefaults(
    defineProps<{
        x: number;
        y: number;
        zoom: number;
        layer?: string;
        disabled?: boolean;
        showControls?: boolean;
    }>(),
    {
        layer: '',
        disabled: false,
        showControls: true,
    },
);

const emit = defineEmits<{
    (e: 'update:x', value: number): void;
    (e: 'update:y', value: number): void;
    (e: 'update:zoom', value: number): void;
    (e: 'update:layer', value: string): void;
}>();

const { game } = useAppConfig();

const settingsStore = useSettingsStore();
const { livemapTileLayer, livemap } = storeToRefs(settingsStore);

const mapRef = useTemplateRef('mapRef');
const mapContainer = useTemplateRef('mapContainer');
const mapResize = useDebounceFn(() => mapRef.value?.leafletObject?.invalidateSize(), 150, { maxWait: 400 });
const center = ref<LatLngExpression>([props.y, props.x]);
const currentZoom = ref(props.zoom);

const currentLayer = computed({
    get: () => props.layer ?? tileLayers[0]!.key,
    set: (value: string) => emit('update:layer', value),
});

const activeLayer = computed(() => mapTileLayers.find((layer) => layer.key === currentLayer.value) ?? tileLayers[0]!);
const iconAnchor = computed<PointExpression>(() => [livemap.value.markerSize / 2, livemap.value.markerSize]);
const iconSize = computed<[number, number]>(() => [livemap.value.markerSize, livemap.value.markerSize]);
const backgroundColor = computed(() => getMapBackgroundColor(activeLayer.value.key));

function scheduleResize(): void {
    nextTick(() => {
        requestAnimationFrame(() => mapResize());
    });
}

function syncMapView(): void {
    center.value = [props.y, props.x];
    currentZoom.value = props.zoom;
    mapRef.value?.leafletObject?.setView(center.value, currentZoom.value, { animate: false });
}

watch(
    () => [props.x, props.y, props.zoom, currentLayer.value],
    () => {
        syncMapView();
        scheduleResize();
    },
    { immediate: true },
);

useResizeObserver(mapContainer, () => mapResize());

function onMapReady(): void {
    syncMapView();
    scheduleResize();
}

function updatePosition(latlng: { lat: number; lng: number }): void {
    emit('update:x', latlng.lng);
    emit('update:y', latlng.lat);
}

function onZoomEnd(): void {
    const nextZoom = mapRef.value?.leafletObject?.getZoom();
    if (typeof nextZoom === 'number' && nextZoom !== props.zoom) {
        currentZoom.value = nextZoom;
    }
}

function onMapClick(event: { latlng?: { lat: number; lng: number } }): void {
    if (!event.latlng || props.disabled) return;

    updatePosition(event.latlng);
    center.value = [event.latlng.lat, event.latlng.lng];
    mapRef.value?.leafletObject?.setView(center.value, currentZoom.value, { animate: false });
    scheduleResize();
}

function onZoomUpdate(value: number): void {
    currentZoom.value = value;
}

watch(currentZoom, (value) => {
    if (value === props.zoom) return;
    emit('update:zoom', value);
});
</script>

<template>
    <ClientOnly>
        <div
            ref="mapContainer"
            class="relative h-64 w-full overflow-hidden rounded-lg border border-neutral-300 dark:border-neutral-700"
            :style="{ backgroundColor }"
        >
            <LMap
                ref="mapRef"
                class="absolute inset-0 block h-full w-full"
                :style="{ backgroundColor }"
                :bounds="mapBounds"
                :max-bounds="mapMaxBounds"
                :min-zoom="1"
                :max-zoom="7"
                :zoom="currentZoom"
                :center="center"
                :crs="customMapCRS"
                :inertia="false"
                :options="{ zoomControl: false }"
                use-global-leaflet
                @ready="onMapReady"
                @click="onMapClick"
                @update:zoom="onZoomUpdate"
                @moveend="onZoomEnd"
            >
                <LTileLayer
                    :url="activeLayer.url"
                    layer-type="base"
                    :name="$t(activeLayer.label)"
                    no-wrap
                    tms
                    :visible="true"
                    :min-zoom="1"
                    :max-zoom="activeLayer.options?.maxZoom || 7"
                    :attribution="activeLayer.options?.attribution || undefined"
                />

                <MapCayoPerico
                    v-if="game.livemap?.enableCayoPerico"
                    :tile-layer="livemapTileLayer"
                    :bounds="overlayCayoPericoBounds"
                />

                <LMarker :lat-lng="[props.y, props.x]">
                    <LIcon
                        :icon-size="iconSize"
                        :icon-anchor="iconAnchor"
                        class-name="pointer-events-none!"
                        :options="{ pmIgnore: true }"
                    >
                        <MapMarkerDownIcon class="size-full text-primary-500 drop-shadow-sm" />
                    </LIcon>

                    <LPopup :options="{ closeButton: false }">
                        <div class="text-xs">
                            <div class="font-medium">
                                {{ $t('common.coordinate') }}: {{ props.x.toFixed(2) }}, {{ props.y.toFixed(2) }}
                            </div>
                        </div>
                    </LPopup>
                </LMarker>

                <MapZoomControls v-if="showControls && !disabled" v-model="currentZoom" />
                <TileLayerSelector v-if="showControls && !disabled" v-model="currentLayer" />
            </LMap>
        </div>
    </ClientOnly>
</template>
