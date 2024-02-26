<script setup lang="ts">
import { zxcvbn, zxcvbnOptions } from '@zxcvbn-ts/core';

const props = defineProps<{
    input: string;
    showFeedback?: boolean;
}>();

zxcvbnOptions.setOptions({});

const percent = ref<number>(0);
const feedback = ref<string | null>('');
const color = ref<string>('bg-base-700');

const result = computed(() => zxcvbn(props.input));

watch(result, () => {
    percent.value = (result.value.score * 100) / 4;
    feedback.value = result.value.feedback.warning;

    switch (result.value.score) {
        case 0:
            color.value = 'bg-error-500';
            break;
        case 1:
            color.value = 'bg-error-500';
            break;
        case 2:
            color.value = 'bg-warn-500';
            break;
        case 3:
            color.value = 'bg-success-500';
            break;
        case 4:
            color.value = 'bg-success-500';
            break;
        default:
            color.value = 'bg-base-700';
            break;
    }

    if (props.input === '') {
        color.value = 'bg-base-700';
    }
});
</script>

<template>
    <div>
        <div :class="['h-2 w-full rounded-full transition-colors', color]"></div>
        <p v-if="showFeedback && feedback !== null" class="my-1 text-sm text-base-300">
            {{ feedback }}
        </p>
    </div>
</template>
