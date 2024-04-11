<script lang="ts" setup>
import { useLivemapStore } from '~/store/livemap';

const livemapStore = useLivemapStore();
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
        <LazyCentrumLivemapCentrumSidebar />
    </ClientOnly>
</template>
