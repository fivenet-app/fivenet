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
            <div class="w-full border-t-2 border-neutral/20"></div>
        </div>
        <div v-if="label" class="relative mt-2 flex justify-center">
            <span class="rounded-lg border-2 border-neutral/20 px-3 text-sm font-semibold leading-6" :class="labelClass">{{
                label
            }}</span>
        </div>
    </div>
</template>
