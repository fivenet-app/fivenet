<script lang="ts" setup>
withDefaults(
    defineProps<{
        title?: string;
        description?: string;
        cancel?: () => Promise<unknown> | unknown;
        confirm: () => Promise<unknown> | unknown;
        icon?: string;
    }>(),
    {
        title: undefined,
        description: undefined,
        cancel: undefined,
        icon: 'i-mdi-warning-circle',
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
            icon: { base: 'text-red-500 dark:text-red-400' } as any,
            footer: { base: 'ml-16' },
        }"
        @update:model-value="cancel && cancel()"
    >
        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton
                    color="red"
                    :label="$t('common.confirm')"
                    @click="
                        confirm();
                        isOpen = false;
                    "
                />
                <UButton
                    block
                    class="flex-1"
                    color="white"
                    :label="$t('common.cancel')"
                    @click="
                        if (cancel) {
                            cancel();
                        }
                        isOpen = false;
                    "
                />
            </UButtonGroup>
        </template>
    </UDashboardModal>
</template>
