<script lang="ts" setup>
import { RefreshIcon } from 'mdi-vue3';
</script>

<template>
    <div class="absolute z-[800] top-10 left-1/2 transform -translate-x-1/2 -translate-y-1/2 inline-flex flex-col">
        <div
            class="inset-0 inline-flex items-center justify-center rounded-md border-2 border-black/20 bg-neutral bg-clip-padding p-1.5 text-black text-primary-500 hover:bg-[#f4f4f4] focus:outline-none"
        >
            <div class="flex items-center gap-1">
                <RefreshIcon class="h-6 w-6 animate-spin" />
                {{ $t('components.livemap.restarting_datastream') }}
            </div>
        </div>
    </div>
</template>
