<script lang="ts" setup>
import { default as CentrumSidebar } from '~/components/centrum/livemap/Sidebar.vue';
import { setWaypoint } from '~/composables/nui';
import { useLivemapStore } from '~/store/livemap';

const livemapStore = useLivemapStore();
const { location } = storeToRefs(livemapStore);

function goto(e: Coordinate) {
    location.value = { x: e.x, y: e.y };

    // Set in-game waypoint via NUI
    setWaypoint(e.x, e.y);
}
</script>

<template>
    <CentrumSidebar @goto="goto($event)" />
</template>
