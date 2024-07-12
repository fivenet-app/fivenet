<script lang="ts" setup>
import { LLayerGroup } from '@vue-leaflet/vue-leaflet';
import MapMarkerMarker from '~/components/livemap/MapMarkerMarker.vue';
import { useLivemapStore } from '~/store/livemap';
import { useSettingsStore } from '~/store/settings';
import { type MarkerMarker } from '~~/gen/ts/resources/livemap/livemap';

defineEmits<{
    (e: 'markerSelected', marker: MarkerMarker): void;
}>();

const livemapStore = useLivemapStore();
const { jobsMarkers, markersMarkers } = storeToRefs(livemapStore);

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);
</script>

<template>
    <LLayerGroup
        v-for="job in jobsMarkers"
        :key="job.name"
        :name="`${$t('common.marker', 2)} ${job.label}`"
        layer-type="overlay"
        :visible="livemap.activeLayers.length === 0 || livemap.activeLayers.includes(`${$t('common.marker', 2)} ${job.label}`)"
    >
        <MapMarkerMarker
            v-for="marker in [...markersMarkers.values()].filter((p) => p.info?.job === job.name)"
            :key="marker.info!.id"
            :marker="marker"
            :size="livemap.markerSize"
            @selected="$emit('markerSelected', marker)"
        />
    </LLayerGroup>
</template>
