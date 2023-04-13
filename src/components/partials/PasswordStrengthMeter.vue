<script setup lang="ts">
import zxcvbn from 'zxcvbn';

const props = defineProps<{
    input: string,
    showFeedback?: boolean,
}>();

const result = computed(() => zxcvbn(props.input));

const percent = ref<number>(0);
const feedback = ref<string>('');
const color = ref<string>('bg-base-700');

watch(result, () => {
    percent.value = result.value.score * 100 / 4;
    feedback.value = result.value.feedback.warning;


    switch (result.value.score) {
        case 0: color.value = 'bg-error-500'; break;
        case 1: color.value = 'bg-error-500'; break;
        case 2: color.value = 'bg-warn-500'; break;
        case 3: color.value = 'bg-success-500'; break;
        case 4: color.value = 'bg-success-500'; break;
        default: color.value = 'bg-base-700'; break;
    }

    if (props.input === '') color.value = 'bg-base-700';
})
</script>

<template>
    <div>
        <div :class="['h-2 w-full rounded-full transition-colors', color]"></div>
        <p v-if="showFeedback" class="text-sm text-base-300 my-1">{{ feedback }}</p>
    </div>
</template>
