<script setup lang="ts">
import type { Editor } from '@tiptap/core';
import TemplateVarForm from './TemplateVarForm.vue';

const props = defineProps<{
    editor: Editor;
    disabled?: boolean;
}>();

const open = ref(false);

function insert(value: { value: string; leftTrim: boolean; rightTrim: boolean }): void {
    props.editor?.commands.insertTemplateVar({
        value: value.value,
        leftTrim: value.leftTrim,
        rightTrim: value.rightTrim,
    });

    open.value = false;
}
</script>

<template>
    <UPopover v-model:open="open">
        <UTooltip :text="$t('components.partials.tiptap_editor.extensions.template_var.title')">
            <UButton color="neutral" variant="ghost" icon="i-mdi-variable" :disabled="disabled" />
        </UTooltip>

        <template #content>
            <div class="flex w-full max-w-86 flex-col gap-2 p-4">
                <h3 class="block font-medium">
                    {{ $t('components.partials.tiptap_editor.extensions.template_var.title') }}
                </h3>

                <TemplateVarForm @submit="insert" />
            </div>
        </template>
    </UPopover>
</template>
