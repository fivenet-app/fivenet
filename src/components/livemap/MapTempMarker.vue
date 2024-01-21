<script lang="ts" setup>
import { LIcon, LMarker } from '@vue-leaflet/vue-leaflet';
import { useTimeoutFn } from '@vueuse/core';
import type { PointExpression } from 'leaflet';
import { MapMarkerDownIcon } from 'mdi-vue3';
import { useLivemapStore } from '~/store/livemap';
import { useSettingsStore } from '~/store/settings';

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);

const livemapStore = useLivemapStore();
const { location, showLocationMarker } = storeToRefs(livemapStore);

const iconAnchor: PointExpression = [livemap.value.markerSize / 2, livemap.value.markerSize];

const showMarker = ref(false);

const { start, stop } = useTimeoutFn(() => {
    if (!showLocationMarker.value) {
        showMarker.value = false;
    }
}, 6000);

watch(location, () => {
    showMarker.value = showLocationMarker.value || true;
    start();
});

onBeforeUnmount(() => {
    stop();
    if (!showLocationMarker.value) {
        showMarker.value = false;
    }
});
</script>

<template>
    <LMarker v-if="location && showMarker" :lat-lng="[location.y, location.x]" :z-index-offset="1000">
        <LIcon :icon-size="[livemap.markerSize, livemap.markerSize]" :icon-anchor="iconAnchor">
            <MapMarkerDownIcon class="h-full w-full animate-pulse fill-primary-500" />
        </LIcon>
    </LMarker>
</template>
