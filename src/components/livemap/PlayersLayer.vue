<script lang="ts" setup>
import { LControl, LLayerGroup } from '@vue-leaflet/vue-leaflet';
import { computedAsync, useTimeoutFn } from '@vueuse/core';
import { useLivemapStore } from '~/store/livemap';
import { useSettingsStore } from '~/store/settings';
import { UserMarker } from '~~/gen/ts/resources/livemap/livemap';
import PlayerMarker from '~/components/livemap/PlayerMarker.vue';

withDefaults(
    defineProps<{
        centerSelectedMarker?: boolean;
        showUserFilter?: boolean;
        showUnitNames?: boolean;
        showUnitStatus?: boolean;
    }>(),
    {
        centerSelectedMarker: true,
        filterPlayers: true,
        showUserFilter: true,
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

const { start, stop } = useTimeoutFn(async () => startStream(), 150);

onBeforeMount(async () => start());

onBeforeUnmount(async () => {
    stop();
    stopStream();
    livemapStore.$reset();
});

const playerQueryRaw = ref<string>('');
const playerQuery = computed(() => playerQueryRaw.value.toLowerCase().trim());

const playerMarkersFiltered = computedAsync(async () =>
    [...markersUsers.value.values()].filter(
        (m) =>
            playerQuery.value === '' || (m.user?.firstname + ' ' + m.user?.lastname).toLowerCase().includes(playerQuery.value),
    ),
);
</script>

<template>
    <LLayerGroup
        v-for="job in jobsUsers"
        :key="job.name"
        :name="`${$t('common.employee', 2)} ${job.label}`"
        layer-type="overlay"
        :visible="
            livemap.activeLayers.length === 0 || livemap.activeLayers.includes(`${$t('common.employee', 2)} ${job.label}`)
        "
    >
        <PlayerMarker
            v-for="marker in playerMarkersFiltered.filter((m) => m.info?.job === job.name)"
            :key="marker.info!.id"
            :marker="marker"
            :size="livemap.markerSize"
            :show-unit-names="showUnitNames || livemap.showUnitNames"
            :show-unit-status="showUnitStatus || livemap.showUnitStatus"
            @selected="$emit('userSelected', marker)"
            @goto="$emit('goto', $event)"
        />
    </LLayerGroup>

    <LControl v-if="showUserFilter" position="bottomleft">
        <div class="flex flex-col gap-2">
            <input
                v-model="playerQueryRaw"
                class="w-full max-w-44 rounded-md border-2 border-black/20 bg-clip-padding p-0.5 px-1"
                type="text"
                name="searchPlayer"
                :placeholder="`${$t('common.employee', 1)} ${$t('common.filter')}`"
                @focusin="focusTablet(true)"
                @focusout="focusTablet(false)"
            />
        </div>
    </LControl>
</template>
