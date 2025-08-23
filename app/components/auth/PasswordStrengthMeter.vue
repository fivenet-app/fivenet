<script setup lang="ts">
import type { ProgressProps } from '@nuxt/ui';
import { zxcvbn, zxcvbnOptions } from '@zxcvbn-ts/core';

const props = defineProps<{
    input: string;
    showFeedback?: boolean;
}>();

zxcvbnOptions.setOptions({
    maxLength: 70,
});

const percent = ref<number>(0);
const feedback = ref<string | null>('');
const color = ref<ProgressProps['color']>('neutral');

const result = computed(() => zxcvbn(props.input));

watch(result, () => {
    percent.value = (result.value.score * 100) / 3;
    feedback.value = result.value.feedback.warning;

    if (props.input.trimEnd() === '') {
        color.value = 'neutral';
        return;
    }

    switch (result.value.score) {
        case 0:
        case 1:
            color.value = 'error';
            break;
        case 2:
            color.value = 'warning';
            break;
        case 3:
        case 4:
            color.value = 'success';
            break;
        default:
            color.value = 'neutral';
            break;
    }
});
</script>

<template>
    <div>
        <!-- @vue-expect-error seems that the `color` prop is not using the `ProgressColor` type -->
        <UProgress :color="color" :value="percent" />

        <p v-if="showFeedback && feedback !== null" class="my-1 text-sm">
            {{ feedback }}
        </p>
    </div>
</template>
