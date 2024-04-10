<script lang="ts" setup>
import { useLivemapStore } from '~/store/livemap';

const livemapStore = useLivemapStore();
const { goto } = livemapStore;
const { offsetLocationZoom } = storeToRefs(livemapStore);

onMounted(() => {
    offsetLocationZoom.value = true;
});

onBeforeUnmount(() => {
    offsetLocationZoom.value = false;
});
</script>

<template>
    <ClientOnly>
        <LazyCentrumLivemapCentrumSidebar @goto="goto($event)" />
    </ClientOnly>
</template>
