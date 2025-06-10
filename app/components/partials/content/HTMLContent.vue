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
        class="tiptap prose prose-sm sm:prose-base lg:prose-lg dark:prose-invert hyphen-auto min-w-full max-w-full whitespace-pre-line break-words"
        :class="[
            'hover:prose-a:text-blue-500',
            'dark:hover:prose-a:text-blue-300',
            'prose-headings:my-0.5',
            'prose-lead:my-0.5',
            'prose-h1:my-0.5',
            'prose-h2:my-0.5',
            'prose-h3:my-0.5',
            'prose-h4:my-0.5',
            'prose-p:my-0.5',
            'prose-a:my-0.5',
            'prose-blockquote:my-0.5',
            'prose-figure:my-0.5',
            'prose-figcaption:my-0.5',
            'prose-strong:my-0.5',
            'prose-em:my-0.5',
            'prose-kbd:my-0.5',
            'prose-code:my-0.5',
            'prose-pre:my-0.5',
            'prose-ol:my-0.5',
            'prose-ul:my-0.5',
            'prose-li:my-0.5',
            'prose-table:my-0.5',
            'prose-thead:my-0.5',
            'prose-tr:my-0.5',
            'prose-th:my-0.5',
            'prose-td:my-0.5',
            'prose-img:my-0.5',
            'prose-video:my-0.5',
            'prose-hr:my-0.5',
        ]"
        :value="value"
        v-bind="$attrs"
    />
</template>
