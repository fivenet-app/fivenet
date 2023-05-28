<script lang="ts" setup>
import { PlusIcon } from '@heroicons/vue/24/solid';

withDefaults(
    defineProps<{
        callback: Function;
        title?: string;
    }>(),
    {
        title: '',
    }
);
</script>

<template>
    <button
        type="button"
        :title="title"
        @click="callback()"
        class="fixed flex items-center justify-center w-12 h-12 rounded-full z-90 bottom-24 right-8 bg-primary-500 shadow-float text-neutral hover:bg-primary-400"
    >
        <PlusIcon class="w-10 h-auto" />
    </button>
</template>
