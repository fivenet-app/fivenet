<script lang="ts" setup>
import { PaginationResponse } from '@fivenet/gen/resources/common/database/database_pb';

const props = defineProps<{
    pagination: undefined | PaginationResponse,
}>();

defineEmits<{
    (e: 'offsetChange', offset: number): void,
}>();

const offset = computed(() => props.pagination?.getOffset() ?? 0);
const total = computed(() => props.pagination?.getTotalCount() ?? 0);
const pageSize = computed(() => props.pagination?.getPageSize() ?? 0);
const end = computed(() => props.pagination?.getEnd() ?? 0);
</script>

<template>
    <nav class="flex items-center justify-between px-4 py-3 border-t sm:px-6" aria-label="Pagination">
        <div class="hidden sm:block">
            <p class="text-sm text-gray-300"
                v-html="$t('components.partials.table_pagination.showing_results', [total == 0 ? offset : offset + 1, end, total])" />
        </div>
        <div class="flex justify-between flex-1 sm:justify-end">
            <button :class="[offset <= 0 ? 'disabled' : '']" :disabled="offset <= 0"
                v-on:click="$emit('offsetChange', offset - pageSize)" type="button"
                class="relative inline-flex items-center px-3 py-2 text-sm font-semibold rounded-md cursor-pointer bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500">{{
                    $t('common.previous') }}
            </button>
            <button :class="[offset >= total ? 'disabled' : '']" :disabled="(end + offset) >= total"
                v-on:click="$emit('offsetChange', end)" type="button"
                class="relative inline-flex items-center px-3 py-2 ml-3 text-sm font-semibold rounded-md cursor-pointer bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500">{{
                    $t('common.next') }}
            </button>
        </div>
    </nav>
</template>
