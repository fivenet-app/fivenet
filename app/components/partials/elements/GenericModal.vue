<script lang="ts" setup>
import type { ClassProp } from '~/typings';

withDefaults(
    defineProps<{
        open: boolean;
        dialogClass?: ClassProp;
        title?: string;
    }>(),
    {
        dialogClass: '',
        title: undefined,
    },
);

defineEmits<{
    (e: 'close'): void;
}>();
</script>

<template>
    <UModal
        :class="dialogClass"
        :model-value="open"
        :ui="{ width: 'w-full sm:max-w-5xl' }"
        @update:model-value="$emit('close')"
    >
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template v-if="title" #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ title }}
                    </h3>

                    <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="$emit('close')" />
                </div>
            </template>

            <div>
                <slot />
            </div>

            <template #footer>
                <UButton class="flex-1" block color="black" @click="$emit('close')">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </UModal>
</template>
