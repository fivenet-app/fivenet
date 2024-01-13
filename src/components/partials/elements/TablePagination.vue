<script lang="ts" setup>
import { useThrottleFn } from '@vueuse/core';
import { RefreshIcon } from 'mdi-vue3';
import { Popover, PopoverButton, PopoverPanel } from '@headlessui/vue';
import { Float } from '@headlessui-float/vue';
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
const pageSize = computed(() => props.pagination?.pageSize ?? 0n);
const end = computed(() => props.pagination?.end ?? 0n);

const totalPages = computed(() => bigIntCeil(total.value, pageSize.value));
const currentPage = computed(() => offset.value / pageSize.value + 1n);

function calculateOffset(pageCount: number): bigint {
    const pageC = BigInt(pageCount ?? 1);
    if (pageC > totalPages.value) {
        return (totalPages.value - 1n) * pageSize.value;
    } else if (pageC < 1) {
        return 0n;
    }

    const o = pageSize.value * (pageC - 1n);
    console.log(o, pageSize.value, pageCount, offset.value);
    if (o < 0) {
        return 0n;
    }
    return o;
}

const paginationLength = 2;
const beforePages = computed(() => {
    const curPage = parseInt(currentPage.value.toString());
    const start = curPage - paginationLength;

    if (curPage <= 1) {
        return [];
    } else if (curPage <= 2) {
        return [1];
    }

    return range(paginationLength, start);
});
const afterPages = computed(() => {
    const curPage = parseInt(currentPage.value.toString());

    if (currentPage.value >= totalPages.value) {
        return [];
    } else if (currentPage.value >= totalPages.value - 1n) {
        return [parseInt(totalPages.value.toString())];
    }

    return range(2, curPage + 1);
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async () => {
    if (props.refresh === undefined) {
        return;
    }

    canSubmit.value = false;
    await props.refresh().finally(() => setTimeout(() => (canSubmit.value = true), 400));
}, 1000);

const pageNumber = ref(currentPage.value.toString());
</script>

<template>
    <nav class="flex items-center justify-between px-4 py-3 border-t sm:px-6" aria-label="Pagination">
        <div class="hidden sm:block">
            <I18nT keypath="components.partials.table_pagination.page_count" tag="p" class="text-sm text-gray-300">
                <template #current>
                    <span class="font-medium text-neutral">
                        {{ currentPage.toString() }}
                    </span>
                </template>
                <template #total>
                    <span class="font-medium text-neutral">
                        {{ totalPages.toString() }}
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

            <nav class="isolate inline-flex -space-x-px rounded-md shadow-sm" aria-label="Pagination">
                <button
                    :disabled="offset <= 0n"
                    type="button"
                    :class="[
                        offset <= 0n
                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                            : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                        'relative inline-flex items-center px-3 py-2 ml-3 text-sm font-semibold rounded-l-md cursor-pointer text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2',
                    ]"
                    @click="$emit('offsetChange', calculateOffset(parseInt((currentPage - 1n).toString())))"
                >
                    {{ $t('common.previous') }}
                </button>

                <button
                    v-for="page in beforePages"
                    :key="page"
                    type="button"
                    class="relative inline-flex items-center px-4 py-2 text-sm font-semibold text-neutral bg-accent-500 hover:bg-base-400 focus:z-20 focus:outline-offset-0"
                    @click="$emit('offsetChange', calculateOffset(page))"
                >
                    {{ page }}
                </button>

                <button
                    type="button"
                    class="relative inline-flex items-center px-4 py-2 text-sm font-semibold text-neutral bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500 underline"
                    disabled
                >
                    {{ currentPage }}
                </button>

                <button
                    v-for="page in afterPages"
                    :key="page"
                    type="button"
                    class="relative inline-flex items-center px-4 py-2 text-sm font-semibold text-neutral bg-accent-500 hover:bg-base-400 focus:z-20 focus:outline-offset-0"
                    @click="$emit('offsetChange', calculateOffset(page))"
                >
                    {{ page }}
                </button>

                <Popover v-if="totalPages > 4n" class="relative">
                    <Float portal placement="top-start" :offset="12">
                        <PopoverButton
                            class="relative inline-flex items-center px-4 py-2 text-sm font-semibold text-neutral bg-accent-500 hover:bg-base-400 focus:z-20 focus:outline-offset-0"
                            @click="pageNumber = ''"
                        >
                            ...
                        </PopoverButton>

                        <PopoverPanel
                            focus
                            class="absolute z-5 w-24 max-w-24 min-w-fit text-sm text-gray-400 transition-opacity bg-gray-800 border border-gray-600 rounded-lg shadow-sm"
                        >
                            <div class="p-3">
                                <form @submit.prevent="$emit('offsetChange', calculateOffset(parseInt(pageNumber)))">
                                    <input
                                        v-model="pageNumber"
                                        type="number"
                                        min="1"
                                        :max="parseInt(totalPages.toString())"
                                        class="remove-arrow block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        name="page_number"
                                        :placeholder="$t('common.page')"
                                    />
                                </form>
                            </div>
                        </PopoverPanel>
                    </Float>
                </Popover>

                <button
                    v-if="currentPage <= totalPages - 2n"
                    type="button"
                    class="relative inline-flex items-center px-4 py-2 text-sm font-semibold text-neutral bg-accent-500 hover:bg-base-400 focus:z-20 focus:outline-offset-0"
                    @click="$emit('offsetChange', calculateOffset(parseInt(totalPages.toString())))"
                >
                    {{ totalPages }}
                </button>

                <button
                    :disabled="total - end <= 0n"
                    type="button"
                    :class="[
                        total - end <= 0n
                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                            : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                        'relative inline-flex items-center px-3 py-2 ml-3 text-sm font-semibold rounded-r-md cursor-pointer text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2',
                    ]"
                    @click="$emit('offsetChange', end)"
                >
                    {{ $t('common.next') }}
                </button>
            </nav>
        </div>
    </nav>
</template>
