<script lang="ts" setup>
import { LMarkerClusterGroup } from 'vue-leaflet-markercluster';
import DispatchDetails from '~/components/centrum/dispatches/DispatchDetails.vue';
import DispatchMarker from '~/components/centrum/livemap/DispatchMarker.vue';
import { useCentrumStore } from '~/store/centrum';
import { useSettingsStore } from '~/store/settings';
import { Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';

defineProps<{
    showAllDispatches?: boolean;
}>();

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const centrumStore = useCentrumStore();
const { dispatches, ownDispatches } = storeToRefs(centrumStore);
const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);

const selectedDispatch = ref<Dispatch | undefined>();
const open = ref(false);
</script>

<template>
    <template v-if="selectedDispatch">
        <DispatchDetails :dispatch="selectedDispatch" :open="open" @close="open = false" @goto="$emit('goto', $event)" />
    </template>

    <LMarkerClusterGroup
        key="your_dispatches"
        :name="$t('common.your_dispatches')"
        layer-type="overlay"
        :visible="true"
        :max-cluster-radius="0"
        :disable-clustering-at-zoom="1"
        :chunked-loading="true"
        :animate="true"
    >
        <DispatchMarker
            v-for="dispatch in ownDispatches"
            :key="dispatch"
            :dispatch="dispatches.get(dispatch)!"
            :size="livemap.markerSize"
            @selected="
                selectedDispatch = $event;
                open = true;
            "
        />
    </LMarkerClusterGroup>

    <LMarkerClusterGroup
        key="all_dispatches"
        :name="$t('common.dispatch', 2)"
        layer-type="overlay"
        :visible="showAllDispatches"
        :max-cluster-radius="15"
        :disable-clustering-at-zoom="2"
        :chunked-loading="true"
        :animate="true"
    >
        <DispatchMarker
            v-for="[id, dispatch] in dispatches"
            :key="id"
            :dispatch="dispatch"
            :size="livemap.markerSize"
            @selected="
                selectedDispatch = $event;
                open = true;
            "
        />
    </LMarkerClusterGroup>
</template>
