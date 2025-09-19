<script lang="ts" setup>
import DispatchDetailsSlideover from '~/components/centrum/dispatches/DispatchDetailsSlideover.vue';
import DispatchMarker from '~/components/centrum/livemap/DispatchMarker.vue';
import { useCentrumStore } from '~/stores/centrum';
import { useSettingsStore } from '~/stores/settings';
import type { Dispatch } from '~~/gen/ts/resources/centrum/dispatches';

const props = defineProps<{
    showAllDispatches?: boolean;
    dispatchList?: Dispatch[];
}>();

const { t } = useI18n();

const overlay = useOverlay();

const centrumStore = useCentrumStore();
const { dispatches, acls, ownDispatches, settings } = storeToRefs(centrumStore);

const settingsStore = useSettingsStore();
const { addOrUpdateLivemapLayer, addOrUpdateLivemapCategory, removeLivemapLayer } = settingsStore;
const { livemap, livemapLayers } = storeToRefs(settingsStore);

const dispatchQueryRaw = ref<string>('');
const dispatchQuery = computed(() => dispatchQueryRaw.value.trim().toLowerCase());

const activeJobLayers = computed(() =>
    livemapLayers.value
        .filter((l) => l.key.startsWith('dispatches_job_') && l.visible)
        .map((l) => l.key.replace('dispatches_job_', '')),
);

const dispatchesFiltered = computedAsync(async () =>
    [...(props.dispatchList ?? dispatches.value.values() ?? [])].filter(
        (m) =>
            !ownDispatches.value.includes(m.id) &&
            m.jobs?.jobs.some((j) => activeJobLayers.value.includes(j.name)) &&
            (m.id.toString().startsWith(dispatchQuery.value) ||
                m.message.toLowerCase().includes(dispatchQuery.value) ||
                (m.creator?.firstname + ' ' + m.creator?.lastname).toLowerCase().includes(dispatchQuery.value)),
    ),
);

watch(settings, () => {
    if (!settings.value?.enabled) return;

    removeLivemapLayer('dispatches_all');

    addOrUpdateLivemapLayer({
        key: 'dispatches_own',
        category: 'dispatches',
        label: t('common.your_dispatches'),
        perm: 'centrum.CentrumService/Stream',
        disabled: true,
        order: 1,
    });
});

watch(acls, () =>
    acls.value?.dispatches?.jobs.forEach((job) =>
        addOrUpdateLivemapLayer({
            key: `dispatches_job_${job.job}`,
            category: 'dispatches',
            label: job.jobLabel ?? job.job ?? '',
            perm: 'centrum.CentrumService/Stream',
            order: 10,
        }),
    ),
);

onBeforeMount(async () => {
    addOrUpdateLivemapCategory({
        key: 'dispatches',
        label: t('common.dispatch', 2),
        order: 0,
    });
});

const dispatchDetailsSlideover = overlay.create(DispatchDetailsSlideover);
</script>

<template>
    <LLayerGroup key="dispatches_own" :name="$t('common.your_dispatches')" layer-type="overlay" visible>
        <DispatchMarker
            v-for="dispatch in ownDispatches"
            :key="dispatch"
            :dispatch="dispatches.get(dispatch)!"
            :size="livemap.markerSize"
            @selected="
                dispatchDetailsSlideover.open({
                    dispatchId: $event.id,
                })
            "
        />
    </LLayerGroup>

    <LLayerGroup key="dispatches_all" :name="$t('common.dispatch', 2)" layer-type="overlay" visible>
        <DispatchMarker
            v-for="dispatch in dispatchesFiltered"
            :key="dispatch.id"
            :dispatch="dispatch"
            :size="livemap.markerSize"
            @selected="
                dispatchDetailsSlideover.open({
                    dispatchId: $event.id,
                })
            "
        />
    </LLayerGroup>

    <LControl position="bottomleft">
        <div v-if="settings?.enabled" class="flex flex-col gap-2">
            <UInput
                v-model="dispatchQueryRaw"
                class="max-w-40"
                type="text"
                name="searchPlayer"
                size="xs"
                :placeholder="`${$t('common.dispatch', 2)} ${$t('common.filter')}`"
                autocomplete="off"
                leading-icon="i-mdi-car-emergency"
                :ui="{ trailing: 'pe-1' }"
            >
                <template #trailing>
                    <UButton
                        v-if="dispatchQueryRaw !== ''"
                        color="neutral"
                        variant="link"
                        icon="i-mdi-close"
                        aria-controls="search"
                        @click="dispatchQueryRaw = ''"
                    />
                </template>
            </UInput>
        </div>
    </LControl>
</template>
