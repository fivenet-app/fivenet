<script lang="ts" setup>
import { LControl, LLayerGroup } from '@vue-leaflet/vue-leaflet';
import { computedAsync } from '@vueuse/core';
import DispatchDetails from '~/components/centrum/dispatches/DispatchDetails.vue';
import DispatchMarker from '~/components/centrum/livemap/DispatchMarker.vue';
import { useCentrumStore } from '~/store/centrum';
import { useSettingsStore } from '~/store/settings';
import { Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';

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

const dispatchQuery = ref<string>('');
const dispatchesFiltered = computedAsync(async () =>
    [...dispatches.value.values()].filter(
        (m) =>
            !ownDispatches.value.includes(m.id) &&
            (m.creator?.firstname + ' ' + m.creator?.lastname).toLowerCase().includes(dispatchQuery.value.toLowerCase()),
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
            :key="dispatch.toString()"
            :dispatch="dispatches.get(dispatch)!"
            :size="livemap.markerSize"
            @selected="
                selectedDispatch = $event;
                open = true;
            "
        />
    </LLayerGroup>

    <LLayerGroup key="all_dispatches" :name="$t('common.dispatch', 2)" layer-type="overlay" :visible="showAllDispatches">
        <DispatchMarker
            v-for="dispatch in dispatchesFiltered"
            :key="dispatch.id.toString()"
            :dispatch="dispatch"
            :size="livemap.markerSize"
            @selected="
                selectedDispatch = $event;
                open = true;
            "
        />
    </LLayerGroup>

    <LControl position="bottomleft">
        <div class="form-control flex flex-col gap-2">
            <div>
                <input
                    v-model="dispatchQuery"
                    class="w-full max-w-[11rem] p-0.5 px-1 bg-clip-padding rounded-md border-2 border-black/20"
                    type="text"
                    name="searchPlayer"
                    :placeholder="`${$t('common.dispatch', 2)} ${$t('common.filter')}`"
                    @focusin="focusTablet(true)"
                    @focusout="focusTablet(false)"
                />
            </div>
        </div>
    </LControl>
</template>
