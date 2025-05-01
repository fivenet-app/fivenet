<script lang="ts" setup>
import type { ButtonColor } from '#ui/types';

withDefaults(
    defineProps<{
        title?: string;
        description?: string;
        cancel?: () => Promise<unknown> | unknown;
        confirm: () => Promise<unknown> | unknown;
        icon?: string;
        color?: ButtonColor;
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

const { isOpen } = useModal();
</script>

<template>
    <UDashboardModal
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
                color="white"
                :label="$t('common.cancel')"
                @click="
                    if (cancel) {
                        cancel();
                    }
                    isOpen = false;
                "
            />
        </template>
    </UDashboardModal>
</template>
