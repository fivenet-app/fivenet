<script lang="ts" setup>
import { useThrottleFn } from '@vueuse/core';
import { RefreshIcon } from 'mdi-vue3';
import { PaginationResponse } from '~~/gen/ts/resources/common/database/database';

const props = defineProps<{
    pagination: undefined | PaginationResponse;
    refresh?: () => Promise<any>;
}>();

defineEmits<{
    (e: 'offsetChange', offset: bigint): void;
}>();

const offset = computed(() => props.pagination?.offset ?? 0n);
const total = computed(() => props.pagination?.totalCount ?? 0n);
const start = computed(() => (total.value === 0n ? offset.value : offset.value + 1n));
const pageSize = computed(() => props.pagination?.pageSize ?? 0n);
const end = computed(() => props.pagination?.end ?? 0n);

function calculateOffset(): bigint {
    const o = offset.value - pageSize.value;
    if (o < 0) {
        return 0n;
    }
    return o;
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async () => {
    if (props.refresh === undefined) {
        return;
    }

    canSubmit.value = false;
    await props.refresh().finally(() => setTimeout(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <nav class="flex items-center justify-between px-4 py-3 border-t sm:px-6" aria-label="Pagination">
        <div class="hidden sm:block">
            <I18nT keypath="components.partials.table_pagination.showing_results" tag="p" class="text-sm text-gray-300">
                <template #start>
                    <span class="font-medium text-neutral">
                        {{ start.toString() }}
                    </span>
                </template>
                <template #end>
                    <span class="font-medium text-neutral">
                        {{ end.toString() }}
                    </span>
                </template>
                <template #total>
                    <span class="font-medium text-neutral">
                        {{ total.toString() }}
                    </span>
                </template>
            </I18nT>
        </div>
        <div class="flex justify-between flex-1 sm:justify-end">
            <button
                v-if="refresh !== undefined"
                type="button"
                class="bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500 relative inline-flex items-center px-3 py-2 ml-3 text-sm font-semibold rounded-md cursor-pointer text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                :disabled="!canSubmit"
                :class="
                    !canSubmit
                        ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                        : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500'
                "
                @click="onSubmitThrottle()"
            >
                <RefreshIcon class="h-5 w-5" :class="!canSubmit ? 'animate-spin' : ''" />
            </button>
            <button
                :disabled="offset <= 0n"
                type="button"
                :class="[
                    offset <= 0n
                        ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                        : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                    'relative inline-flex items-center px-3 py-2 ml-3 text-sm font-semibold rounded-md cursor-pointer text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2',
                ]"
                @click="$emit('offsetChange', calculateOffset())"
            >
                {{ $t('common.previous') }}
            </button>
            <button
                :disabled="total - end <= 0"
                type="button"
                :class="[
                    total - end <= 0
                        ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                        : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                    'relative inline-flex items-center px-3 py-2 ml-3 text-sm font-semibold rounded-md cursor-pointer text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2',
                ]"
                @click="$emit('offsetChange', end)"
            >
                {{ $t('common.next') }}
            </button>
        </div>
    </nav>
</template>
