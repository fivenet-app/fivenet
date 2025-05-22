<script lang="ts" setup>
import { LLayerGroup } from '@vue-leaflet/vue-leaflet';
import MapMarkerMarker from '~/components/livemap/MapMarkerMarker.vue';
import { useLivemapStore } from '~/stores/livemap';
import { useSettingsStore } from '~/stores/settings';
import type { MarkerMarker } from '~~/gen/ts/resources/livemap/livemap';

defineEmits<{
    (e: 'markerSelected', marker: MarkerMarker): void;
}>();

const { t } = useI18n();

const livemapStore = useLivemapStore();
const { jobsMarkers, markersMarkers } = storeToRefs(livemapStore);

const settingsStore = useSettingsStore();
const { addOrUpdateLivemapLayer, addOrUpdateLivemapCategory } = settingsStore;
const { livemap, livemapLayers } = storeToRefs(settingsStore);

watch(jobsMarkers, () =>
    jobsMarkers.value.forEach((job) =>
        addOrUpdateLivemapLayer({
            key: `markers_${job.name}`,
            category: 'markers',
            label: job.label,
            perm: 'LivemapperService.Stream',
            attr: {
                key: 'Markers',
                val: job.name,
            },
        }),
    ),
);

onBeforeMount(async () =>
    addOrUpdateLivemapCategory({
        key: 'markers',
        label: t('common.marker', 2),
        order: 3,
    }),
);
</script>

<template>
    <LLayerGroup
        v-for="job in jobsMarkers"
        :key="job.name"
        :name="`${$t('common.marker', 2)} ${job.label}`"
        layer-type="overlay"
        :visible="livemapLayers.find((l) => l.key === `markers_${job.name}`)?.visible"
    >
        <MapMarkerMarker
            v-for="marker in [...markersMarkers.values()].filter((p) => p.job === job.name)"
            :key="marker.id"
            :marker="marker"
            :size="livemap.markerSize"
            @selected="$emit('markerSelected', marker)"
        />
    </LLayerGroup>
</template>
