<script lang="ts" setup>
import { useThrottleFn, useTimeoutFn } from '@vueuse/core';
import { RefreshIcon } from 'mdi-vue3';
import { Popover, PopoverButton, PopoverPanel } from '@headlessui/vue';
import { PaginationResponse } from '~~/gen/ts/resources/common/database/database';

const props = withDefaults(
    defineProps<{
        pagination: undefined | PaginationResponse;
        refresh?: () => Promise<any>;
        showBorder?: boolean;
    }>(),
    {
        refresh: undefined,
        showBorder: true,
    },
);

defineEmits<{
    (e: 'offsetChange', offset: number): void;
}>();

const offset = computed(() => props.pagination?.offset ?? 0);
const total = computed(() => props.pagination?.totalCount ?? 0);
const pageSize = computed(() => props.pagination?.pageSize ?? 1);
const end = computed(() => props.pagination?.end ?? 0);

const totalPages = computed(() => Math.ceil(total.value / pageSize.value));
const currentPage = computed(() => offset.value / pageSize.value + 1);

function calculateOffset(pageCount?: number): number {
    const pageC = pageCount ?? 1;
    if (pageC > totalPages.value) {
        return (totalPages.value - 1) * pageSize.value;
    } else if (pageC < 1) {
        return 0;
    }

    const o = pageSize.value * (pageC - 1);
    if (o < 0) {
        return 0;
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
    } else if (currentPage.value >= totalPages.value - 1) {
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
    await props.refresh().finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const pageNumber = ref(currentPage.value.toString());
</script>

<template>
    <nav
        v-if="pagination !== undefined"
        class="flex items-center justify-between px-4 py-3 sm:px-6"
        :class="showBorder ? 'border-t' : ''"
        aria-label="Pagination"
    >
        <div v-if="total > -1" class="hidden sm:block">
            <I18nT keypath="components.partials.table_pagination.page_count" tag="p" class="text-sm text-gray-300">
                <template #current>
                    <span class="font-medium">
                        {{ currentPage.toString() }}
                    </span>
                </template>
                <template #total>
                    <span class="font-medium">
                        {{ total.toString() }}
                    </span>
                </template>
                <template #maxPage>
                    <span class="font-medium">
                        {{ totalPages === 0 ? 1 : totalPages }}
                    </span>
                </template>
                <template #size>
                    <span class="font-medium">
                        {{ pageSize.toString() }}
                    </span>
                </template>
            </I18nT>
        </div>
        <div class="flex flex-1 justify-between sm:justify-end">
            <UButton
                v-if="refresh !== undefined"
                class="relative ml-3 inline-flex cursor-pointer items-center rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                :disabled="!canSubmit"
                :class="
                    !canSubmit
                        ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                        : 'bg-primary-500 hover:bg-primary-400'
                "
                @click="onSubmitThrottle()"
            >
                <RefreshIcon class="size-5" :class="!canSubmit ? 'animate-spin' : ''" />
            </UButton>

            <nav class="isolate inline-flex -space-x-px rounded-md shadow-sm" aria-label="Pagination">
                <template v-if="total > -1">
                    <UButton
                        :disabled="offset <= 0n"
                        :class="[
                            offset <= 0n
                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                : 'bg-primary-500 hover:bg-primary-400',
                            'relative ml-3 inline-flex cursor-pointer items-center rounded-l-md px-3 py-2 text-sm font-semibold focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2',
                        ]"
                        @click="$emit('offsetChange', calculateOffset(currentPage - 1))"
                    >
                        {{ $t('common.previous') }}
                    </UButton>

                    <UButton
                        v-for="page in beforePages"
                        :key="page"
                        class="relative inline-flex items-center bg-secondary-500 px-4 py-2 text-sm font-semibold hover:bg-base-400 focus:z-20 focus:outline-offset-0"
                        @click="$emit('offsetChange', calculateOffset(page))"
                    >
                        {{ page }}
                    </UButton>

                    <UButton
                        class="relative inline-flex items-center bg-primary-500 px-4 py-2 text-sm font-semibold underline hover:bg-primary-400"
                        disabled
                    >
                        {{ currentPage }}
                    </UButton>

                    <UButton
                        v-for="page in afterPages"
                        :key="page"
                        class="relative inline-flex items-center bg-secondary-500 px-4 py-2 text-sm font-semibold hover:bg-base-400 focus:z-20 focus:outline-offset-0"
                        @click="$emit('offsetChange', calculateOffset(page))"
                    >
                        {{ page }}
                    </UButton>

                    <Popover v-if="totalPages > 1n" class="relative">
                        <PopoverButton
                            class="relative inline-flex items-center bg-secondary-500 px-4 py-2 text-sm font-semibold hover:bg-base-400 focus:z-20 focus:outline-offset-0"
                            @click="pageNumber = ''"
                        >
                            ...
                        </PopoverButton>

                        <PopoverPanel
                            focus
                            class="absolute z-5 w-24 min-w-fit max-w-24 rounded-lg border border-gray-600 bg-gray-800 text-sm text-gray-400 shadow-sm transition-opacity"
                        >
                            <div class="p-3">
                                <form @submit.prevent="$emit('offsetChange', calculateOffset(parseInt(pageNumber)))">
                                    <UInput
                                        v-model="pageNumber"
                                        type="number"
                                        min="1"
                                        :max="parseInt(totalPages.toString())"
                                        class="remove-arrow block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        name="page_number"
                                        :placeholder="$t('common.page')"
                                    />
                                </form>
                            </div>
                        </PopoverPanel>
                    </Popover>

                    <UButton
                        v-if="currentPage <= totalPages - 1"
                        class="relative inline-flex items-center bg-secondary-500 px-4 py-2 text-sm font-semibold hover:bg-base-400 focus:z-20 focus:outline-offset-0"
                        @click="$emit('offsetChange', calculateOffset(totalPages))"
                    >
                        {{ totalPages }}
                    </UButton>

                    <UButton
                        :disabled="total - end <= 0n"
                        :class="[
                            total - end <= 0n
                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                : 'bg-primary-500 hover:bg-primary-400',
                            'relative ml-3 inline-flex cursor-pointer items-center rounded-r-md px-3 py-2 text-sm font-semibold focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2',
                        ]"
                        @click="$emit('offsetChange', end)"
                    >
                        {{ $t('common.next') }}
                    </UButton>
                </template>
                <template v-else>
                    <UButton
                        :disabled="offset <= 0n"
                        :class="[
                            offset <= 0n
                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                : 'bg-primary-500 hover:bg-primary-400',
                            'relative ml-3 inline-flex cursor-pointer items-center rounded-l-md px-3 py-2 text-sm font-semibold focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2',
                        ]"
                        @click="$emit('offsetChange', offset - pageSize < 0 ? 0 : offset - pageSize)"
                    >
                        {{ $t('common.previous') }}
                    </UButton>
                    <UButton
                        :disabled="end < pageSize"
                        :class="[
                            end < pageSize
                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                : 'bg-primary-500 hover:bg-primary-400',
                            'relative ml-3 inline-flex cursor-pointer items-center rounded-r-md px-3 py-2 text-sm font-semibold focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2',
                        ]"
                        @click="$emit('offsetChange', offset + pageSize)"
                    >
                        {{ $t('common.next') }}
                    </UButton>
                </template>
            </nav>
        </div>
    </nav>
</template>
