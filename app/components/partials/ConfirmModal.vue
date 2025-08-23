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

const { isOpen } = useOverlay();
</script>

<template>
    <UModal
        :title="title ?? $t('components.partials.confirm_dialog.title')"
        :description="description ?? $t('components.partials.confirm_dialog.description')"
        :icon="icon"
        :ui="{
            icon: { base: iconClass },
            footer: { base: 'ml-16' },
        }"
        @update:model-value="cancel && cancel()"
    >
        <template #footer>
            <UButton
                :color="color"
                :label="$t('common.confirm')"
                @click="
                    confirm();
                    isOpen = false;
                "
            />
            <UButton
                color="neutral"
                :label="$t('common.cancel')"
                @click="
                    if (cancel) {
                        cancel();
                    }
                    isOpen = false;
                "
            />
        </template>
    </UModal>
</template>
