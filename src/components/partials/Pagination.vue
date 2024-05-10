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
</script>

<template>
    <div
        class="flex justify-between px-3 py-3.5"
        :class="!disableBorder ? 'border-t border-gray-200 dark:border-gray-700' : ''"
    >
        <div class="inline-flex items-center gap-2">
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
                trailing-icon="i-mdi-refresh"
                :title="$t('common.refresh')"
                :disabled="loading"
                :loading="loading"
                @click="refresh()"
            >
                {{ $t('common.refresh') }}
            </UButton>
        </div>

        <UPagination
            v-model="page"
            :page-count="pagination?.pageSize ?? 0"
            :total="!infinite ? pagination?.totalCount ?? 0 : (page + 1) * (pagination?.pageSize ?? 0)"
        />
    </div>
</template>
