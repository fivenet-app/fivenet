<script lang="ts" setup>
import { LControl, LLayerGroup } from '@vue-leaflet/vue-leaflet';
import DispatchMarker from '~/components/centrum/livemap/DispatchMarker.vue';
import { useCentrumStore } from '~/store/centrum';
import { useSettingsStore } from '~/store/settings';
import { Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';

defineProps<{
    showAllDispatches?: boolean;
}>();

defineEmits<{
    (e: 'selected', dsp: Dispatch): void;
}>();

const centrumStore = useCentrumStore();
const { dispatches, ownDispatches } = storeToRefs(centrumStore);
const settingsStore = useSettingsStore();
const { livemap } = storeToRefs(settingsStore);

const query = ref<string>('');
const dispatchesFiltered = computed(() =>
    dispatches.value.filter((m) => (m.user?.firstname + ' ' + m.user?.lastname).includes(query.value)),
);
</script>

<template>
    <LLayerGroup key="your_dispatches" :name="$t('common.your_dispatches')" layer-type="overlay" :visible="true">
        <DispatchMarker
            v-for="dispatch in ownDispatches"
            :dispatch="dispatch"
            @selected="$emit('selected', $event)"
            :size="livemap.markerSize"
        />
    </LLayerGroup>

    <LLayerGroup key="all_dispatches" :name="$t('common.dispatch', 2)" layer-type="overlay" :visible="showAllDispatches">
        <DispatchMarker
            v-for="dispatch in dispatchesFiltered.filter((d) => !ownDispatches.includes(d))"
            :key="dispatch.id.toString()"
            :dispatch="dispatch"
            @selected="$emit('selected', $event)"
            :size="livemap.markerSize"
        />
    </LLayerGroup>

    <LControl position="bottomleft">
        <div class="form-control flex flex-col gap-2">
            <div>
                <input
                    v-model="query"
                    class="w-full p-0.5 px-1 bg-clip-padding rounded-md border-2 border-black/20"
                    type="text"
                    name="searchPlayer"
                    :placeholder="`${$t('common.dispatch', 2)} ${$t('common.filter')}`"
                />
            </div>
        </div>
    </LControl>
</template>
