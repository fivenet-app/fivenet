<script lang="ts" setup>
import { LControl, LLayerGroup } from '@vue-leaflet/vue-leaflet';
import { computedAsync } from '@vueuse/core';
import DispatchDetails from '~/components/centrum/dispatches/DispatchDetails.vue';
import DispatchMarker from '~/components/centrum/livemap/DispatchMarker.vue';
import { useCentrumStore } from '~/store/centrum';
import { useSettingsStore } from '~/store/settings';
import { Dispatch } from '~~/gen/ts/resources/centrum/dispatches';

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

const dispatchQueryRaw = ref<string>('');
const dispatchQuery = computed(() => dispatchQueryRaw.value.trim().toLowerCase());

const dispatchesFiltered = computedAsync(async () =>
    [...dispatches.value.values()].filter(
        (m) =>
            !ownDispatches.value.includes(m.id) &&
            (m.id.startsWith(dispatchQuery.value) ||
                m.message.toLowerCase().includes(dispatchQuery.value) ||
                (m.creator?.firstname + ' ' + m.creator?.lastname).toLowerCase().includes(dispatchQuery.value)),
    ),
);

const selectedDispatch = ref<Dispatch | undefined>();
const open = ref(false);
</script>

<template>
    <template v-if="selectedDispatch">
        <DispatchDetails :dispatch="selectedDispatch" :open="open" @close="open = false" @goto="$emit('goto', $event)" />
    </template>

    <LLayerGroup key="your_dispatches" :name="$t('common.your_dispatches')" layer-type="overlay" :visible="true">
        <DispatchMarker
            v-for="dispatch in ownDispatches"
            :key="dispatch"
            :dispatch="dispatches.get(dispatch)!"
            :size="livemap.markerSize"
            @selected="
                selectedDispatch = $event;
                open = true;
            "
            @goto="$emit('goto', $event)"
        />
    </LLayerGroup>

    <LLayerGroup
        key="all_dispatches"
        :name="$t('common.dispatch', 2)"
        layer-type="overlay"
        :visible="showAllDispatches || dispatchQueryRaw.length > 0"
    >
        <DispatchMarker
            v-for="dispatch in dispatchesFiltered"
            :key="dispatch.id"
            :dispatch="dispatch"
            :size="livemap.markerSize"
            @selected="
                selectedDispatch = $event;
                open = true;
            "
            @goto="$emit('goto', $event)"
        />
    </LLayerGroup>

    <LControl position="bottomleft">
        <div class="flex flex-col gap-2">
            <UInput
                v-model="dispatchQueryRaw"
                type="text"
                name="searchPlayer"
                :placeholder="`${$t('common.dispatch', 2)} ${$t('common.filter')}`"
                @focusin="focusTablet(true)"
                @focusout="focusTablet(false)"
            />
        </div>
    </LControl>
</template>
