<script lang="ts" setup>
import { LIcon, LMarker } from '@vue-leaflet/vue-leaflet';
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

const { start } = useTimeoutFn(() => (showMarker.value = false), 5000, { immediate: false });

watch(location, () => {
    showMarker.value = showLocationMarker.value;
    start();
});

watch(showLocationMarker, () => {
    if (!showLocationMarker.value) {
        showMarker.value = false;
    }
});
</script>

<template>
    <LMarker v-if="location && showMarker" :lat-lng="[location.y, location.x]" :z-index-offset="90">
        <LIcon
            :icon-size="[livemap.markerSize, livemap.markerSize]"
            :icon-anchor="iconAnchor"
            class-name="!pointer-events-none"
        >
            <MapMarkerDownIcon class="text-primary-500 size-full animate-pulse" />
        </LIcon>
    </LMarker>
</template>
