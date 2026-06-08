<script lang="ts" setup>
import type { Editor } from '@tiptap/core';
import type { MapBlockAttrs } from '~/composables/tiptap/extensions/MapBlock';
import MapBlockEditorForm from './MapBlockEditorForm.vue';
import { tileLayers } from '~/types/livemap';

const props = defineProps<{
    editor: Editor;
    disabled?: boolean;
}>();

const settingsStore = useSettingsStore();
const { livemapTileLayer } = storeToRefs(settingsStore);

const open = ref(false);
const draft = ref<MapBlockAttrs>({
    x: 0,
    y: 0,
    zoom: 2,
    postal: '',
    layer: livemapTileLayer.value || tileLayers[0]!.key,
});

watch(open, (isOpen) => {
    if (isOpen) return;

    draft.value = {
        x: 0,
        y: 0,
        zoom: 2,
        postal: '',
        layer: livemapTileLayer.value || tileLayers[0]!.key,
    };
});

function insertMapBlock(): void {
    props.editor?.chain().focus().insertMapBlock(draft.value).run();
    open.value = false;
    draft.value = {
        x: 0,
        y: 0,
        zoom: 2,
        postal: '',
        layer: livemapTileLayer.value || tileLayers[0]!.key,
    };
}
</script>

<template>
    <UPopover v-model:open="open">
        <UTooltip text="Map">
            <UButton color="neutral" variant="ghost" icon="i-mdi-map-marker" :disabled="disabled" />
        </UTooltip>

        <template #content>
            <div class="flex w-[28rem] max-w-[calc(100vw-2rem)] flex-col gap-3 p-4">
                <h3 class="block font-medium">{{ $t('common.map') }}</h3>

                <MapBlockEditorForm v-model="draft" :disabled="disabled" />

                <UButton
                    block
                    color="neutral"
                    icon="i-mdi-map-marker-plus"
                    :label="$t('common.insert')"
                    :disabled="disabled"
                    @click="insertMapBlock"
                />
            </div>
        </template>
    </UPopover>
</template>
