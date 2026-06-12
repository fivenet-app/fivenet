<script lang="ts" setup>
import { useLivemapStore } from '~/stores/livemap';
import { useSettingsStore } from '~/stores/settings';
import { groupByJob } from '~/utils/livemap/groupByJob';
import type { MarkerMarker } from '~~/gen/ts/resources/livemap/markers/marker_marker';
import MarkerSearchDrawer from './markers/MarkerSearchDrawer.vue';
import MarkersJobLayer from '~/components/livemap/MarkersJobLayer.vue';

defineEmits<{
    (e: 'markerSelected', marker: MarkerMarker): void;
}>();

const { t } = useI18n();

const overlay = useOverlay();

const livemapStore = useLivemapStore();
const { jobsMarkers, markersMarkers } = storeToRefs(livemapStore);

const settingsStore = useSettingsStore();
const { addOrUpdateLivemapLayer, addOrUpdateLivemapCategory } = settingsStore;
const { livemapLayers } = storeToRefs(settingsStore);

watch(jobsMarkers, (val) =>
    val.forEach((job) =>
        addOrUpdateLivemapLayer({
            key: `markers_${job.name}`,
            category: 'markers',
            label: job.label,
            perm: 'livemap.LivemapService/Stream',
            attr: {
                key: 'Markers',
                val: job.name,
            },
        }),
    ),
);

const markersByJob = computed(() => groupByJob<MarkerMarker>(markersMarkers.value.values()));

onBeforeMount(async () =>
    addOrUpdateLivemapCategory({
        key: 'markers',
        label: t('common.marker', 2),
        order: 3,
    }),
);

const markerSearchDrawer = overlay.create(MarkerSearchDrawer);
</script>

<template>
    <MarkersJobLayer
        v-for="job in jobsMarkers"
        :key="job.name"
        :job="job"
        :markers="markersByJob.get(job.name) ?? []"
        :visible="livemapLayers.find((l) => l.key === `markers_${job.name}`)?.visible === true"
        @marker-selected="$emit('markerSelected', $event)"
    />

    <LControl position="topright">
        <div class="flex flex-col gap-2">
            <UTooltip :text="$t('common.marker', 2)">
                <UButton icon="i-mdi-map-markers" @click="markerSearchDrawer.open()" />
            </UTooltip>
        </div>
    </LControl>
</template>
