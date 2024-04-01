<script lang="ts" setup>
withDefaults(
    defineProps<{
        open: boolean;
        dialogClass?: unknown;
        unmount?: boolean;
        title?: string;
    }>(),
    {
        dialogClass: '' as any,
        unmount: true,
        title: undefined,
    },
);

defineEmits<{
    (e: 'close'): void;
}>();
</script>

<template>
    <UModal
        :model-value="open"
        :class="dialogClass"
        :ui="{ width: 'w-full sm:max-w-5xl' }"
        @update:model-value="$emit('close')"
    >
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }" :unmount="unmount">
            <template v-if="title" #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ title }}
                    </h3>

                    <UButton
                        color="gray"
                        variant="ghost"
                        icon="i-heroicons-x-mark-20-solid"
                        class="-my-1"
                        @click="$emit('close')"
                    />
                </div>
            </template>

            <div>
                <slot />
            </div>

            <div class="mt-5 gap-2 sm:mt-4 sm:flex">
                <UButton block @click="$emit('close')">
                    {{ $t('common.close', 1) }}
                </UButton>
            </div>
        </UCard>
    </UModal>
</template>
