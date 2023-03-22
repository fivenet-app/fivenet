<script lang="ts" setup>
defineProps({
    title: {
        required: true,
        type: String,
    },
});
</script>

<template>
    <header class="bg-gray-700 shadow-sm">
        <div class="mx-auto max-w-7xl py-4 px-4 sm:px-6 lg:px-8">
            <div class="grid gap-4 grid-cols-2 grid-rows-1">
                <h1 class="text-lg font-semibold leading-6 text-white">{{ title }}</h1>
            </div>
        </div>
    </header>
</template>
