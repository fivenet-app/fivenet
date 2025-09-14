<script lang="ts" setup>
import type { ButtonProps } from '@nuxt/ui';

withDefaults(
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

defineEmits<{
    (e: 'close', v: boolean): void;
}>();
</script>

<template>
    <UModal
        :title="title ?? $t('components.partials.confirm_dialog.title')"
        :description="description ?? $t('components.partials.confirm_dialog.description')"
        @update:model-value="cancel && cancel()"
    >
        <template #footer>
            <UButton
                :color="color"
                :label="$t('common.confirm')"
                @click="
                    confirm();
                    $emit('close', true);
                "
            />
            <UButton
                color="neutral"
                :label="$t('common.cancel')"
                @click="
                    if (cancel) {
                        cancel();
                    }
                    $emit('close', false);
                "
            />
        </template>
    </UModal>
</template>
