<script lang="ts" setup>
import { useLivemapStore } from '~/stores/livemap';
import MapZoomControls from './MapZoomControls.vue';

const livemapStore = useLivemapStore();
const { zoom } = storeToRefs(livemapStore);
</script>

<template>
    <MapZoomControls v-model="zoom" />
</template>
