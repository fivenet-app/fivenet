<script lang="ts" setup>
import type { Extensions } from '@tiptap/core';
import { NodeType, type Content } from '~~/gen/ts/resources/common/content/content';
import HTMLContent from './HTMLContent.vue';
import TiptapContent from './TiptapContent.vue';

withDefaults(
    defineProps<{
        value: Content;
        extensions?: Extensions;
    }>(),
    {
        extensions: () => [],
    },
);
</script>

<template>
    <TiptapContent v-if="value?.tiptapJson" :value="value.tiptapJson" :extensions="extensions" v-bind="$attrs" />
    <HTMLContent v-else :value="value.content ?? { attrs: {}, content: [], tag: '', type: NodeType.DOC }" v-bind="$attrs" />
</template>
