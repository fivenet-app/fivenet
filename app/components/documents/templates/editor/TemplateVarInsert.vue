<script setup lang="ts">
import type { Editor } from '@tiptap/core';
import TemplateVarForm from './TemplateVarForm.vue';

const props = defineProps<{
    editor: Editor;
    disabled?: boolean;
}>();

function insert(value: { value: string; leftTrim: boolean; rightTrim: boolean }): void {
    props.editor?.commands.insertTemplateVar({
        value: value.value,
        leftTrim: value.leftTrim,
        rightTrim: value.rightTrim,
    });
}
</script>

<template>
    <UPopover>
        <UTooltip :text="$t('components.partials.tiptap_editor.extensions.template_var.title')">
            <UButton color="neutral" variant="ghost" icon="i-mdi-variable" :disabled="disabled" />
        </UTooltip>

        <template #content>
            <div class="flex flex-col gap-2 p-4">
                <h3 class="block font-medium">
                    {{ $t('components.partials.tiptap_editor.extensions.template_var.title') }}
                </h3>

                <TemplateVarForm @submit="insert" />
            </div>
        </template>
    </UPopover>
</template>
