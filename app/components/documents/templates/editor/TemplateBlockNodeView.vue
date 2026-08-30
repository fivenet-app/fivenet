<script setup lang="ts">
import { NodeViewWrapper, nodeViewProps } from '@tiptap/vue-3';
import TemplateBlockEndForm from '~/components/documents/templates/editor/TemplateBlockEndForm.vue';
import TemplateBlockInsertForm from '~/components/documents/templates/editor/TemplateBlockInsertForm.vue';
import { isTemplateBlockActionValid } from '~/composables/tiptap/extensions/TemplateBlockValidation';

const props = defineProps(nodeViewProps);

const open = ref(false);
const editable = computed(() => props.editor?.isEditable ?? false);

const isEnd = computed(() => props.node.type.name === 'templateBlockEnd');
const isValid = computed(() => isEnd.value || isTemplateBlockActionValid(props.node.attrs['data-template-block']));
const displayValue = computed(() => {
    const opening = props.node.attrs['data-left-trim'] ? '{{-' : '{{';
    const closing = props.node.attrs['data-right-trim'] ? '-}}' : '}}';
    return `${opening} ${isEnd.value ? 'end' : props.node.attrs['data-template-block']} ${closing}`;
});

function save(value: { value: string; leftTrim: boolean; rightTrim: boolean }): void {
    props.updateAttributes({
        ...(isEnd.value ? {} : { 'data-template-block': value.value }),
        'data-left-trim': value.leftTrim,
        'data-right-trim': value.rightTrim,
    });
    open.value = false;
}

function deleteNode(): void {
    open.value = false;
    props.deleteNode();
}
</script>

<template>
    <NodeViewWrapper
        as="span"
        class="template-block inline-block rounded border border-dashed border-amber-400 !bg-amber-950 px-0.5 font-mono !text-amber-100"
        :class="{
            'cursor-pointer': editable,
            'template-block-invalid border-red-400 !bg-red-950 !text-red-100 motion-safe:animate-pulse': !isValid,
        }"
        contenteditable="false"
    >
        <template v-if="!editable">
            <span>{{ displayValue }}</span>
        </template>
        <UPopover v-else v-model:open="open" :content="{ side: 'top', sideOffset: 8 }">
            <button type="button" class="font-mono" contenteditable="false" @mousedown.prevent>
                {{ displayValue }}
            </button>

            <template #content>
                <div class="w-full max-w-86 p-4">
                    <div class="mb-3 font-medium">
                        {{ $t('components.partials.tiptap_editor.extensions.template_block.title') }}
                    </div>

                    <UAlert class="mb-3" icon="i-mdi-information-outline" variant="subtle">
                        <template #description>
                            <I18nT keypath="components.partials.tiptap_editor.extensions.template_block.close_info">
                                <template #end>
                                    <code class="font-mono" v-text="'{{ end }}'" />
                                </template>
                            </I18nT>
                        </template>
                    </UAlert>

                    <UAlert
                        v-if="!isValid"
                        class="mb-3"
                        color="error"
                        icon="i-mdi-alert-outline"
                        :description="$t('components.partials.tiptap_editor.extensions.template_block.invalid')"
                        variant="subtle"
                    />

                    <TemplateBlockInsertForm
                        v-if="!isEnd"
                        editing
                        :value="node.attrs['data-template-block']"
                        :left-trim="node.attrs['data-left-trim']"
                        :right-trim="node.attrs['data-right-trim']"
                        @delete="deleteNode"
                        @submit="save"
                    />
                    <TemplateBlockEndForm
                        v-else
                        :left-trim="node.attrs['data-left-trim']"
                        :right-trim="node.attrs['data-right-trim']"
                        @delete="deleteNode"
                        @submit="save"
                    />
                </div>
            </template>
        </UPopover>
    </NodeViewWrapper>
</template>
