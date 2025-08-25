<script setup lang="ts">
defineProps<{
    label?: string;
}>();
</script>

<template>
    <div class="absolute top-10 left-1/2 z-10 inline-flex -translate-x-1/2 -translate-y-1/2 flex-col md:top-8">
        <div
            class="inset-0 inline-flex items-center justify-center rounded-md border border-black/20 bg-default bg-clip-padding p-1.5 text-primary-500"
        >
            <div class="mx-0.5 flex items-center gap-1 text-xs">
                <UIcon class="size-4 shrink-0 animate-spin" name="i-mdi-refresh" />
                <span>{{ label }}</span>
            </div>
        </div>
    </div>
</template>
