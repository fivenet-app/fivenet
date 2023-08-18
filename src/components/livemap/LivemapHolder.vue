<script lang="ts" setup>
import CentrumSidebar from './CentrumSidebar.vue';
import Livemap from './Livemap.vue';

const livemapComponent = ref<InstanceType<typeof Livemap>>();

function goto(e: { x: number; y: number }) {
    if (livemapComponent.value) {
        livemapComponent.value.location = { x: e.x, y: e.y };
    }
}
</script>

<template>
    <Livemap ref="livemapComponent" :enable-centrum="true">
        <template v-slot:afterMap>
            <div v-if="can('CentrumService.Stream')" class="lg:inset-y-0 lg:flex lg:w-50 lg:flex-col">
                <CentrumSidebar @goto="goto($event)" />
            </div>
        </template>
    </Livemap>
</template>
