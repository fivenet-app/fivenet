<script lang="ts" setup>
withDefaults(
    defineProps<{
        open: boolean;
        dialogClass?: unknown;
        title?: string;
    }>(),
    {
        dialogClass: '' as any,
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
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template v-if="title" #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ title }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="$emit('close')" />
                </div>
            </template>

            <div>
                <slot />
            </div>

            <template #footer>
                <UButton block class="flex-1" color="black" @click="$emit('close')">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </UModal>
</template>
