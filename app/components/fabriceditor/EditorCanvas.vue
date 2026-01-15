<script setup lang="ts">
import type { FabricObject } from 'fabric';
import { useFabricEditor } from '~/composables/useFabricEditor';

const props = defineProps<{
    maxHeight?: number;
    maxWidth?: number;
    backgroundColor?: string;
    disabled?: boolean;
}>();

const svgData = defineModel<string | undefined>({ required: true });

const canvasEl = ref<HTMLCanvasElement | null>(null);

const fabric = await import('fabric');

// Get composable state and methods
const { canvas, documentSize, initCanvas } = useFabricEditor();

watch(svgData, async (value) => {
    if (!canvas || !value) return;

    await fabric.loadSVGFromString(value, () => canvas.value?.renderAll());
});

onMounted(async () => {
    // Initialize Fabric canvas
    const canvasElement = canvasEl.value!;
    // Set canvas dimensions to container size
    const parent = canvasElement.parentElement;
    const width = props.maxWidth || parent?.clientWidth || 800;
    const height = props.maxHeight || parent?.clientHeight || 600;

    // If max dimensions are provided, set them in document size and disable resize
    if (props.maxHeight) {
        documentSize.value.height = props.maxHeight;
        documentSize.value.disabled = true;
    }
    if (props.maxWidth) {
        documentSize.value.width = props.maxWidth;
        documentSize.value.disabled = true;
    }
    if (props.backgroundColor) {
        documentSize.value.fill = props.backgroundColor;
    }

    initCanvas(canvasElement, {
        width,
        height,
    });

    if (svgData.value) {
        const loadedSvg = await fabric.loadSVGFromString(svgData.value);

        if (loadedSvg.objects) {
            if (loadedSvg.options['width'] && loadedSvg.options['height']) {
                documentSize.value.width = parseInt(loadedSvg.options['width']);
                documentSize.value.height = parseInt(loadedSvg.options['height']);
            }

            canvas.value?.add(...loadedSvg.objects.filter((obj): obj is FabricObject => !!obj));
        }
    }

    const emitChange = () => (svgData.value = canvas.value?.toSVG());
    canvas.value?.on('object:added', emitChange);
    canvas.value?.on('object:modified', emitChange);
    canvas.value?.on('object:removed', emitChange);
});
</script>

<template>
    <!-- The canvas element for Fabric.js -->
    <canvas ref="canvasEl" class="h-full w-full"></canvas>
</template>
