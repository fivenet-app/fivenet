<script lang="ts" setup>
import { PaginationResponse } from '~~/gen/ts/resources/common/database/database';

const props = defineProps<{
    pagination: undefined | PaginationResponse;
}>();

defineEmits<{
    (e: 'offsetChange', offset: bigint): void;
}>();

const offset = computed(() => props.pagination?.offset ?? 0n);
const total = computed(() => props.pagination?.totalCount ?? 0n);
const pageSize = computed(() => props.pagination?.pageSize ?? 0n);
const end = computed(() => props.pagination?.end ?? 0n);

function calculateOffset(): bigint {
    const o = offset.value - pageSize.value;
    if (o < 0) {
        return 0n;
    }
    return o;
}
</script>

<template>
    <nav class="flex items-center justify-between px-4 py-3 border-t sm:px-6" aria-label="Pagination">
        <div class="hidden sm:block">
            <p
                class="text-sm text-gray-300"
                v-html="
                    $t('components.partials.table_pagination.showing_results', [
                        (total === 0n ? offset : offset + 1n).toString(),
                        end.toString(),
                        total.toString(),
                    ])
                "
            />
        </div>
        <div class="flex justify-between flex-1 sm:justify-end">
            <button
                :disabled="offset <= 0n"
                v-on:click="$emit('offsetChange', calculateOffset())"
                type="button"
                :class="[
                    offset <= 0n
                        ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                        : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                    'relative inline-flex items-center px-3 py-2 ml-3 text-sm font-semibold rounded-md cursor-pointer text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2',
                ]"
            >
                {{ $t('common.previous') }}
            </button>
            <button
                :disabled="total - end <= 0"
                v-on:click="$emit('offsetChange', end)"
                type="button"
                :class="[
                    total - end <= 0
                        ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                        : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                    'relative inline-flex items-center px-3 py-2 ml-3 text-sm font-semibold rounded-md cursor-pointer text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2',
                ]"
            >
                {{ $t('common.next') }}
            </button>
        </div>
    </nav>
</template>
