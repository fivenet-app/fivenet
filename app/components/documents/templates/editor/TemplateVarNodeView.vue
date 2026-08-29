<script setup lang="ts">
import { NodeViewWrapper, nodeViewProps } from '@tiptap/vue-3';
import TemplateVarForm from './TemplateVarForm.vue';

const props = defineProps(nodeViewProps);

const open = ref(false);

const displayValue = computed(() => {
    const opening = props.node.attrs['data-left-trim'] ? '{{-' : '{{';
    const closing = props.node.attrs['data-right-trim'] ? '-}}' : '}}';
    return `${opening} ${props.node.attrs['data-template-var']} ${closing}`;
});

function save(value: { value: string; leftTrim: boolean; rightTrim: boolean }): void {
    props.updateAttributes({
        'data-template-var': value.value,
        'data-left-trim': value.leftTrim,
        'data-right-trim': value.rightTrim,
    });
    open.value = false;
}
</script>

<template>
    <NodeViewWrapper as="span" class="template-var" contenteditable="false">
        <UPopover v-model:open="open" :content="{ side: 'top', sideOffset: 8 }">
            <button type="button" class="font-mono" contenteditable="false" @mousedown.prevent>
                {{ displayValue }}
            </button>

            <template #content>
                <div class="w-80 p-4">
                    <div class="mb-3 font-medium">
                        {{ $t('components.partials.tiptap_editor.extensions.template_var.title') }}
                    </div>
                    <TemplateVarForm
                        editing
                        :value="node.attrs['data-template-var']"
                        :left-trim="node.attrs['data-left-trim']"
                        :right-trim="node.attrs['data-right-trim']"
                        @submit="save"
                    />
                </div>
            </template>
        </UPopover>
    </NodeViewWrapper>
</template>
