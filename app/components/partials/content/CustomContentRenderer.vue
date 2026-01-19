<script lang="ts" setup>
import type { Extensions } from '@tiptap/core';
import { NodeType, type Content } from '~~/gen/ts/resources/common/content/content';
import HTMLContent from './HTMLContent.vue';
import TiptapContent from './TiptapContent.vue';

withDefaults(
    defineProps<{
        value?: Content;
        extensions?: Extensions;
        placeholder?: string;
    }>(),
    {
        value: undefined,
        extensions: () => [],
        placeholder: undefined,
    },
);
</script>

<template>
    <TiptapContent
        v-if="value?.tiptapJson && !isEmptyDoc(value?.tiptapJson)"
        :value="value?.tiptapJson"
        :extensions="extensions"
        v-bind="$attrs"
    />
    <HTMLContent
        v-else-if="value?.content && !isEmptyRichContentDoc(value?.content)"
        :value="value?.content ?? { attrs: {}, content: [], tag: '', type: NodeType.DOC }"
        v-bind="$attrs"
    />
    <span v-else-if="placeholder" class="text-neutral-500 italic dark:text-neutral-400">{{ placeholder }}</span>
</template>
