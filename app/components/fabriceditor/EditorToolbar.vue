<script setup lang="ts">
import type { DropdownMenuItem } from '@nuxt/ui';
import { useFabricEditor } from '~/composables/useFabricEditor';

// Get methods from composable to manipulate canvas
const {
    exportJSON,
    exportSVG,
    zoom,
    resetZoom,
    fitDocumentToView,
    undo,
    redo,
    groupObject,
    ungroupObject,
    bringForward,
    sendBackward,
    bringToFront,
    sendToBack,
    pickingColor,
    pickedColor,
} = useFabricEditor();

const exportMenuItems = computed<DropdownMenuItem[]>(() => [
    {
        label: 'Export JSON',
        icon: 'i-mdi-code-json',
        onClick: () => exportJSON(),
    },
    {
        label: 'Export SVG',
        icon: 'i-mdi-svg',
        onClick: () => exportSVG(),
    },
]);
</script>

<template>
    <div
        class="mx-auto flex w-full max-w-(--breakpoint-xl) flex-1 snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto"
    >
        <UTooltip>
            <UButton
                variant="ghost"
                icon="i-mdi-select-color"
                :style="{ backgroundColor: pickedColor }"
                @click="pickingColor = !pickingColor"
            />
        </UTooltip>

        <UButtonGroup>
            <UTooltip :text="$t('common.undo')">
                <UButton variant="ghost" icon="i-mdi-undo" @click="undo" />
            </UTooltip>
            <UTooltip :text="$t('common.redo')">
                <UButton variant="ghost" icon="i-mdi-redo" @click="redo" />
            </UTooltip>
        </UButtonGroup>

        <UButtonGroup>
            <UButton variant="ghost" icon="i-mdi-group" @click="groupObject">Group</UButton>
            <UButton variant="ghost" icon="i-mdi-ungroup" @click="ungroupObject">Ungroup</UButton>
        </UButtonGroup>

        <UButtonGroup>
            <UButton variant="ghost" icon="i-mdi-arrange-bring-forward" @click="bringToFront" />
            <UButton variant="ghost" icon="i-mdi-arrange-bring-to-front" @click="bringForward" />
            <UButton variant="ghost" icon="i-mdi-arrange-send-to-back" @click="sendToBack" />
            <UButton variant="ghost" icon="i-mdi-arrange-send-backward" @click="sendBackward" />
        </UButtonGroup>

        <UDropdownMenu :items="exportMenuItems">
            <UButton variant="ghost" icon="i-mdi-export">Export</UButton>
        </UDropdownMenu>

        <UFormField>
            <div class="inline-flex gap-2">
                <UInputNumber
                    v-model="zoom"
                    class="w-30"
                    :step="0.1"
                    :min="0.1"
                    :max="3"
                    decrement-icon="i-mdi-zoom-out"
                    increment-icon="i-mdi-zoom-in"
                    :format-options="{
                        style: 'percent',
                    }"
                />

                <UButton icon="i-mdi-fit-to-screen" @click="fitDocumentToView">Zoom to Fit</UButton>
                <UButton variant="outline" icon="i-mdi-number-zero-circle-outline" @click="resetZoom">Reset</UButton>
            </div>
        </UFormField>
    </div>
</template>
