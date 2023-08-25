<script lang="ts" setup>
import { PauseIcon, RefreshIcon } from 'mdi-vue3';

withDefaults(
    defineProps<{
        message: string;
        paused?: boolean;
    }>(),
    {
        paused: false,
    },
);
</script>

<template>
    <button
        type="button"
        disabled
        class="mt-8 relative block w-full p-12 text-center border-2 border-dotted rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
    >
        <PauseIcon v-if="paused" class="w-12 h-12 mx-auto text-neutral" />
        <RefreshIcon v-else class="w-12 h-12 mx-auto text-neutral animate-spin" />
        <span class="block mt-2 text-sm font-semibold text-gray-300">
            {{ message }}
        </span>
    </button>
</template>
