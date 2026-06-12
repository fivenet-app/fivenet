<script lang="ts" setup>
import MapUserMarker from '~/components/livemap/MapUserMarker.vue';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';
import type { UserMarker } from '~~/gen/ts/resources/livemap/markers/user_marker';

defineProps<{
    job: Job;
    markers: UserMarker[];
    visible?: boolean;
    showUnitNames?: boolean;
    showUnitStatus?: boolean;
    useUnitColor?: boolean;
}>();

defineEmits<{
    (e: 'userSelected', marker: UserMarker): void;
}>();

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);
</script>

<template>
    <LLayerGroup
        :key="job.name"
        :name="`${$t('common.employee', 2)} ${job.label}`"
        layer-type="overlay"
        :visible="visible"
        :options="{ name: `users_${job.name}` }"
    >
        <MapUserMarker
            v-for="marker in markers"
            :key="marker.userId"
            :marker="marker"
            :size="livemap.markerSize"
            :show-unit-names="showUnitNames || livemap.showUnitNames"
            :show-unit-status="showUnitStatus || livemap.showUnitStatus"
            :use-unit-color="useUnitColor || livemap.useUnitColor"
            @selected="$emit('userSelected', marker)"
        />
    </LLayerGroup>
</template>
