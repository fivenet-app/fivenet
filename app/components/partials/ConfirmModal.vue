<script lang="ts" setup>
import type { ButtonProps } from '@nuxt/ui';

const props = withDefaults(
    defineProps<{
        title?: string;
        description?: string;
        cancel?: () => Promise<unknown> | unknown;
        confirm: () => Promise<unknown> | unknown;
        icon?: string;
        color?: ButtonProps['color'];
        iconClass?: string;
    }>(),
    {
        title: undefined,
        description: undefined,
        cancel: undefined,
        icon: 'i-mdi-warning-circle',
        color: 'error',
        iconClass: 'text-red-500 dark:text-red-400',
    },
);

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

async function handleConfirm(): Promise<void> {
    await props.confirm();
    emit('close', true);
}

async function handleCancel(): Promise<void> {
    await props.cancel?.();
    emit('close', false);
}
</script>

<template>
    <UModal
        :title="props.title ?? $t('components.partials.confirm_dialog.title')"
        :description="props.description ?? $t('components.partials.confirm_dialog.description')"
        @update:model-value="props.cancel && props.cancel()"
    >
        <template #footer>
            <UButton :color="props.color" :label="$t('common.confirm')" @click="handleConfirm" />
            <UButton color="neutral" :label="$t('common.cancel')" @click="handleCancel" />
        </template>
    </UModal>
</template>
