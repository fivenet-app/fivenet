<script lang="ts" setup>
import type { Editor } from '@tiptap/core';
import { createMapBlockAttrs, defaultMapBlockLayerKey, type MapBlockAttrs } from '~/composables/tiptap/extensions/MapBlock';
import MapBlockEditorForm from './MapBlockEditorForm.vue';
import MapBlockEditorPopoverShell from './MapBlockEditorPopoverShell.vue';

const props = defineProps<{
    editor: Editor;
    disabled?: boolean;
}>();

const settingsStore = useSettingsStore();
const { livemapTileLayer } = storeToRefs(settingsStore);

const open = ref(false);
function createDraft(): MapBlockAttrs {
    return createMapBlockAttrs({
        x: 0,
        y: 0,
        zoom: 2,
        postal: '',
        layer: livemapTileLayer.value || defaultMapBlockLayerKey,
    });
}

const draft = ref<MapBlockAttrs>(createDraft());

function resetDraft(): void {
    draft.value = createDraft();
}

watch(open, (isOpen) => {
    if (isOpen) return;

    resetDraft();
});

function insertMapBlock(): void {
    props.editor?.chain().focus().insertMapBlock(draft.value).run();
    open.value = false;
    resetDraft();
}
</script>

<template>
    <MapBlockEditorPopoverShell v-model:open="open" :title="$t('common.map')">
        <template #trigger>
            <UTooltip text="Map">
                <UButton color="neutral" variant="ghost" icon="i-mdi-map-marker" :disabled="disabled" />
            </UTooltip>
        </template>

        <MapBlockEditorForm v-model="draft" :disabled="disabled" />

        <UButton
            block
            color="neutral"
            icon="i-mdi-map-marker-plus"
            :label="$t('common.insert')"
            :disabled="disabled"
            @click="insertMapBlock"
        />
    </MapBlockEditorPopoverShell>
</template>
