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
        class="relative block w-full rounded-lg border-2 border-dotted border-base-300 p-8 text-center hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
    >
        <PauseIcon v-if="paused" class="mx-auto h-12 w-12 text-neutral" />
        <RefreshIcon v-else class="mx-auto h-12 w-12 animate-spin text-neutral" />
        <span class="mt-2 block text-sm font-semibold text-gray-300">
            {{ message }}
        </span>
    </button>
</template>
