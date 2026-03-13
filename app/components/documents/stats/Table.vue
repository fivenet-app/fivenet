<script setup lang="ts">
import { UBadge, UButton, ULink } from '#components';
import type { TableColumn } from '@nuxt/ui';
import { getGroupedRowModel } from '@tanstack/vue-table';
import { type CategoryValue, type GetStatsResponse, StatsCategory, type StatsPeriod } from '~~/gen/ts/services/documents/stats';
import type { Range } from './helpers';

const props = defineProps<{
    range: Range;
    period: StatsPeriod;
    category: StatsCategory;
    stats?: GetStatsResponse;
}>();

const { t } = useI18n();

const { format: formatNumber } = useIntlNumberFormatWithOptions({
    style: 'decimal',
    currency: undefined,
});

const categoryColumns = computed<TableColumn<CategoryValue>[]>(() => [
    {
        header: t('common.name'),
        cell: ({ row }) =>
            row.original.id
                ? h(ULink, { to: `/documents?categories=[${row.original.id}]` }, [row.original.name])
                : row.original.name,
    },
    {
        header: t('common.count'),
        cell: ({ row }) => formatNumber(row.original.value),
    },
]);

const categoryData = computed(() => props.stats?.documentsByCategory ?? []);

type TopLawRow = {
    lawBook: string;
    law: string;
    value: number;
};

const splitTopLawKey = (key: string): { lawBook: string; law: string } => {
    const separator = '::';
    const idx = key.indexOf(separator);
    if (idx === -1) {
        return {
            lawBook: t('common.na'),
            law: key,
        };
    }

    return {
        lawBook: key.slice(0, idx),
        law: key.slice(idx + separator.length),
    };
};

const topLawColumns = computed<TableColumn<TopLawRow>[]>(() => [
    {
        accessorKey: 'lawBook',
        header: t('common.law_book', 1),
    },
    {
        accessorKey: 'law',
        header: t('common.law', 1),
    },
    {
        accessorKey: 'value',
        header: t('common.count'),
        aggregateFn: 'sum',
        aggregatedCell: ({ getValue }) => formatNumber(getValue<number>()),
        cell: ({ row, getValue }) => {
            if (row.getIsGrouped()) {
                return h(UBadge, { color: 'neutral', variant: 'subtle' }, () => formatNumber(getValue<number>() ?? 0));
            }

            return formatNumber(getValue<number>());
        },
    },
]);

const topLawData = computed<TopLawRow[]>(() =>
    (props.stats?.topLaws ?? []).map((item) => {
        const split = splitTopLawKey(item.key);
        return {
            lawBook: split.lawBook,
            law: split.law,
            value: item.value,
        };
    }),
);
</script>

<template>
    <UCard :ui="{ root: 'overflow-visible', body: '!px-0 !pt-0 !pb-0' }">
        <UTable v-if="props.category === StatsCategory.DOCUMENTS_BY_CATEGORY" :columns="categoryColumns" :data="categoryData" />
        <UTable
            v-else-if="props.category === StatsCategory.TOP_LAWS"
            :columns="topLawColumns"
            :data="topLawData"
            :grouping="['lawBook']"
            :grouping-options="{
                getGroupedRowModel: getGroupedRowModel(),
            }"
            :sorting="[
                {
                    id: 'value',
                    desc: true,
                },
            ]"
            :ui="{ root: 'min-w-full', td: 'empty:p-0' }"
        >
            <template #lawBook-cell="{ row }">
                <div v-if="row.getIsGrouped()" class="flex items-center gap-2">
                    <span class="inline-block" :style="{ width: `calc(${row.depth} * 1rem)` }" />

                    <UButton
                        variant="outline"
                        color="neutral"
                        class="mr-1"
                        size="xs"
                        :icon="row.getIsExpanded() ? 'i-mdi-minus' : 'i-mdi-plus'"
                        @click="row.toggleExpanded()"
                    />

                    <strong>{{ row.original.lawBook }}</strong>
                </div>
                <span v-else />
            </template>
        </UTable>
    </UCard>
</template>
