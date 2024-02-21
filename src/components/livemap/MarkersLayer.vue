<script lang="ts" setup>
import { LMarkerClusterGroup } from 'vue-leaflet-markercluster';
import 'vue-leaflet-markercluster/dist/style.css';
import { useLivemapStore } from '~/store/livemap';
import { useSettingsStore } from '~/store/settings';
import { MarkerMarker as Marker } from '~~/gen/ts/resources/livemap/livemap';
import MarkerMarker from '~/components/livemap/MarkerMarker.vue';

defineEmits<{
    (e: 'markerSelected', marker: Marker): void;
    (e: 'goto', loc: Coordinate): void;
}>();

const livemapStore = useLivemapStore();
const { jobsMarkers, markersMarkers } = storeToRefs(livemapStore);

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);
</script>

<template>
    <LMarkerClusterGroup
        v-for="job in jobsMarkers"
        :key="job.name"
        :name="`${$t('common.marker', 2)} ${job.label}`"
        layer-type="overlay"
        :visible="livemap.activeLayers.length === 0 || livemap.activeLayers.includes(`${$t('common.marker', 2)} ${job.label}`)"
        :max-cluster-radius="0"
        :disable-clustering-at-zoom="0"
        :single-marker-mode="false"
        :chunked-loading="true"
        :animate="false"
    >
        <MarkerMarker
            v-for="marker in [...markersMarkers.values()].filter((p) => p.info?.job === job.name)"
            :key="marker.info!.id"
            :marker="marker"
            :size="livemap.markerSize"
            @selected="$emit('markerSelected', marker)"
            @goto="$emit('goto', $event)"
        />
    </LMarkerClusterGroup>
</template>
