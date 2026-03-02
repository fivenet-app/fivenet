<script lang="ts" setup>
import EditorSidebar from './EditorSidebar.vue';

const props = withDefaults(
    defineProps<{
        maxHeight?: number;
        maxWidth?: number;
        backgroundColor?: string;
        disabled?: boolean;
    }>(),
    {
        maxHeight: 350,
        maxWidth: 900,
        backgroundColor: undefined,
        disabled: false,
    },
);

const svgData = defineModel<string | undefined>({ required: true });

const canvasContainer = useTemplateRef('canvasContainer');
const canvasEl = useTemplateRef('canvasEl');

// Get composable state and methods
const { canvas, documentSize, initCanvas, importSVG, fitToView } = useFabricEditor();

onMounted(async () => {
    if (!canvasContainer.value) {
        console.error('EditorCanvas requires a container element');
        return;
    }
    if (!canvasEl.value) {
        console.error('EditorCanvas requires a canvas element');
        return;
    }

    // Initialize Fabric canvas
    const containerElement = canvasContainer.value;
    // Set canvas dimensions to container size
    const width = containerElement?.clientWidth || 800;
    const height = containerElement?.clientHeight || 600;

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

    initCanvas(canvasContainer.value, canvasEl.value, { width, height });

    if (svgData.value && svgData.value !== '') importSVG(svgData.value);

    fitToView();

    const emitChange = () => (svgData.value = canvas.value?.toSVG());
    canvas.value?.on('object:added', emitChange);
    canvas.value?.on('object:modified', emitChange);
    canvas.value?.on('object:removed', emitChange);
});
</script>

<template>
    <!-- Container: full-screen flex layout with top toolbar and content area -->
    <div class="flex h-full max-w-screen flex-col">
        <!-- Main content: canvas and sidebar -->
        <div class="flex flex-1 flex-col overflow-hidden lg:flex-row">
            <!-- Canvas container fills remaining space -->
            <div ref="canvasContainer" class="min-w-0 flex-1 overflow-hidden">
                <!-- The canvas element for Fabric.js -->
                <canvas ref="canvasEl" class="h-full w-full" v-bind="$attrs" />
            </div>

            <!-- Sidebar on the right with fixed width -->
            <EditorSidebar class="w-full shrink-0 border-l border-l-default bg-default lg:max-w-sm lg:min-w-64">
                <template #sidebar-top>
                    <slot name="sidebar-top" />
                </template>
            </EditorSidebar>
        </div>
    </div>
</template>
