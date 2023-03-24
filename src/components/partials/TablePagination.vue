<script lang="ts" setup>
defineProps({
    offset: {
        required: true,
        type: Number,
    },
    entries: {
        required: true,
        type: Number,
    },
    end: {
        required: true,
        type: Number,
    },
    total: {
        required: true,
        type: Number,
    },
    callback: {
        required: true,
        type: Function,
    }
});
</script>

<template>
    <nav class="flex items-center justify-between border-t px-4 py-3 sm:px-6" aria-label="Pagination">
        <div class="hidden sm:block">
            <p class="text-sm text-gray-300">
                Showing
                {{ ' ' }}
                <span class="font-medium text-neutral">{{ total == 0 ? offset : offset + 1 }}</span>
                {{ ' ' }}
                to
                {{ ' ' }}
                <span class="font-medium text-neutral">{{ end }}</span>
                {{ ' ' }}
                of
                {{ ' ' }}
                <span class="font-medium text-neutral">{{ total }}</span>
                {{ ' ' }}
                results
            </p>
        </div>
        <div class="flex flex-1 justify-between sm:justify-end">
            <button :class="[offset <= 0 ? 'disabled' : '']" :disabled="offset <= 0" v-on:click="callback(offset - entries)"
                type="button"
                class="relative inline-flex items-center rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral shadow-sm hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500 cursor-pointer">Previous</button>
            <button :class="[offset >= total ? 'disabled' : '']" :disabled="(offset + end) >= total"
                v-on:click="callback(end)" type="button"
                class="relative ml-3 inline-flex items-center rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral shadow-sm hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500 cursor-pointer">Next</button>
        </div>
    </nav>
</template>
