<script lang="ts" setup>
import { setWaypoint } from '~/composables/nui';
import { useLivemapStore } from '~/store/livemap';

const livemapStore = useLivemapStore();
const { location, offsetLocationZoom } = storeToRefs(livemapStore);

function goto(e: Coordinate) {
    location.value = { x: e.x, y: e.y };

    // Set in-game waypoint via NUI
    setWaypoint(e.x, e.y);
}

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
