<script lang="ts" setup>
import { ClipboardPlusIcon } from 'mdi-vue3';

withDefaults(
    defineProps<{
        callback: Function;
        title?: string;
    }>(),
    {
        title: '',
    },
);
</script>

<template>
    <button
        type="button"
        :title="title"
        class="fixed bottom-16 right-6 z-10 flex size-12 items-center justify-center rounded-full bg-accent-500 text-neutral shadow-float hover:bg-accent-400"
        @click="callback()"
    >
        <ClipboardPlusIcon class="h-auto w-10" aria-hidden="true" />
    </button>
</template>
