<script lang="ts" setup>
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
        class="inline-flex w-full flex-col items-center gap-1 rounded-lg border border-dotted border-base-300 p-6 text-center hover:border-base-400"
    >
        <UIcon v-if="paused" name="i-mdi-pause" class="mx-auto size-12" />
        <UIcon v-else name="i-mdi-refresh" class="mx-auto size-12 animate-spin" />

        <span class="block text-sm font-semibold">
            {{ message }}
        </span>
    </UButton>
</template>
