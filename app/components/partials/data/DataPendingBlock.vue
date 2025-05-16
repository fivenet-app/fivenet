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
        class="my-2 inline-flex w-full flex-col items-center gap-2 rounded-lg border border-dotted border-base-300 p-6 text-center hover:border-base-400"
        disabled
        block
    >
        <UIcon v-if="paused" class="mx-auto size-12" name="i-mdi-pause" />
        <UIcon v-else class="mx-auto size-12 animate-spin" name="i-mdi-refresh" />

        <span class="block text-sm font-semibold">
            {{ message }}
        </span>
    </UButton>
</template>
