<script lang="ts" setup>
import EditorCanvas from './EditorCanvas.vue';
import EditorSidebar from './EditorSidebar.vue';

withDefaults(
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

const svgData = defineModel<string>({ required: true });
</script>

<template>
    <!-- Container: full-screen flex layout with top toolbar and content area -->
    <div class="flex h-full max-w-screen flex-col">
        <!-- Main content: canvas and sidebar -->
        <div class="flex min-h-0 flex-1 flex-col overflow-hidden lg:flex-row">
            <div ref="canvasContainer" class="min-h-0 min-w-0 flex-1 overflow-hidden">
                <!-- Canvas area fills remaining space -->
                <EditorCanvas
                    v-model="svgData"
                    :max-height="maxHeight"
                    :max-width="maxWidth"
                    :background-color="backgroundColor"
                    :disabled="disabled"
                    v-bind="$attrs"
                />
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
