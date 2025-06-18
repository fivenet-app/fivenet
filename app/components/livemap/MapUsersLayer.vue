<script lang="ts" setup>
import MapUserMarker from '~/components/livemap/MapUserMarker.vue';
import { useLivemapStore } from '~/stores/livemap';
import { useSettingsStore } from '~/stores/settings';
import type { UserMarker } from '~~/gen/ts/resources/livemap/user_marker';

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

const { t } = useI18n();

const livemapStore = useLivemapStore();
const { jobsUsers, markersUsers, ownMarker } = storeToRefs(livemapStore);
const { startStream, stopStream, goto } = livemapStore;

const settingsStore = useSettingsStore();
const { addOrUpdateLivemapLayer, addOrUpdateLivemapCategory } = settingsStore;
const { livemap, livemapLayers } = storeToRefs(settingsStore);

watch(jobsUsers, () =>
    jobsUsers.value.forEach((job) =>
        addOrUpdateLivemapLayer({
            key: `users_${job.name}`,
            category: 'users',
            label: job.label,
            perm: 'livemap.LivemapService/Stream',
            attr: {
                key: 'Players',
                val: job.name,
            },
        }),
    ),
);

onBeforeMount(async () => {
    addOrUpdateLivemapCategory({
        key: 'users',
        label: t('common.employee', 2),
        order: 2,
    });

    useTimeoutFn(async () => {
        try {
            startStream();
        } catch (e) {
            logger.error('exception during map users stream', e);
        }
    }, 50);
});

onBeforeRouteLeave(async (to) => {
    // Don't end livemap stream if user is switching to livemap/centrum page
    if (to.path.startsWith('/livemap') || to.path === '/centrum') return;

    await stopStream();
});

const playerQueryRaw = ref<string>('');
const playerQuery = computed(() => playerQueryRaw.value.toLowerCase());

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
        :visible="livemapLayers.find((l) => l.key === `users_${job.name}`)?.visible === true"
    >
        <MapUserMarker
            v-for="marker in playerMarkersFiltered?.filter((m) => m.job === job.name)"
            :key="marker.userId"
            :marker="marker"
            :size="livemap.markerSize"
            :show-unit-names="showUnitNames || livemap.showUnitNames"
            :show-unit-status="showUnitStatus || livemap.showUnitStatus"
            @selected="$emit('userSelected', marker)"
        />
    </LLayerGroup>

    <LControl position="topleft">
        <UTooltip v-if="ownMarker" :text="$t('common.my_location')" :popper="{ placement: 'right' }">
            <UButton
                class="inset-0 border border-black/20 bg-clip-padding p-1.5 hover:bg-[#f4f4f4]"
                icon="i-mdi-my-location"
                block
                @click="async () => await goto({ x: ownMarker!.x, y: ownMarker!.y }, false)"
            />
        </UTooltip>
    </LControl>

    <LControl v-if="showUserFilter" position="bottomleft">
        <div class="flex flex-col gap-2">
            <UInput
                v-model.trim="playerQueryRaw"
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
