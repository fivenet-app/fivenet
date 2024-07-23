<!-- eslint-disable vue/multi-word-component-names -->
<script lang="ts" setup>
import type { PaginationResponse } from '~~/gen/ts/resources/common/database/database';

const props = withDefaults(
    defineProps<{
        modelValue: number;
        pagination: PaginationResponse | undefined | null;
        infinite?: boolean;
        disableBorder?: boolean;
        refresh?: () => Promise<void>;
        loading?: boolean;
    }>(),
    {
        infinite: false,
        disableBorder: false,
        loading: false,
    },
);

const emits = defineEmits<{
    (e: 'update:modelValue', offset: number): void;
}>();

const page = useVModel(props, 'modelValue', emits);

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
            class="@sm:flex-row flex flex-col-reverse justify-between gap-1 px-3 py-3.5 sm:items-center"
            :class="!disableBorder ? 'border-t border-gray-200 dark:border-gray-700' : ''"
        >
            <div class="@sm:flex-row inline-flex flex-col items-center gap-2">
                <I18nT keypath="components.partials.table_pagination.page_count" tag="p" class="text-sm">
                    <template #current>
                        <span class="text-neutral font-medium">
                            {{ page }}
                        </span>
                    </template>
                    <template #total>
                        <span class="text-neutral font-medium">
                            {{ total }}
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

                <UButton
                    v-if="refresh"
                    variant="link"
                    icon="i-mdi-refresh"
                    :title="$t('common.refresh')"
                    :disabled="loading || loadingState"
                    :loading="loading || loadingState"
                    class="p-px"
                    @click="refresh()"
                >
                    <span class="hidden sm:block">
                        {{ $t('common.refresh') }}
                    </span>
                </UButton>
            </div>

            <UPagination
                v-model="page"
                :page-count="pagination?.pageSize ?? 0"
                :total="!infinite ? pagination?.totalCount ?? 0 : (page + 1) * (pagination?.pageSize ?? 0)"
            />
        </div>
    </div>
</template>
