<script lang="ts" setup>
import type { JSONNode } from '~~/gen/ts/resources/common/content/content';
import HTMLContentRenderer from './HTMLContentRenderer.vue';

const props = defineProps<{
    value: JSONNode;
}>();

const contentRef = useTemplateRef('contentRef');

function disableCheckboxes(): void {
    console.log('disableCheckboxes', contentRef.value);
    if (contentRef.value === null) {
        return;
    }

    const checkboxes: NodeListOf<HTMLInputElement> = contentRef.value.querySelectorAll('input[type=checkbox]');
    checkboxes.forEach((el) => {
        el.setAttribute('disabled', 'disabled');
        el.classList.add('form-checkbox');
    });
}

watchOnce(
    () => props.value,
    () => useTimeoutFn(disableCheckboxes, 50),
);

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <div ref="contentRef">
        <HTMLContentRenderer
            v-bind="$attrs"
            :value="value"
            class="tiptap ProseMirror prose prose-sm sm:prose-base lg:prose-lg xl:prose-2xl dark:prose-invert min-w-full max-w-full break-words"
        />
    </div>
</template>
