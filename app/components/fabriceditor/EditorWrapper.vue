<script lang="ts" setup>
import EditorCanvas from './EditorCanvas.vue';
import EditorSidebar from './EditorSidebar.vue';

withDefaults(
    defineProps<{
        maxHeight?: number;
        maxWidth?: number;
        backgroundColor?: string;
    }>(),
    {
        maxHeight: 350,
        maxWidth: 900,
        backgroundColor: undefined,
    },
);
</script>

<template>
    <!-- Container: full-screen flex layout with top toolbar and content area -->
    <div class="flex h-full max-w-screen flex-col">
        <!-- Main content: canvas and sidebar -->
        <div class="flex flex-1 overflow-hidden">
            <div ref="canvasContainer" class="flex-1 overflow-hidden">
                <!-- Canvas area fills remaining space -->
                <EditorCanvas :min-height="350" :min-width="900" :background-color="backgroundColor" v-bind="$attrs" />
            </div>

            <!-- Sidebar on the right with fixed width -->
            <EditorSidebar class="w-full max-w-sm min-w-64 shrink-0 border-l border-l-default bg-default" />
        </div>
    </div>
</template>
