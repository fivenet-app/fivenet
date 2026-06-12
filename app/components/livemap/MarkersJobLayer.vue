<script lang="ts" setup>
import MapMarkerMarker from '~/components/livemap/MapMarkerMarker.vue';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';
import type { MarkerMarker } from '~~/gen/ts/resources/livemap/markers/marker_marker';

defineProps<{
    job: Job;
    markers: MarkerMarker[];
    visible?: boolean;
}>();

defineEmits<{
    (e: 'markerSelected', marker: MarkerMarker): void;
}>();

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);
</script>

<template>
    <LLayerGroup
        :key="job.name"
        :name="`${$t('common.marker', 2)} ${job.label}`"
        layer-type="overlay"
        :visible="visible"
        :options="{ name: `markers_${job.name}` }"
    >
        <MapMarkerMarker
            v-for="marker in markers"
            :key="marker.id"
            :marker="marker"
            :size="livemap.markerSize"
            :icon-key="`markers_${marker.id}_${marker.job}_${visible}`"
            @selected="$emit('markerSelected', marker)"
        />
    </LLayerGroup>
</template>
