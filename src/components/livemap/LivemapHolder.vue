<script lang="ts" setup>
import DispatchesLayer from '~/components/centrum/livemap/DispatchesLayer.vue';
import { default as CentrumSidebar } from '~/components/centrum/livemap/Sidebar.vue';
import { setWaypoint } from '~/components/centrum/nui';
import { useLivemapStore } from '~/store/livemap';
import Livemap from './Livemap.vue';

const livemapStore = useLivemapStore();
const { location } = storeToRefs(livemapStore);

function goto(e: { x: number; y: number }) {
    location.value = { x: e.x, y: e.y };

    // Set in-game waypoint via NUI
    setWaypoint(e.x, e.y);
}
</script>

<template>
    <Livemap>
        <template v-slot:default>
            <DispatchesLayer />
        </template>
        <template v-slot:afterMap>
            <div v-if="can('CentrumService.Stream')" class="lg:inset-y-0 lg:flex lg:w-50 lg:flex-col">
                <CentrumSidebar @goto="goto($event)" />
            </div>
        </template>
    </Livemap>
</template>
