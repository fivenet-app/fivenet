<script lang="ts" setup>
import type { Extensions, JSONContent } from '@tiptap/core';
import { h, type ComputedRef, type Ref } from 'vue';
import PenaltyCalculatorContentView from '~/components/quickbuttons/penaltycalculator/PenaltyCalculatorContentView.vue';
import { Struct } from '~~/gen/ts/google/protobuf/struct';
import type { DocumentData } from '~~/gen/ts/resources/documents/data/data';
import TiptapContentRenderer from './TiptapContentRenderer.vue';

const props = withDefaults(
    defineProps<{
        value: Struct | undefined;
        extensions?: Extensions;
    }>(),
    {
        extensions: () => [],
    },
);

const builtInExtensions = useTiptapEditor();

const localDocumentData = ref<DocumentData | undefined>();
const documentData = inject<ComputedRef<DocumentData | undefined> | Ref<DocumentData | undefined>>(
    'documents:content:data',
    localDocumentData,
);
provide('documents:content:data', documentData);

const options = computed(() => ({
    nodeMapping: {
        penaltyCalculator: () => h(PenaltyCalculatorContentView),
    },
}));
</script>

<template>
    <TiptapContentRenderer
        class="tiptap hyphen-auto prose prose-sm max-w-full min-w-full break-words whitespace-pre-line sm:prose-base lg:prose-lg dark:prose-invert"
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
        :extensions="[...builtInExtensions, ...props.extensions]"
        :options="options"
        :value="props.value ? (Struct.toJson(props.value) as JSONContent) : undefined"
        v-bind="$attrs"
    />
</template>
