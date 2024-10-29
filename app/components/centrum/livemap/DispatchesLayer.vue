<script lang="ts" setup>
import { LControl, LLayerGroup } from '@vue-leaflet/vue-leaflet';
import DispatchDetailsSlideover from '~/components/centrum/dispatches/DispatchDetailsSlideover.vue';
import DispatchMarker from '~/components/centrum/livemap/DispatchMarker.vue';
import { useCentrumStore } from '~/store/centrum';
import { useSettingsStore } from '~/store/settings';

defineProps<{
    showAllDispatches?: boolean;
}>();

const slideover = useSlideover();

const centrumStore = useCentrumStore();
const { dispatches, ownDispatches } = storeToRefs(centrumStore);

const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);

const dispatchQueryRaw = ref<string>('');
const dispatchQuery = computed(() => dispatchQueryRaw.value.trim().toLowerCase());

const dispatchesFiltered = computedAsync(async () =>
    [...(dispatches.value.values() ?? [])].filter(
        (m) =>
            !ownDispatches.value.includes(m.id) &&
            (m.id.startsWith(dispatchQuery.value) ||
                m.message.toLowerCase().includes(dispatchQuery.value) ||
                (m.creator?.firstname + ' ' + m.creator?.lastname).toLowerCase().includes(dispatchQuery.value)),
    ),
);
</script>

<template>
    <LLayerGroup key="your_dispatches" :name="$t('common.your_dispatches')" layer-type="overlay" :visible="true">
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
        key="all_dispatches"
        :name="$t('common.dispatch', 2)"
        layer-type="overlay"
        :visible="
            showAllDispatches ||
            livemap.activeLayers.length === 0 ||
            livemap.activeLayers.includes($t('common.dispatch', 2)) ||
            dispatchQuery.length > 0
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
        <div class="flex flex-col gap-2">
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
