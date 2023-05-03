<script lang="ts" setup>
import { TemplateSchema } from '@fivenet/gen/resources/documents/templates_pb';

const props = defineProps<{
    modelValue: TemplateSchema,
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', payload: TemplateSchema): void,
}>();
</script>

<template>
    Test
</template>
