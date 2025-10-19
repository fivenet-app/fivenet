<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">
import SvgCanvas from './SvgCanvas.client.vue';
import SvgLayers from './SvgLayers.vue';
import SvgProperties from './SvgProperties.vue';
import SvgShapes from './SvgShapes.vue';
import SvgToolbar from './SvgToolbar.vue';

const props = withDefaults(
    defineProps<{
        config?: Record<string, any>;
        initialSvg?: string;
    }>(),
    {
        config: () => ({}),
        initialSvg: '',
    },
);

const api = ref<any | null>(null);
const zoom = ref(1);

function onReady(a: any) {
    api.value = a;
    zoom.value = a.zoom ?? 1;
    a.bind('zoomed', (z: number) => {
        console.log('Zoom changed', z);
        zoom.value = z;
    });
    console.log('SVG Canvas ready', a);
}

watch(zoom, (v) => {
    console.log('Setting zoom to', v);
    api.value?.setZoom(v);
});
</script>

<template>
    <div class="flex h-[80vh] flex-col">
        <SvgToolbar v-model:zoom="zoom" :api="api" />

        <div class="flex grow gap-3 overflow-hidden p-3">
            <div class="w-64 shrink-0 space-y-3">
                <SvgShapes :api="api" />
                <SvgLayers :api="api" />
            </div>

            <div class="min-w-0 grow overflow-hidden border border-default">
                <SvgCanvas :config="props.config" :initial-svg="props.initialSvg" @ready="onReady" />
            </div>

            <div class="w-64 shrink-0">
                <SvgProperties :api="api" />
            </div>
        </div>
    </div>
</template>
