<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
    props: {
        start: {
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
        }
    },
});
</script>

<template>
    <nav class="flex items-center justify-between border-t px-4 py-3 sm:px-6" aria-label="Pagination">
        <div class="hidden sm:block">
            <p class="text-sm text-gray-300">
                Showing
                {{ ' ' }}
                <span class="font-medium text-white">{{ start + 1 }}</span>
                {{ ' ' }}
                to
                {{ ' ' }}
                <span class="font-medium text-white">{{ end }}</span>
                {{ ' ' }}
                of
                {{ ' ' }}
                <span class="font-medium text-white">{{ total }}</span>
                {{ ' ' }}
                results
            </p>
        </div>
        <div class="flex flex-1 justify-between sm:justify-end">
            <button :class="[start <= 0 ? 'disabled' : '']" :disabled="start <= 0" v-on:click="callback(start - end)"
                class="relative inline-flex items-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus-visible:outline-offset-0">Previous</button>
            <button :class="[start >= total ? 'disabled' : '']" :disabled="start >= total" v-on:click="callback(end)"
                class="relative ml-3 inline-flex items-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus-visible:outline-offset-0">Next</button>
        </div>
    </nav>
</template>
