<script lang="ts" setup>
import type { LeafletMouseEvent, Map, PointExpression, LatLngBoundsExpression } from 'leaflet';
import { customMapCRS, getMapBackgroundColor, mapBounds, mapMaxBounds } from '~/composables/livemap/useMapProjection';
import MapCayoPerico from './MapCayoPerico.vue';

const props = withDefaults(
    defineProps<{
        backgroundLayer: string;
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        mapOptions?: Record<string, any>;
        bounds?: LatLngBoundsExpression;
        maxBounds?: LatLngBoundsExpression;
        minZoom?: number;
        maxZoom?: number;
        containerClass?: string;
        mapClass?: string;
        showCayoPerico?: boolean;
        cayoTileLayer?: string;
        cayoBounds?: LatLngBoundsExpression;
    }>(),
    {
        mapOptions: () => ({ zoomControl: false }),
        bounds: () => mapBounds,
        maxBounds: () => mapMaxBounds,
        minZoom: 1,
        maxZoom: 7,
        containerClass: 'relative h-64 w-full overflow-hidden rounded-lg border border-neutral-300 dark:border-neutral-700',
        mapClass: 'absolute inset-0 block h-full w-full',
        showCayoPerico: false,
        cayoTileLayer: '',
        cayoBounds: undefined,
    },
);

const center = defineModel<PointExpression>('center', { required: true });
const zoom = defineModel<number>('zoom', { required: true });

const emit = defineEmits<{
    (e: 'ready', map: Map): void;
    (e: 'click', event: LeafletMouseEvent): void;
    (e: 'movestart'): void;
    (e: 'moveend'): void;
}>();

const mapRef = useTemplateRef('mapRef');
const mapContainer = useTemplateRef('mapContainer');
const mapResize = useDebounceFn(() => mapRef.value?.leafletObject?.invalidateSize(), 150, { maxWait: 400 });
const backgroundColor = computed(() => getMapBackgroundColor(props.backgroundLayer));

function scheduleResize(): void {
    nextTick(() => {
        requestAnimationFrame(() => mapResize());
    });
}

function onMapReady(map: Map): void {
    scheduleResize();
    emit('ready', map);
}

function onMapClick(event: LeafletMouseEvent): void {
    emit('click', event);
}

useResizeObserver(mapContainer, () => mapResize());

defineExpose({
    mapResize,
});
</script>

<template>
    <div ref="mapContainer" :class="props.containerClass" :style="{ backgroundColor }">
        <LMap
            ref="mapRef"
            v-model:zoom="zoom"
            v-model:center="center"
            :style="{ backgroundColor }"
            :bounds="props.bounds"
            :max-bounds="props.maxBounds"
            :min-zoom="props.minZoom"
            :max-zoom="props.maxZoom"
            :crs="customMapCRS"
            :inertia="false"
            :options="props.mapOptions"
            use-global-leaflet
            @ready="onMapReady"
            @click="onMapClick"
            @movestart="emit('movestart')"
            @moveend="emit('moveend')"
        >
            <slot name="layers" />

            <MapCayoPerico
                v-if="props.showCayoPerico && props.cayoBounds"
                :tile-layer="props.cayoTileLayer"
                :bounds="props.cayoBounds"
            />

            <slot />
        </LMap>
    </div>
</template>
