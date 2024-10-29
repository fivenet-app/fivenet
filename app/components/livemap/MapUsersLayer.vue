<script lang="ts" setup>
import { LControl, LLayerGroup } from '@vue-leaflet/vue-leaflet';
import MapUserMarker from '~/components/livemap/MapUserMarker.vue';
import { useLivemapStore } from '~/store/livemap';
import { useSettingsStore } from '~/store/settings';
import type { UserMarker } from '~~/gen/ts/resources/livemap/livemap';

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
}>();

const livemapStore = useLivemapStore();
const { jobsUsers, markersUsers } = storeToRefs(livemapStore);
const { startStream, stopStream } = livemapStore;

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);

onBeforeMount(async () => {
    useTimeoutFn(async () => {
        try {
            startStream();
        } catch (e) {
            logger.error('exception during map users stream', e);
        }
    }, 50);
});

onBeforeUnmount(async () => stopStream());

const playerQueryRaw = ref<string>('');
const playerQuery = computed(() => playerQueryRaw.value.toLowerCase().trim());

const playerMarkersFiltered = computedAsync(async () =>
    [...(markersUsers.value.values() ?? [])].filter(
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
        <MapUserMarker
            v-for="marker in playerMarkersFiltered?.filter((m) => m.info?.job === job.name)"
            :key="marker.info!.id"
            :marker="marker"
            :size="livemap.markerSize"
            :show-unit-names="showUnitNames || livemap.showUnitNames"
            :show-unit-status="showUnitStatus || livemap.showUnitStatus"
            @selected="$emit('userSelected', marker)"
        />
    </LLayerGroup>

    <LControl v-if="showUserFilter" position="bottomleft">
        <div class="flex flex-col gap-2">
            <UInput
                v-model="playerQueryRaw"
                class="max-w-40"
                type="text"
                name="searchPlayer"
                size="xs"
                :placeholder="`${$t('common.employee', 1)} ${$t('common.filter')}`"
                autocomplete="off"
                :ui="{ icon: { trailing: { pointer: '' } } }"
                leading-icon="i-mdi-user-multiple"
            >
                <template #trailing>
                    <UButton
                        v-show="playerQueryRaw !== ''"
                        color="gray"
                        variant="link"
                        icon="i-mdi-close"
                        :padded="false"
                        @click="playerQueryRaw = ''"
                    />
                </template>
            </UInput>
        </div>
    </LControl>
</template>
