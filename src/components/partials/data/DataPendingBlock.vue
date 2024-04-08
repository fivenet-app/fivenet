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
    <UButton
        disabled
        block
        class="inline-flex w-full flex-col items-center gap-1 rounded-lg border border-dotted border-base-300 p-8 text-center hover:border-base-400 focus:ring-1 focus:ring-offset-2"
    >
        <PauseIcon v-if="paused" class="mx-auto size-12" />
        <RefreshIcon v-else class="mx-auto size-12 animate-spin" />
        <span class="block text-sm font-semibold">
            {{ message }}
        </span>
    </UButton>
</template>
