<script lang="ts" setup>
import { CodeDiff as VCodeDiff, type CodeDiffProps } from 'v-code-diff';
import { computed, useAttrs } from 'vue';

defineOptions({
    inheritAttrs: false,
});

const props = defineProps<CodeDiffProps>();
const attrs = useAttrs();

const forwardedAttrs = computed(() => {
    const rest: Record<string, unknown> = {};

    for (const [key, value] of Object.entries(attrs)) {
        if (key !== 'class' && key !== 'style') {
            rest[key] = value;
        }
    }

    return rest;
});

const codeDiffAttrs = computed(() => ({
    ...forwardedAttrs.value,
    ...props,
}));
</script>

<template>
    <ClientOnly>
        <VCodeDiff v-bind="codeDiffAttrs" :class="attrs.class" :style="attrs.style" />
    </ClientOnly>
</template>
