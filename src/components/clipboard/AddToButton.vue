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
        class="fixed bottom-16 right-6 z-10 flex h-12 w-12 items-center justify-center rounded-full bg-primary-500 text-neutral shadow-float hover:bg-primary-400"
        @click="callback()"
    >
        <ClipboardPlusIcon class="w-10 h-auto" />
    </button>
</template>
