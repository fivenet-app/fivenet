<script lang="ts" setup>
import { PaginationResponse } from '@arpanet/gen/resources/common/database/database_pb';
import { computed } from 'vue';

const props = defineProps({
    pagination: {
        required: true,
        type: PaginationResponse,
    },
    callback: {
        required: true,
        type: Function,
    },
});

const offset = computed(() => props.pagination.getOffset()!);
const total = computed(() => props.pagination.getTotalCount()!);
const pageSize = computed(() => props.pagination.getPageSize()!);
const end = computed(() => props.pagination.getEnd()!);
</script>

<template>
    <nav class="flex items-center justify-between px-4 py-3 border-t sm:px-6" aria-label="Pagination">
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
        <div class="flex justify-between flex-1 sm:justify-end">
            <button :class="[offset <= 0 ? 'disabled' : '']" :disabled="offset <= 0" v-on:click="callback(offset - pageSize)"
                type="button"
                class="relative inline-flex items-center px-3 py-2 text-sm font-semibold rounded-md cursor-pointer bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500">Previous</button>
            <button :class="[offset >= total ? 'disabled' : '']" :disabled="(end + offset) >= total"
                v-on:click="callback(end)" type="button"
                class="relative inline-flex items-center px-3 py-2 ml-3 text-sm font-semibold rounded-md cursor-pointer bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500">Next</button>
        </div>
    </nav>
</template>
