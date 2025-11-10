<script lang="ts" setup>
import { UIcon, UTooltip } from '#components';
import type { TableColumn } from '@nuxt/ui';
import { h } from 'vue';
import { useI18n } from 'vue-i18n';
import type { SelectedPenalty } from '~/components/quickbuttons/penaltycalculator/PenaltyCalculator.vue';
import type { LawBook } from '~~/gen/ts/resources/laws/laws';

const props = defineProps<{
    lawBooks: LawBook[];
    selectedLaws: SelectedPenalty[];
    reduction: number;
}>();

const { t } = useI18n();

const { display } = useAppConfig();

function getNameForLawBookId(id: number): string | undefined {
    return props.lawBooks?.filter((b) => b.id === id)[0]?.name;
}

const leeway = computed(() => props.reduction / 100);

const formatter = new Intl.NumberFormat(display.intlLocale, {
    style: 'currency',
    currency: display.currencyName,
});

const columns = computed(
    () =>
        [
            {
                accessorKey: 'law',
                header: t('common.law'),
                cell: ({ row }) =>
                    h('div', { class: 'inline-flex items-center gap-2' }, [
                        h(
                            'p',
                            { class: 'whitespace-pre-line text-highlighted' },
                            `${getNameForLawBookId(row.original.law.lawbookId)} - ${row.original.law.name}`,
                        ),
                        row.original.law.hint
                            ? h(UTooltip, { text: row.original.law.hint }, () =>
                                  h(UIcon, { class: 'size-5', name: 'i-mdi-information-outline' }),
                              )
                            : null,
                    ]),
            },
            {
                accessorKey: 'fine',
                header: t('common.fine'),
                cell: ({ row }) => {
                    return h('span', null, [
                        formatter.format(row.original.law.fine ? row.original.law.fine * row.original.count : 0),
                        (row.original.law.fine ?? 0) * row.original.count > 0 && leeway.value > 0
                            ? h(
                                  'span',
                                  null,
                                  ` (${formatter.format(-Math.abs(row.original.law.fine ?? 0) * row.original.count * leeway.value)})`,
                              )
                            : null,
                    ]);
                },
            },
            {
                accessorKey: 'detentionTime',
                header: t('common.detention_time'),
                cell: ({ row }) =>
                    h('span', null, [
                        `${row.original.law.detentionTime ? row.original.law.detentionTime * row.original.count : 0}`,
                        (row.original.law.detentionTime ?? 0) * row.original.count > 0 && leeway.value > 0
                            ? h(
                                  'span',
                                  null,
                                  ` (-${((row.original.law.detentionTime ?? 0) * row.original.count * leeway.value).toFixed(0)})`,
                              )
                            : null,
                    ]),
            },
            {
                accessorKey: 'trafficInfractionPoints',
                header: t('common.traffic_infraction_points', 2),
                cell: ({ row }) =>
                    h('span', null, [
                        `${row.original.law.stvoPoints ? row.original.law.stvoPoints * row.original.count : 0}`,
                        (row.original.law.stvoPoints ?? 0) * row.original.count > 0 && leeway.value > 0
                            ? h(
                                  'span',
                                  null,
                                  ` (-${((row.original.law.stvoPoints ?? 0) * row.original.count * leeway.value).toFixed(0)})`,
                              )
                            : null,
                    ]),
            },
            {
                accessorKey: 'description',
                header: t('common.description'),
                cell: ({ row }) =>
                    h(
                        'p',
                        {
                            class: 'line-clamp-2 w-full max-w-sm break-all whitespace-normal hover:line-clamp-none',
                        },
                        row.original.law.description,
                    ),
            },
            {
                accessorKey: 'count',
                header: t('common.count'),
                cell: ({ row }) => row.original.count,
            },
        ] as TableColumn<SelectedPenalty>[],
);
</script>

<template>
    <UAlert
        v-if="selectedLaws.length === 0"
        class="h-[64px] max-h-[64px] w-full items-center p-2"
        :title="$t('common.none_selected', [$t('common.crime')])"
        icon="i-mdi-calculator"
        variant="outline"
    />

    <UTable
        v-else
        class="max-w-full divide-y divide-default"
        :columns="columns"
        :data="selectedLaws"
        :empty="$t('common.none_selected', [`${$t('common.crime')}`])"
    />
</template>
