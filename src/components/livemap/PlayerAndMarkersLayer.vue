<script lang="ts" setup>
import { LControl } from '@vue-leaflet/vue-leaflet';
import { computedAsync } from '@vueuse/core';
import { LMarkerClusterGroup } from 'vue-leaflet-markercluster';
import 'vue-leaflet-markercluster/dist/style.css';
import { useAuthStore } from '~/store/auth';
import { useLivemapStore } from '~/store/livemap';
import { useSettingsStore } from '~/store/settings';
import { Marker, UserMarker } from '~~/gen/ts/resources/livemap/livemap';
import MarkerMarker from '~/components/livemap/MarkerMarker.vue';
import PlayerMarker from '~/components/livemap/PlayerMarker.vue';

withDefaults(
    defineProps<{
        centerSelectedMarker?: boolean;
        filterPlayers?: boolean;
        showUnitNames?: boolean;
        showUnitStatus?: boolean;
    }>(),
    {
        centerSelectedMarker: true,
        filterPlayers: true,
    },
);

defineEmits<{
    (e: 'userSelected', marker: UserMarker): void;
    (e: 'markerSelected', marker: Marker): void;
}>();

const livemapStore = useLivemapStore();
const { jobsUsers, jobsMarkers, markersMarkers, markersUsers } = storeToRefs(livemapStore);
const { startStream, stopStream } = livemapStore;

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const playerQuery = ref<string>('');
const playerMarkersFiltered = computedAsync(async () =>
    markersUsers.value.filter((m) =>
        (m.user?.firstname + ' ' + m.user?.lastname).toLowerCase().includes(playerQuery.value.toLowerCase()),
    ),
);

onBeforeMount(async () => startStream());

onBeforeUnmount(async () => {
    stopStream();
    livemapStore.$reset();
});
</script>

<template>
    <LMarkerClusterGroup
        v-for="job in jobsUsers"
        :key="`users_${job.name}`"
        :name="`${$t('common.employee', 2)} ${job.label}`"
        layer-type="overlay"
        :visible="true"
        :max-cluster-radius="15"
        :disable-clustering-at-zoom="2"
        :chunked-loading="true"
        :animate="true"
        :animate-adding-markers="true"
    >
        <PlayerMarker
            v-for="marker in playerMarkersFiltered.filter((p) => p.user?.job === job.name)"
            :key="marker.userId"
            :marker="marker"
            :active-char="activeChar"
            :size="livemap.markerSize"
            :show-unit-names="showUnitNames || livemap.showUnitNames"
            :show-unit-status="showUnitStatus || livemap.showUnitStatus"
            @selected="$emit('userSelected', marker)"
        />
    </LMarkerClusterGroup>

    <LMarkerClusterGroup
        v-for="job in jobsMarkers"
        :key="`markers_${job.name}`"
        :name="`${$t('common.marker', 2)} ${job.label}`"
        layer-type="overlay"
        :visible="true"
        :max-cluster-radius="0"
        :disable-clustering-at-zoom="1"
        :chunked-loading="true"
        :animate="true"
        :animate-adding-markers="true"
    >
        <MarkerMarker
            v-for="marker in markersMarkers"
            :key="marker.info!.id"
            :marker="marker"
            :size="livemap.markerSize"
            @selected="$emit('markerSelected', marker)"
        />
    </LMarkerClusterGroup>

    <LControl v-if="filterPlayers" position="bottomleft">
        <div class="form-control flex flex-col gap-2">
            <input
                v-model="playerQuery"
                class="w-full max-w-[11rem] p-0.5 px-1 bg-clip-padding rounded-md border-2 border-black/20"
                type="text"
                name="searchPlayer"
                :placeholder="`${$t('common.employee', 1)} ${$t('common.filter')}`"
                @focusin="focusTablet(true)"
                @focusout="focusTablet(false)"
            />
        </div>
    </LControl>
</template>
