<script setup lang="ts">
import type { Editor } from '@tiptap/vue-3';
import TemplateBlockInsert from './TemplateBlockInsert.vue';
import TemplateVarInsert from './TemplateVarInsert.vue';

defineProps<{
    editor: Editor;
    disabled?: boolean;
}>();
</script>

<template>
    <TemplateVarInsert :editor="editor" :disabled="disabled" />
    <TemplateBlockInsert :editor="editor" :disabled="disabled" />
</template>
