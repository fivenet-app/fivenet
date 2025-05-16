<!-- eslint-disable vue/multi-word-component-names -->
<script lang="ts" setup>
import type { PaginationResponse } from '~~/gen/ts/resources/common/database/database';

const props = withDefaults(
    defineProps<{
        modelValue?: number;
        pagination?: PaginationResponse | undefined | null;
        disableBorder?: boolean;
        refresh?: () => Promise<void>;
        loading?: boolean;
        hideText?: boolean;
        hideButtons?: boolean;
    }>(),
    {
        modelValue: 0,
        pagination: undefined,
        disableBorder: false,
        refresh: undefined,
        loading: false,
        hideText: false,
        hideButtons: false,
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
    () => props.loading,
    () => {
        if (props.loading) {
            loadingState.value = true;
        }
    },
);
watchDebounced(
    () => props.loading,
    () => {
        if (!props.loading) {
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
    if (!canGoFirstOrPrev.value) {
        return;
    }

    currentPage.value--;
}

function onClickNext() {
    if (!canGoLastOrNext.value) {
        return;
    }

    currentPage.value++;
}
</script>

<template>
    <div class="@container">
        <div
            class="@md:flex-row flex justify-between gap-1 px-3 py-3 md:items-center"
            :class="!disableBorder ? 'border-t border-gray-200 dark:border-gray-700' : ''"
        >
            <div v-if="!hideText" class="flex flex-col items-center gap-2">
                <I18nT
                    v-if="!isInfinite"
                    class="@md:block hidden truncate text-sm"
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
                    class="@md:block hidden truncate text-sm"
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
                    :disabled="loading || loadingState"
                    :loading="loading || loadingState"
                    @click="refresh()"
                >
                    <span class="@md:block hidden">
                        {{ $t('common.refresh') }}
                    </span>
                </UButton>
            </UTooltip>

            <template v-if="!hideButtons">
                <UPagination
                    v-if="!isInfinite"
                    v-model="currentPage"
                    :page-count="pagination?.pageSize ?? 0"
                    :total="pagination?.totalCount ?? 0"
                />
                <UButtonGroup v-else>
                    <UButton
                        :label="$t('common.previous')"
                        color="gray"
                        icon="i-mdi-chevron-left"
                        :disabled="!canGoFirstOrPrev || loading"
                        @click="onClickPrev"
                    />
                    <UButton :label="currentPage.toString()" color="white" disabled />
                    <UButton
                        :label="$t('common.next')"
                        color="gray"
                        :disabled="!canGoLastOrNext || loading"
                        trailing
                        trailing-icon="i-mdi-chevron-right"
                        @click="onClickNext"
                    />
                </UButtonGroup>
            </template>
            <div v-else></div>

            <slot />
        </div>
    </div>
</template>
