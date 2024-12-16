<script lang="ts" setup>
import type { JSONNode } from '~~/gen/ts/resources/common/content/content';
import HTMLContentRenderer from './HTMLContentRenderer.vue';

defineProps<{
    value: JSONNode;
}>();

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <HTMLContentRenderer
        v-bind="$attrs"
        :value="value"
        class="tiptap ProseMirror prose prose-sm sm:prose-base lg:prose-lg xl:prose-2xl dark:prose-invert min-w-full max-w-full break-words"
        :class="[
            'prose-headings:mt-0.5',
            'prose-lead:mt-0.5',
            'prose-h1:mt-0.5',
            'prose-h2:mt-0.5',
            'prose-h3:mt-0.5',
            'prose-h4:mt-0.5',
            'prose-p:mt-0.5',
            'prose-a:mt-0.5',
            'prose-blockquote:mt-0.5',
            'prose-figure:mt-0.5',
            'prose-figcaption:mt-0.5',
            'prose-strong:mt-0.5',
            'prose-em:mt-0.5',
            'prose-kbd:mt-0.5',
            'prose-code:mt-0.5',
            'prose-pre:mt-0.5',
            'prose-ol:mt-0.5',
            'prose-ul:mt-0.5',
            'prose-li:mt-0.5',
            'prose-table:mt-0.5',
            'prose-thead:mt-0.5',
            'prose-tr:mt-0.5',
            'prose-th:mt-0.5',
            'prose-td:mt-0.5',
            'prose-img:mt-0.5',
            'prose-video:mt-0.5',
            'prose-hr:mt-0.5',
        ]"
    />
</template>
