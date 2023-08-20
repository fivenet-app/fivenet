<script lang="ts" setup>
import { LLayerGroup } from '@vue-leaflet/vue-leaflet';
import DispatchMarker from '~/components/centrum/livemap/DispatchMarker.vue';
import { useCentrumStore } from '~/store/centrum';
import { Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';

defineEmits<{
    (e: 'select', dsp: Dispatch): void;
}>();

const centrumStore = useCentrumStore();
const { dispatches } = storeToRefs(centrumStore);
</script>

<template>
    <LLayerGroup :name="$t('common.dispatch', 2)" layer-type="overlay" :visible="true">
        <DispatchMarker v-for="dispatch in dispatches" :dispatch="dispatch" @select="$emit('select', $event)" />
    </LLayerGroup>
</template>
