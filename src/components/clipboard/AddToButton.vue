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
        @click="callback()"
        class="fixed flex items-center justify-center w-12 h-12 rounded-full z-10 bottom-16 right-6 bg-primary-500 shadow-float text-neutral hover:bg-primary-400"
    >
        <ClipboardPlusIcon class="w-10 h-auto" />
    </button>
</template>
