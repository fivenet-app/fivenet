<script lang="ts" setup>
import { ArrowPathIcon } from '@heroicons/vue/24/solid';

defineProps({
    error: {
        type: String,
        required: false,
    },
    message: {
        type: String,
        required: false,
        default: 'Loading your data...',
    },
});
</script>

<template>
    <button type="button" disabled
        class="relative block w-full p-12 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2">
        <ArrowPathIcon class="w-12 h-12 mx-auto text-neutral animate-spin" />
        <span class="block mt-2 text-sm font-semibold text-gray-300">
            {{ message }}
        </span>
    </button>
</template>
