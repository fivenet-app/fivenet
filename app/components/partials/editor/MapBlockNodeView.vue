<script lang="ts" setup>
import { NodeViewWrapper, type NodeViewProps } from '@tiptap/vue-3';
import type { MapBlockAttrs } from '~/composables/tiptap/extensions/MapBlock';
import MapBlockContentView from '~/components/partials/content/MapBlockContentView.vue';
import MapBlockEditorForm from './MapBlockEditorForm.vue';
import { tileLayers } from '~/types/livemap';

const props = defineProps<NodeViewProps>();

const open = ref(false);
const draft = ref<MapBlockAttrs>({
    x: Number(props.node.attrs.x ?? 0),
    y: Number(props.node.attrs.y ?? 0),
    zoom: Number(props.node.attrs.zoom ?? 2),
    postal: String(props.node.attrs.postal ?? ''),
    layer: String(props.node.attrs.layer ?? ''),
});

watch(
    () => props.node.attrs,
    (attrs) => {
        if (open.value) return;

        draft.value = {
            x: Number(attrs.x ?? 0),
            y: Number(attrs.y ?? 0),
            zoom: Number(attrs.zoom ?? 2),
            postal: String(attrs.postal ?? ''),
            layer: String(attrs.layer ?? ''),
        };
    },
    { deep: true, immediate: true },
);

watch(open, (isOpen) => {
    if (isOpen) return;

    draft.value = {
        x: Number(props.node.attrs.x ?? 0),
        y: Number(props.node.attrs.y ?? 0),
        zoom: Number(props.node.attrs.zoom ?? 2),
        postal: String(props.node.attrs.postal ?? ''),
        layer: String(props.node.attrs.layer ?? ''),
    };
});

function save(): void {
    props.updateAttributes({
        x: draft.value.x,
        y: draft.value.y,
        zoom: draft.value.zoom,
        postal: draft.value.postal ?? '',
        layer: draft.value.layer || tileLayers[0]!.key,
    });
    open.value = false;
}
</script>

<template>
    <NodeViewWrapper class="inline-flex align-middle" as="span">
        <span class="group relative inline-flex align-middle">
            <MapBlockContentView
                :x="Number(props.node.attrs.x ?? 0)"
                :y="Number(props.node.attrs.y ?? 0)"
                :zoom="Number(props.node.attrs.zoom ?? 2)"
                :postal="String(props.node.attrs.postal ?? '') || undefined"
                :layer="String(props.node.attrs.layer ?? '') || undefined"
                :show-goto-coords="false"
            />

            <div
                v-if="editor.isEditable"
                class="absolute top-2 right-2 flex gap-1 opacity-100 transition-opacity group-hover:opacity-100"
            >
                <UPopover v-model:open="open">
                    <UTooltip :text="$t('common.edit')">
                        <UButton color="neutral" variant="soft" size="xs" icon="i-mdi-pencil" :label="$t('common.edit')" />
                    </UTooltip>

                    <template #content>
                        <div class="flex w-[28rem] max-w-[calc(100vw-2rem)] flex-col gap-3 p-4">
                            <h3 class="block font-medium">{{ $t('common.map') }}</h3>

                            <MapBlockEditorForm v-model="draft" :disabled="!editor.isEditable" />

                            <UFieldGroup>
                                <UButton
                                    block
                                    color="neutral"
                                    icon="i-mdi-close"
                                    :label="$t('common.cancel')"
                                    @click="open = false"
                                />

                                <UButton
                                    block
                                    color="neutral"
                                    icon="i-mdi-content-save"
                                    :label="$t('common.save')"
                                    @click="save"
                                />
                            </UFieldGroup>
                        </div>
                    </template>
                </UPopover>

                <UButton
                    color="error"
                    variant="soft"
                    size="xs"
                    icon="i-mdi-delete"
                    :label="$t('common.delete')"
                    @click="props.deleteNode()"
                />
            </div>
        </span>
    </NodeViewWrapper>
</template>
