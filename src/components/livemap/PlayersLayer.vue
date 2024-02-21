<script lang="ts" setup>
import { LControl } from '@vue-leaflet/vue-leaflet';
import { computedAsync, useTimeoutFn } from '@vueuse/core';
import { LMarkerClusterGroup } from 'vue-leaflet-markercluster';
import 'vue-leaflet-markercluster/dist/style.css';
import { useAuthStore } from '~/store/auth';
import { useLivemapStore } from '~/store/livemap';
import { useSettingsStore } from '~/store/settings';
import { UserMarker } from '~~/gen/ts/resources/livemap/livemap';
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
    (e: 'goto', loc: Coordinate): void;
}>();

const livemapStore = useLivemapStore();
const { jobsUsers, markersUsers } = storeToRefs(livemapStore);
const { startStream, stopStream } = livemapStore;

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const playerQuery = ref<string>('');
const playerMarkersFiltered = computedAsync(async () =>
    [...markersUsers.value.values()].filter(
        (m) =>
            playerQuery.value === '' ||
            (m.user?.firstname + ' ' + m.user?.lastname).toLowerCase().includes(playerQuery.value.toLowerCase()),
    ),
);

const { start, stop } = useTimeoutFn(async () => startStream(), 650);

onBeforeMount(async () => start());

onBeforeUnmount(async () => {
    stop();
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
        :visible="
            livemap.activeLayers.length === 0 || livemap.activeLayers.includes(`${$t('common.employee', 2)} ${job.label}`)
        "
        :max-cluster-radius="15"
        :disable-clustering-at-zoom="2"
        :single-marker-mode="true"
        :chunked-loading="true"
        :animate="true"
    >
        <PlayerMarker
            v-for="marker in playerMarkersFiltered.filter((p) => p.user?.job === job.name)"
            :key="`user_${marker.userId}`"
            :marker="marker"
            :active-char="activeChar"
            :size="livemap.markerSize"
            :show-unit-names="showUnitNames || livemap.showUnitNames"
            :show-unit-status="showUnitStatus || livemap.showUnitStatus"
            @selected="$emit('userSelected', marker)"
            @goto="$emit('goto', $event)"
        />
    </LMarkerClusterGroup>

    <LControl v-if="filterPlayers" position="bottomleft">
        <div class="form-control flex flex-col gap-2">
            <input
                v-model="playerQuery"
                class="w-full max-w-[11rem] rounded-md border-2 border-black/20 bg-clip-padding p-0.5 px-1"
                type="text"
                name="searchPlayer"
                :placeholder="`${$t('common.employee', 1)} ${$t('common.filter')}`"
                @focusin="focusTablet(true)"
                @focusout="focusTablet(false)"
            />
        </div>
    </LControl>
</template>
