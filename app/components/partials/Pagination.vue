<!-- eslint-disable vue/multi-word-component-names -->
<script lang="ts" setup>
import type { PaginationResponse } from '~~/gen/ts/resources/common/database/database';

const props = withDefaults(
    defineProps<{
        modelValue?: number;
        pagination?: PaginationResponse | undefined | null;
        infinite?: boolean;
        disableBorder?: boolean;
        refresh?: () => Promise<void>;
        loading?: boolean;
        hideText?: boolean;
        hideButtons?: boolean;
    }>(),
    {
        modelValue: 0,
        pagination: undefined,
        infinite: false,
        disableBorder: false,
        refresh: undefined,
        loading: false,
        hideText: false,
        hideButtons: false,
    },
);

const emit = defineEmits<{
    (e: 'update:modelValue', offset: number): void;
}>();

const page = useVModel(props, 'modelValue', emit);

const total = computed(() => props.pagination?.totalCount ?? 0);
const pageSize = computed(() => props.pagination?.pageSize ?? 1);

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
</script>

<template>
    <div class="@container">
        <div
            class="@md:flex-row flex justify-between gap-1 px-3 py-3 md:items-center"
            :class="!disableBorder ? 'border-t border-gray-200 dark:border-gray-700' : ''"
        >
            <div v-if="!hideText" class="flex flex-col items-center gap-2">
                <I18nT
                    keypath="components.partials.table_pagination.page_count"
                    tag="p"
                    class="@md:block hidden truncate text-sm"
                >
                    <template #current>
                        <span class="text-neutral font-medium">
                            {{ page }}
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
            </div>
            <div v-else></div>

            <UTooltip v-if="refresh" :text="$t('common.refresh')">
                <UButton
                    variant="link"
                    icon="i-mdi-refresh"
                    :disabled="loading || loadingState"
                    :loading="loading || loadingState"
                    class="p-px"
                    @click="refresh()"
                >
                    <span class="@md:block hidden">
                        {{ $t('common.refresh') }}
                    </span>
                </UButton>
            </UTooltip>

            <UPagination
                v-if="!hideButtons"
                v-model="page"
                :page-count="pagination?.pageSize ?? 0"
                :total="!infinite ? (pagination?.totalCount ?? 0) : (page + 1) * (pagination?.pageSize ?? 0)"
            />
            <div v-else></div>

            <slot />
        </div>
    </div>
</template>
