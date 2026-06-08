<script lang="ts" setup>
import { NodeViewWrapper, type NodeViewProps } from '@tiptap/vue-3';
import { createMapBlockAttrs, normalizeMapBlockAttrs, type MapBlockAttrs } from '~/composables/tiptap/extensions/MapBlock';
import MapBlockContentView from '~/components/partials/content/MapBlockContentView.vue';
import MapBlockEditorForm from './MapBlockEditorForm.vue';
import MapBlockEditorPopoverShell from './MapBlockEditorPopoverShell.vue';

const props = defineProps<NodeViewProps>();

const open = ref(false);
const draft = ref<MapBlockAttrs>(createMapBlockAttrs(props.node.attrs));
const contentProps = computed(() => normalizeMapBlockAttrs(props.node.attrs));

watch(
    () => props.node.attrs,
    (attrs) => {
        if (open.value) return;

        draft.value = createMapBlockAttrs(attrs);
    },
    { deep: true, immediate: true },
);

watch(open, (isOpen) => {
    if (isOpen) return;

    draft.value = createMapBlockAttrs(props.node.attrs);
});

function save(): void {
    props.updateAttributes(createMapBlockAttrs(draft.value));
    open.value = false;
}
</script>

<template>
    <NodeViewWrapper class="inline-flex align-middle" as="span">
        <span class="group relative inline-flex align-middle">
            <MapBlockContentView :show-goto-coords="false" v-bind="contentProps" />

            <div
                v-if="editor.isEditable"
                class="absolute top-2 right-2 flex gap-1 opacity-100 transition-opacity group-hover:opacity-100"
            >
                <MapBlockEditorPopoverShell v-model:open="open" :title="$t('common.map')">
                    <template #trigger>
                        <UTooltip :text="$t('common.edit')">
                            <UButton color="neutral" variant="soft" size="xs" icon="i-mdi-pencil" :label="$t('common.edit')" />
                        </UTooltip>
                    </template>

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
                </MapBlockEditorPopoverShell>

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
