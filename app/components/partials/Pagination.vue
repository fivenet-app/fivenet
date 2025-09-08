<!-- eslint-disable vue/multi-word-component-names -->
<script lang="ts" setup>
import type { AsyncDataRequestStatus } from '#app';
import type { PaginationResponse } from '~~/gen/ts/resources/common/database/database';

const props = withDefaults(
    defineProps<{
        modelValue?: number;
        pagination?: PaginationResponse | undefined | null;
        disableBorder?: boolean;
        refresh?: () => Promise<void>;
        status?: AsyncDataRequestStatus;
        hideText?: boolean;
        hideButtons?: boolean;
        compact?: boolean;
    }>(),
    {
        modelValue: 0,
        pagination: undefined,
        disableBorder: false,
        refresh: undefined,
        status: 'pending',
        hideText: false,
        hideButtons: false,
        compact: false,
    },
);

const emit = defineEmits<{
    (e: 'update:modelValue', page: number): void;
}>();

const currentPage = useVModel(props, 'modelValue', emit);

const total = computed(() => props.pagination?.totalCount ?? 0);
const pageSize = computed(() => props.pagination?.pageSize ?? 10);

const totalPages = computed(() => Math.ceil(total.value / pageSize.value));

const loadingState = ref(false);
watch(
    () => props.status,
    () => {
        if (isRequestPending(props.status)) {
            loadingState.value = true;
        }
    },
);
watchDebounced(
    () => props.status,
    () => {
        if (!isRequestPending(props.status)) {
            loadingState.value = false;
        }
    },
    {
        debounce: 750,
        maxWait: 1250,
    },
);

const isInfinite = computed(() => props.pagination && props.pagination.totalCount === -1);
const canGoFirstOrPrev = computed(() => currentPage.value > 1);
const canGoLastOrNext = computed(
    () =>
        (props.pagination &&
            props.pagination.totalCount === -1 &&
            props.pagination?.offset + props.pagination?.pageSize <= props.pagination?.end) ||
        currentPage.value < totalPages.value,
);

function onClickPrev() {
    if (!canGoFirstOrPrev.value) return;

    currentPage.value--;
}

function onClickNext() {
    if (!canGoLastOrNext.value) return;

    currentPage.value++;
}
</script>

<template>
    <div class="@container/pagination">
        <div
            class="flex justify-between gap-1 md:items-center @md/pagination:flex-row"
            :class="[
                !disableBorder ? 'border-t border-neutral-200 dark:border-neutral-700' : '',
                compact ? 'px-1 py-1' : 'px-3 py-3',
            ]"
        >
            <div v-if="!hideText" class="flex flex-col items-center gap-2">
                <I18nT
                    v-if="!isInfinite"
                    class="hidden truncate text-sm @md/pagination:block"
                    keypath="components.partials.table_pagination.page_count_with_total"
                    tag="p"
                >
                    <template #current>
                        <span class="text-neutral font-medium">
                            {{ currentPage }}
                        </span>
                    </template>

                    <template #total>
                        <span class="text-neutral font-medium">
                            {{ total === -1 ? '∞' : total }}
                        </span>
                    </template>

                    <template #maxPage>
                        <span class="text-neutral font-medium">
                            {{ totalPages === 0 ? 1 : totalPages }}
                        </span>
                    </template>

                    <template #size>
                        <span class="text-neutral font-medium">
                            {{ pageSize }}
                        </span>
                    </template>
                </I18nT>

                <I18nT
                    v-else
                    class="hidden truncate text-sm @md/pagination:block"
                    keypath="components.partials.table_pagination.page_count"
                    tag="p"
                >
                    <template #current>
                        <span class="text-neutral font-medium">
                            {{ currentPage }}
                        </span>
                    </template>

                    <template #size>
                        <span class="text-neutral font-medium">
                            {{ pageSize }}
                        </span>
                    </template>
                </I18nT>
            </div>
            <div v-else></div>

            <UTooltip v-if="refresh" :text="$t('common.refresh')">
                <UButton
                    class="p-px"
                    variant="link"
                    icon="i-mdi-refresh"
                    :disabled="loadingState || isRequestPending(status)"
                    :loading="loadingState || isRequestPending(status)"
                    @click="refresh()"
                >
                    <span class="hidden @md/pagination:block">
                        {{ $t('common.refresh') }}
                    </span>
                </UButton>
            </UTooltip>

            <template v-if="!hideButtons">
                <UPagination
                    v-if="!isInfinite"
                    v-model:page="currentPage"
                    :page-count="pagination?.pageSize ?? 0"
                    :total="pagination?.totalCount ?? 0"
                    :show-edges="false"
                    :ui="{ first: 'hidden', last: 'hidden' }"
                />
                <UButtonGroup v-else>
                    <UButton
                        color="neutral"
                        variant="outline"
                        icon="i-mdi-chevron-left"
                        :disabled="!canGoFirstOrPrev || isRequestPending(status)"
                        @click="onClickPrev"
                    />
                    <UButton :label="currentPage.toString()" color="primary" variant="solid" />
                    <UButton
                        color="neutral"
                        variant="outline"
                        :disabled="!canGoLastOrNext || isRequestPending(status)"
                        icon="i-mdi-chevron-right"
                        @click="onClickNext"
                    />
                </UButtonGroup>
            </template>
            <div v-else></div>

            <slot />
        </div>
    </div>
</template>
