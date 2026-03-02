<script setup lang="ts">
import type { DropdownMenuItem } from '@nuxt/ui';
import { useFabricEditor } from '~/composables/useFabricEditor';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import ImportModal from './ImportModal.vue';

const { t } = useI18n();

const overlay = useOverlay();

// Get methods from composable to manipulate canvas
const {
    history,
    redoStack,

    importJSON,
    importSVG,
    exportJSON,
    exportSVG,
    zoom,
    fitToView,
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

const importModal = overlay.create(ImportModal);

const notifications = useNotificationsStore();

const fileMenuItems = computed<DropdownMenuItem[][]>(() => [
    [
        {
            label: t('components.fabric_editor.import.title'),
            icon: 'i-mdi-import',
            type: 'label',
        },
        {
            label: t('components.fabric_editor.import.json'),
            icon: 'i-mdi-code-json',
            onClick: () =>
                importModal.open({
                    modelValue: '',
                    'onUpdate:modelValue': async (value) => {
                        try {
                            await importJSON(value);
                        } catch (e) {
                            notifications.add({
                                type: NotificationType.ERROR,
                                title: {
                                    key: 'components.fabric_editor.import.json_failed.title',
                                },
                                description: {
                                    key: 'components.fabric_editor.import.json_failed.content',
                                    parameters: {
                                        msg: (e as Error).message,
                                    },
                                },
                            });
                            console.error('Invalid JSON', e);
                        }
                    },
                }),
        },
        {
            label: t('components.fabric_editor.import.svg'),
            icon: 'i-mdi-svg',
            onClick: () =>
                importModal.open({
                    modelValue: '',
                    'onUpdate:modelValue': async (value) => {
                        try {
                            await importSVG(value);
                        } catch (e) {
                            notifications.add({
                                type: NotificationType.ERROR,
                                title: {
                                    key: 'components.fabric_editor.import.svg_failed.title',
                                },
                                description: {
                                    key: 'components.fabric_editor.import.svg_failed.content',
                                    parameters: {
                                        msg: (e as Error).message,
                                    },
                                },
                            });
                            console.error('Invalid SVG', e);
                        }
                    },
                }),
        },
        {
            label: t('components.fabric_editor.export.title'),
            icon: 'i-mdi-export',
            type: 'label',
        },
        {
            label: t('components.fabric_editor.export.json'),
            icon: 'i-mdi-code-json',
            onClick: () => exportJSON(),
        },
        {
            label: t('components.fabric_editor.export.svg'),
            icon: 'i-mdi-svg',
            onClick: () => exportSVG(),
        },
    ],
]);
</script>

<template>
    <div
        class="mx-auto flex w-full max-w-(--breakpoint-xl) flex-1 snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto"
    >
        <UFieldGroup>
            <UTooltip :text="$t('common.undo')">
                <UButton variant="ghost" icon="i-mdi-undo" :disabled="history.length === 0" @click="undo" />
            </UTooltip>

            <UTooltip :text="$t('common.redo')">
                <UButton variant="ghost" icon="i-mdi-redo" :disabled="redoStack.length === 0" @click="redo" />
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

        <UTooltip :text="$t('components.fabric_editor.pick_color')">
            <UButton
                variant="ghost"
                icon="i-mdi-select-color"
                :style="{ backgroundColor: pickedColor }"
                @click="pickingColor = !pickingColor"
            />
        </UTooltip>

        <UFormField>
            <div class="inline-flex gap-2">
                <UFieldGroup>
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

                    <UTooltip :text="$t('components.fabric_editor.zoom_to_fit')">
                        <UButton icon="i-mdi-fit-to-screen" @click="fitToView" />
                    </UTooltip>
                </UFieldGroup>
            </div>
        </UFormField>

        <UDropdownMenu :items="fileMenuItems">
            <UButton variant="ghost" trailing-icon="i-mdi-menu-down" :label="$t('components.fabric_editor.file_menu')" />
        </UDropdownMenu>
    </div>
</template>
