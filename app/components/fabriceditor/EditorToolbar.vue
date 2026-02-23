<script setup lang="ts">
import type { DropdownMenuItem } from '@nuxt/ui';
import { useFabricEditor } from '~/composables/useFabricEditor';

const { t } = useI18n();

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
        label: t('components.fabric_editor.export_json'),
        icon: 'i-mdi-code-json',
        onClick: () => exportJSON(),
    },
    {
        label: t('components.fabric_editor.export_svg'),
        icon: 'i-mdi-svg',
        onClick: () => exportSVG(),
    },
]);
</script>

<template>
    <div
        class="mx-auto flex w-full max-w-(--breakpoint-xl) flex-1 snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto"
    >
        <UTooltip :text="$t('components.fabric_editor.pick_color')">
            <UButton
                variant="ghost"
                icon="i-mdi-select-color"
                :style="{ backgroundColor: pickedColor }"
                @click="pickingColor = !pickingColor"
            />
        </UTooltip>

        <UFieldGroup>
            <UTooltip :text="$t('common.undo')">
                <UButton variant="ghost" icon="i-mdi-undo" @click="undo" />
            </UTooltip>
            <UTooltip :text="$t('common.redo')">
                <UButton variant="ghost" icon="i-mdi-redo" @click="redo" />
            </UTooltip>
        </UFieldGroup>

        <UFieldGroup>
            <UTooltip :text="$t('components.fabric_editor.group')">
                <UButton variant="ghost" icon="i-mdi-group" @click="groupObject" />
            </UTooltip>
            <UTooltip :text="$t('components.fabric_editor.ungroup')">
                <UButton variant="ghost" icon="i-mdi-ungroup" @click="ungroupObject" />
            </UTooltip>
        </UFieldGroup>

        <UFieldGroup>
            <UTooltip :text="$t('components.fabric_editor.bring_forward')">
                <UButton variant="ghost" icon="i-mdi-arrange-bring-forward" @click="bringToFront" />
            </UTooltip>
            <UTooltip :text="$t('components.fabric_editor.bring_to_front')">
                <UButton variant="ghost" icon="i-mdi-arrange-bring-to-front" @click="bringForward" />
            </UTooltip>
            <UTooltip :text="$t('components.fabric_editor.send_to_back')">
                <UButton variant="ghost" icon="i-mdi-arrange-send-to-back" @click="sendToBack" />
            </UTooltip>
            <UTooltip :text="$t('components.fabric_editor.send_backward')">
                <UButton variant="ghost" icon="i-mdi-arrange-send-backward" @click="sendBackward" />
            </UTooltip>
        </UFieldGroup>

        <UDropdownMenu :items="exportMenuItems">
            <UButton variant="ghost" icon="i-mdi-export" :label="$t('components.fabric_editor.export')" />
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

                <UFieldGroup>
                    <UButton
                        icon="i-mdi-fit-to-screen"
                        :label="$t('components.fabric_editor.zoom_to_fit')"
                        @click="fitDocumentToView"
                    />
                    <UButton
                        variant="outline"
                        icon="i-mdi-number-zero-circle-outline"
                        :label="$t('components.fabric_editor.reset_zoom')"
                        @click="resetZoom"
                    />
                </UFieldGroup>
            </div>
        </UFormField>
    </div>
</template>
