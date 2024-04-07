<script lang="ts" setup>
import { RefreshIcon } from 'mdi-vue3';
</script>

<template>
    <div class="absolute left-1/2 top-10 z-[800] inline-flex -translate-x-1/2 -translate-y-1/2 flex-col">
        <div
            class="text-primary-500 inset-0 inline-flex items-center justify-center rounded-md border border-black/20 bg-neutral-50 bg-clip-padding p-1.5 hover:bg-[#f4f4f4]"
        >
            <div class="flex items-center gap-1">
                <RefreshIcon class="size-6 animate-spin" />
                {{ $t('components.livemap.restarting_datastream') }}
            </div>
        </div>
    </div>
</template>
