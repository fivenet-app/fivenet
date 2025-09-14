<script lang="ts" setup>
import { UButton, UIcon, UTooltip } from '#components';
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

function getNameForLawBookId(id: number): string | undefined {
    return props.lawBooks?.filter((b) => b.id === id)[0]?.name;
}

const leeway = computed(() => props.reduction / 100);

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
                cell: ({ row }) =>
                    h('span', null, [
                        `$${row.original.law.fine ? row.original.law.fine * row.original.count : 0}`,
                        (row.original.law.fine ?? 0) * row.original.count > 0 && leeway.value > 0
                            ? h(
                                  'span',
                                  null,
                                  ` ($-${((row.original.law.fine ?? 0) * row.original.count * leeway.value).toFixed(0)})`,
                              )
                            : null,
                    ]),
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
    <UButton v-if="selectedLaws.length === 0" class="relative block w-full p-4 text-center" disabled variant="outline">
        <UIcon class="mx-auto size-8" name="i-mdi-calculator" />
        <span class="mt-2 block text-sm font-semibold">
            {{ $t('common.none_selected', [`${$t('common.crime')}`]) }}
        </span>
    </UButton>

    <UTable
        v-else
        class="max-w-full divide-y divide-default"
        :columns="columns"
        :data="selectedLaws"
        :empty="$t('common.none_selected', [`${$t('common.crime')}`])"
    />
</template>
