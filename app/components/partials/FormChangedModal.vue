<script lang="ts" setup>
import type { ButtonProps } from '@nuxt/ui';

const props = withDefaults(
    defineProps<{
        title?: string;
        description?: string;
        confirmLabel?: string;
        cancelLabel?: string;
        color?: ButtonProps['color'];
    }>(),
    {
        title: undefined,
        description: undefined,
        confirmLabel: undefined,
        cancelLabel: undefined,
        color: 'error',
    },
);

defineEmits<{
    (e: 'close', v: boolean): void;
}>();
</script>

<template>
    <UModal
        :overlay="false"
        :title="props.title ?? $t('components.partials.form_changed.title')"
        :description="props.description ?? $t('components.partials.form_changed.description')"
        :ui="{ content: 'z-[100]' }"
        @update:model-value="$emit('close', false)"
    >
        <template #footer>
            <UButton :color="props.color" :label="props.confirmLabel ?? $t('common.leave')" @click="$emit('close', true)" />
            <UButton color="neutral" :label="props.cancelLabel ?? $t('common.back')" @click="$emit('close', false)" />
        </template>
    </UModal>
</template>
