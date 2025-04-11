<script lang="ts" setup>
import { LControl, LLayerGroup } from '@vue-leaflet/vue-leaflet';
import DispatchDetailsSlideover from '~/components/centrum/dispatches/DispatchDetailsSlideover.vue';
import DispatchMarker from '~/components/centrum/livemap/DispatchMarker.vue';
import { useCentrumStore } from '~/stores/centrum';
import { useSettingsStore } from '~/stores/settings';

defineProps<{
    showAllDispatches?: boolean;
}>();

const { t } = useI18n();

const slideover = useSlideover();

const centrumStore = useCentrumStore();
const { dispatches, ownDispatches, settings } = storeToRefs(centrumStore);

const settingsStore = useSettingsStore();
const { addOrUpdateLivemapLayer, addOrUpdateLivemapCategory } = settingsStore;
const { livemap, livemapLayers } = storeToRefs(settingsStore);

const dispatchQueryRaw = ref<string>('');
const dispatchQuery = computed(() => dispatchQueryRaw.value.trim().toLowerCase());

const dispatchesFiltered = computedAsync(async () =>
    [...(dispatches.value.values() ?? [])].filter(
        (m) =>
            !ownDispatches.value.includes(m.id) &&
            (m.id.toString().startsWith(dispatchQuery.value) ||
                m.message.toLowerCase().includes(dispatchQuery.value) ||
                (m.creator?.firstname + ' ' + m.creator?.lastname).toLowerCase().includes(dispatchQuery.value)),
    ),
);

onBeforeMount(() => {
    addOrUpdateLivemapCategory({
        key: 'dispatches',
        label: t('common.dispatch', 2),
    });
    addOrUpdateLivemapLayer({
        key: 'dispatches_own',
        category: 'dispatches',
        label: t('common.your_dispatches'),
        perm: 'CentrumService.Stream',
    });
    addOrUpdateLivemapLayer({
        key: 'dispatches_all',
        category: 'dispatches',
        label: t('common.all_dispatches'),
        perm: 'CentrumService.Stream',
    });
});
</script>

<template>
    <LLayerGroup key="dispatches_own" :name="$t('common.your_dispatches')" layer-type="overlay" :visible="true">
        <DispatchMarker
            v-for="dispatch in ownDispatches"
            :key="dispatch"
            :dispatch="dispatches.get(dispatch)!"
            :size="livemap.markerSize"
            @selected="
                slideover.open(DispatchDetailsSlideover, {
                    dispatchId: $event.id,
                })
            "
        />
    </LLayerGroup>

    <LLayerGroup
        key="dispatches_all"
        :name="$t('common.dispatch', 2)"
        layer-type="overlay"
        :visible="
            showAllDispatches || livemapLayers.find((l) => l.key === `dispatches_all`)?.visible || dispatchQuery.length > 0
        "
    >
        <DispatchMarker
            v-for="dispatch in dispatchesFiltered"
            :key="dispatch.id"
            :dispatch="dispatch"
            :size="livemap.markerSize"
            @selected="
                slideover.open(DispatchDetailsSlideover, {
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
                :ui="{ icon: { trailing: { pointer: '' } } }"
                leading-icon="i-mdi-car-emergency"
            >
                <template #trailing>
                    <UButton
                        v-show="dispatchQueryRaw !== ''"
                        color="gray"
                        variant="link"
                        icon="i-mdi-close"
                        :padded="false"
                        @click="dispatchQueryRaw = ''"
                    />
                </template>
            </UInput>
        </div>
    </LControl>
</template>
