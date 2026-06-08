<script lang="ts" setup>
import type { PointExpression } from 'leaflet';
import { MapMarkerDownIcon } from 'mdi-vue3';
import { mapTileLayers } from '~/composables/livemap/useMapProjection';
import { overlayCayoPericoBounds, tileLayers } from '~/types/livemap';
import LivemapMapShell from './LivemapMapShell.vue';
import MapZoomControls from './controls/MapZoomControls.vue';
import TileLayerSelector from './controls/TileLayerSelector.vue';

const props = withDefaults(
    defineProps<{
        x: number;
        y: number;
        zoom: number;
        layer?: string;
        disabled?: boolean;
        showControls?: boolean;
        containerClass?: string;
    }>(),
    {
        layer: '',
        disabled: false,
        showControls: true,
        containerClass: 'relative h-64 w-full overflow-hidden rounded-lg border border-neutral-300 dark:border-neutral-700',
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

const center = ref<PointExpression>([props.y, props.x]);
const currentZoom = ref(props.zoom);

const currentLayer = computed({
    get: () => props.layer ?? tileLayers[0]!.key,
    set: (value: string) => emit('update:layer', value),
});

const activeLayer = computed(() => mapTileLayers.find((layer) => layer.key === currentLayer.value) ?? tileLayers[0]!);
const iconAnchor = computed<PointExpression>(() => [livemap.value.markerSize / 2, livemap.value.markerSize]);
const iconSize = computed<[number, number]>(() => [livemap.value.markerSize, livemap.value.markerSize]);

watch(
    () => [props.x, props.y, props.zoom],
    () => {
        center.value = [props.y, props.x];
        currentZoom.value = props.zoom;
    },
    { immediate: true },
);

function updatePosition(latlng: { lat: number; lng: number }): void {
    emit('update:x', latlng.lng);
    emit('update:y', latlng.lat);
}

function onMapClick(event: { latlng?: { lat: number; lng: number } }): void {
    if (!event.latlng || props.disabled) return;

    updatePosition(event.latlng);
    center.value = [event.latlng.lat, event.latlng.lng];
}

watch(currentZoom, (value) => {
    if (value === props.zoom) return;
    emit('update:zoom', value);
});
</script>

<template>
    <ClientOnly>
        <LivemapMapShell
            v-model:center="center"
            v-model:zoom="currentZoom"
            :container-class="containerClass"
            map-class="absolute inset-0 block h-full w-full"
            :background-layer="activeLayer.key"
            :show-cayo-perico="game.livemap?.enableCayoPerico"
            :cayo-tile-layer="livemapTileLayer"
            :cayo-bounds="overlayCayoPericoBounds"
            @click="onMapClick"
        >
            <template #layers>
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
            </template>

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
        </LivemapMapShell>
    </ClientOnly>
</template>
