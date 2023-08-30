<script lang="ts" setup>
import { LControl, LLayerGroup } from '@vue-leaflet/vue-leaflet';
import { default as DispatchDetails } from '~/components/centrum/dispatches/Details.vue';
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

const query = ref<string>('');
const dispatchesFiltered = computed(() =>
    dispatches.value.filter((m) => (m.user?.firstname + ' ' + m.user?.lastname).includes(query.value)),
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
            :dispatch="dispatch"
            @selected="
                selectedDispatch = $event;
                open = true;
            "
            :size="livemap.markerSize"
        />
    </LLayerGroup>

    <LLayerGroup key="all_dispatches" :name="$t('common.dispatch', 2)" layer-type="overlay" :visible="showAllDispatches">
        <DispatchMarker
            v-for="dispatch in dispatchesFiltered.filter((d) => !ownDispatches.includes(d))"
            :key="dispatch.id.toString()"
            :dispatch="dispatch"
            @selected="
                selectedDispatch = $event;
                open = true;
            "
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
