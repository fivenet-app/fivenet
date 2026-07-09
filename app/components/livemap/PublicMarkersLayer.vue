<script lang="ts" setup>
import MapMarkerMarker from '~/components/livemap/MapMarkerMarker.vue';
import type { MarkerMarker } from '~~/gen/ts/resources/livemap/markers/marker_marker';

withDefaults(
    defineProps<{
        markers: MarkerMarker[];
        visible?: boolean;
    }>(),
    {
        visible: true,
    },
);

defineEmits<{
    (e: 'markerSelected', marker: MarkerMarker): void;
}>();

const { t } = useI18n();

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);
</script>

<template>
    <LLayerGroup :name="t('common.public')" layer-type="overlay" :visible="visible" :options="{ name: 'markers_public' }">
        <MapMarkerMarker
            v-for="marker in markers"
            :key="marker.id"
            :marker="marker"
            :size="livemap.markerSize"
            :icon-key="`markers_${marker.id}_${marker.job}_${marker.public}`"
            @selected="$emit('markerSelected', marker)"
        />
    </LLayerGroup>
</template>
