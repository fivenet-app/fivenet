<script lang="ts" setup>
import { useVModel } from '@vueuse/core';

const props = withDefaults(
    defineProps<{
        modelValue?: string;
    }>(),
    {
        modelValue: '',
    },
);

const emit = defineEmits(['update:modelValue']);
const data = useVModel(props, 'modelValue', emit);
</script>

<template>
    <UInput v-model="data" @focusin="focusTablet(true)" @focusout="focusTablet(false)" />
</template>
