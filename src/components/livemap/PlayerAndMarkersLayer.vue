<script lang="ts" setup>
import { LControl, LLayerGroup } from '@vue-leaflet/vue-leaflet';
import { computedAsync } from '@vueuse/core';
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
const { markers, jobs } = storeToRefs(livemapStore);
const { startStream, stopStream } = livemapStore;

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const playerQuery = ref<string>('');
const playerMarkersFiltered = computedAsync(async () =>
    markers.value.users.filter((m) =>
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
    <LLayerGroup
        v-for="job in jobs.users"
        :key="`users_${job.name}`"
        :name="`${$t('common.employee', 2)} ${job.label}`"
        layer-type="overlay"
        :visible="true"
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
    </LLayerGroup>

    <LLayerGroup
        v-for="job in jobs.markers"
        :key="`markers_${job.name}`"
        :name="`${$t('common.marker', 2)} ${job.label}`"
        layer-type="overlay"
        :visible="true"
    >
        <MarkerMarker
            v-for="marker in markers.markers"
            :key="marker.info!.id"
            :marker="marker"
            :size="livemap.markerSize"
            @selected="$emit('markerSelected', marker)"
        />
    </LLayerGroup>

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
