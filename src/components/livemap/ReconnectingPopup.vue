<script lang="ts" setup>
import { RefreshIcon } from 'mdi-vue3';
</script>

<template>
    <div class="absolute left-1/2 top-10 z-[800] inline-flex -translate-x-1/2 -translate-y-1/2 flex-col">
        <div
            class="inset-0 inline-flex items-center justify-center rounded-md border-2 border-black/20 bg-neutral bg-clip-padding p-1.5 text-primary-500 hover:bg-[#f4f4f4] focus:outline-none"
        >
            <div class="flex items-center gap-1">
                <RefreshIcon class="size-6 animate-spin" aria-hidden="true" />
                {{ $t('components.livemap.restarting_datastream') }}
            </div>
        </div>
    </div>
</template>
