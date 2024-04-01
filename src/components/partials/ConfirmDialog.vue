<script lang="ts" setup>
withDefaults(
    defineProps<{
        open: boolean;
        title?: string;
        description?: string;
        cancel: (data?: any) => void;
        confirm: (data?: any) => void;
        icon?: string;
    }>(),
    {
        title: undefined,
        description: undefined,
        icon: 'i-heroicons-exclamation-circle',
    },
);

defineEmits<{
    (e: 'close'): void;
}>();
</script>

<template>
    <UDashboardModal
        :model-value="open"
        :title="title ?? $t('components.partials.confirm_dialog.title')"
        :description="description ?? $t('components.partials.confirm_dialog.description')"
        :icon="icon"
        :ui="{
            icon: { base: 'text-red-500 dark:text-red-400' } as any,
            footer: { base: 'ml-16' } as any,
        }"
        @update:model-value="
            cancel();
            $emit('close');
        "
    >
        <template #footer>
            <UButton
                color="red"
                :label="$t('common.confirm')"
                @click="
                    confirm();
                    $emit('close');
                "
            />
            <UButton
                color="white"
                :label="$t('common.cancel')"
                @click="
                    cancel();
                    $emit('close');
                "
            />
        </template>
    </UDashboardModal>
</template>
