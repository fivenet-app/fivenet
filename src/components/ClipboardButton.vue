<script lang="ts" setup>
import { PlusIcon } from '@heroicons/vue/24/solid';

defineProps({
    callback: {
        required: true,
        type: Function,
    }
});
</script>

<template>
    <button title="Actions" @click="callback()"
        class="fixed z-90 bottom-10 right-8 bg-blue-600 w-14 h-14 rounded-full drop-shadow-lg flex justify-center items-center text-white text-4xl hover:bg-blue-700 hover:drop-shadow-2xl">
        <PlusIcon />
    </button>
</template>
