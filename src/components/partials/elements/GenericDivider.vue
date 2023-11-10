<script lang="ts" setup>
withDefaults(
    defineProps<{
        label?: string;
        labelClass?: unknown;
    }>(),
    {
        label: undefined,
        labelClass: 'bg-primary-900 text-gray-200' as any,
    },
);
</script>

<template>
    <div :class="[label ? 'mb-2' : 'mb-4', 'relative']">
        <div class="absolute inset-0 flex items-center" aria-hidden="true">
            <div class="w-full border-t border-gray-300"></div>
        </div>
        <div v-if="label" class="relative flex justify-center">
            <span class="px-3 text-base font-semibold leading-6" :class="labelClass">{{ label }}</span>
        </div>
    </div>
</template>
