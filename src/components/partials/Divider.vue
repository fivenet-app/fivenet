<script lang="ts" setup>
defineProps<{
    label?: string;
}>();
</script>

<template>
    <div :class="[label ? 'mb-2' : 'mb-4', 'relative']">
        <div class="absolute inset-0 flex items-center" aria-hidden="true">
            <div class="w-full border-t border-gray-300"></div>
        </div>
        <div v-if="label" class="relative flex justify-center">
            <span class="bg-primary-900 px-3 text-base font-semibold leading-6 text-gray-200">{{ label }}</span>
        </div>
    </div>
</template>
