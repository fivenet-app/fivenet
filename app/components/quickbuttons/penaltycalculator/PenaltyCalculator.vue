<script lang="ts" setup>
import { UButton, UIcon, UTooltip } from '#components';
import type { TableColumn } from '@nuxt/ui';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import PenaltyStats from '~/components/quickbuttons/penaltycalculator/PenaltyStats.vue';
import PenaltySummaryTable from '~/components/quickbuttons/penaltycalculator/PenaltySummaryTable.vue';
import { useCompletorStore } from '~/stores/completor';
import type { Law } from '~~/gen/ts/resources/laws/laws';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { calculatePenaltySummary, type SelectedPenalty } from './helpers';

const { display, quickButtons } = useAppConfig();

const completorStore = useCompletorStore();
const notifications = useNotificationsStore();

const { t, d, n } = useI18n();

const { data: lawBooks, status, refresh, error } = useLazyAsyncData(`lawbooks`, () => completorStore.listLawBooks());

const formatter = new Intl.NumberFormat(display.intlLocale, {
    style: 'currency',
    currency: display.currencyName,
    trailingZeroDisplay: 'stripIfInteger',
});

const querySearchRaw = ref('');
const querySearch = computed(() => querySearchRaw.value.trim().toLowerCase());

const selectedPenalties = useState<SelectedPenalty[]>('quickButton:penaltyCalculator:selected', () => [] as SelectedPenalty[]);
const reduction = useState<number>('quickButton:penaltyCalculator:reduction', () => 0);

const summary = computed(() => calculatePenaltySummary(selectedPenalties.value));

const filteredLawBooks = computed(() =>
    lawBooks.value
        ?.map((book) => {
            const laws = book.laws
                .filter(
                    (p) =>
                        p.name.toLowerCase().includes(querySearch.value) ||
                        p.description?.toLowerCase().includes(querySearch.value),
                )
                .map((p) => {
                    const show = true;
                    return {
                        ...p,
                        show,
                    };
                });
            return {
                ...book,
                laws,
            };
        })
        .filter((books) => books.laws.length > 0)
        .map((book) => {
            return {
                label: `${book.name}` + (!book.description ? '' : ' - ' + book.description),
                book: book,
            };
        }),
);

const leeway = computed(() => reduction.value / 100);
const maxLeeway = computed(() => quickButtons.penaltyCalculator?.maxLeeway ?? 25);

function getNameForLawBookId(id: number): string | undefined {
    return lawBooks.value?.filter((b) => b.id === id)[0]?.name;
}

function copySummary(): void {
    let text =
        t('components.penaltycalculator.title') +
        ` (` +
        d(new Date(), 'long') +
        `)

${t('common.fine')}: ${n(summary.value.fine, 'currency')}${
            leeway.value > 0 && summary.value.fine > 0
                ? ` ${formatter.format(-Math.abs(summary.value.fine * leeway.value))}`
                : ''
        }
${t('common.detention_time')}: ${summary.value.detentionTime} ${t('common.month', summary.value.detentionTime)}${
            leeway.value > 0 && summary.value.detentionTime > 0
                ? ` (-${summary.value.detentionTime * leeway.value} ${t('common.month', summary.value.detentionTime * leeway.value)})`
                : ''
        }
${t('common.traffic_infraction_points', 2)}: ${summary.value.stvoPoints}${
            leeway.value > 0 && summary.value.stvoPoints > 0
                ? ` (-${summary.value.stvoPoints * leeway.value} ${t('common.points', summary.value.stvoPoints * leeway.value)})`
                : ''
        }
${t('common.reduction')}: ${reduction.value}%
${t('common.total_count')}: ${summary.value.count}
`;

    if (selectedPenalties.value.length > 0) {
        text += `
${t('common.crime', selectedPenalties.value.length)}:
`;

        selectedPenalties.value.forEach((v) => {
            text += `* ${getNameForLawBookId(v.law.lawbookId)} - ${v.law.name} (${v.count}x)
`;
        });
    }

    notifications.add({
        title: { key: 'notifications.penaltycalculator.title', parameters: {} },
        description: { key: 'notifications.penaltycalculator.content', parameters: {} },
        type: NotificationType.INFO,
    });

    copyToClipboardWrapper(text);
}

function updateLaw(selected: SelectedPenalty): void {
    const index = selectedPenalties.value.findIndex((p) => p.law.id === selected.law.id);
    if (index !== -1) {
        if (selected.count === 0) {
            selectedPenalties.value.splice(index, 1);
        } else {
            selectedPenalties.value[index] = selected;
        }
    } else {
        if (selected.count > 0) {
            selectedPenalties.value.push(selected);
        }
    }
}

function reset(): void {
    querySearchRaw.value = '';
    selectedPenalties.value = [];
}

const columns = computed(
    () =>
        [
            {
                accessorKey: 'name',
                header: t('common.law'),
                cell: ({ row }) =>
                    h('div', { class: 'inline-flex items-center gap-2' }, [
                        h('span', { class: 'whitespace-pre-line text-highlighted' }, row.original.name),
                        row.original.hint
                            ? h(UTooltip, { text: row.original.hint }, () =>
                                  h(UIcon, { class: 'size-5', name: 'i-mdi-information-outline' }),
                              )
                            : null,
                    ]),
            },
            {
                accessorKey: 'fine',
                header: t('common.fine'),
                cell: ({ row }) => formatter.format(row.original.fine ?? 0),
            },
            {
                accessorKey: 'detentionTime',
                header: t('common.detention_time'),
                cell: ({ row }) => row.original.detentionTime,
            },
            {
                accessorKey: 'stvoPoints',
                header: t('common.traffic_infraction_points', 2),
                cell: ({ row }) => row.original.stvoPoints,
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
                        row.original.description,
                    ),
            },
            {
                accessorKey: 'count',
                header: t('common.count'),
            },
        ] as TableColumn<Law>[],
);
</script>

<template>
    <div>
        <div class="pb-2 sm:flex sm:items-center">
            <div class="sm:flex-auto">
                <DataPendingBlock
                    v-if="isRequestPending(status)"
                    class="mt-5"
                    :message="$t('common.loading', [$t('common.law', 2)])"
                />
                <DataErrorBlock
                    v-else-if="error"
                    class="mt-5"
                    :title="$t('common.unable_to_load', [$t('common.law', 2)])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock
                    v-else-if="lawBooks && lawBooks.length === 0"
                    class="mt-5"
                    icon="i-mdi-gavel"
                    :type="$t('common.law', 2)"
                />

                <div v-else>
                    <UFormField name="search">
                        <UInput
                            v-model="querySearchRaw"
                            type="text"
                            name="search"
                            class="w-full"
                            :placeholder="$t('common.filter')"
                            :ui="{ trailing: 'pe-1' }"
                        >
                            <template #trailing>
                                <UButton
                                    v-if="querySearchRaw !== ''"
                                    color="red"
                                    variant="link"
                                    icon="i-mdi-clear"
                                    aria-controls="search"
                                    @click="querySearchRaw = ''"
                                />
                            </template>
                        </UInput>
                    </UFormField>

                    <dl class="mt-2">
                        <UAccordion type="multiple" :items="filteredLawBooks">
                            <template #content="{ item: lawBook }">
                                <UTable
                                    :columns="columns"
                                    :data="lawBook.book.laws"
                                    :empty="$t('common.not_found', [$t('common.law', 2)])"
                                    :pagination-options="{ manualPagination: true }"
                                    :sorting-options="{ manualSorting: true }"
                                >
                                    <template #count-cell="{ row }">
                                        <UInputNumber
                                            :model-value="
                                                selectedPenalties.find((p) => p.law.id === row.original.id)?.count ?? 0
                                            "
                                            name="count"
                                            :min="0"
                                            :max="quickButtons.penaltyCalculator?.maxCount ?? 10"
                                            :step="1"
                                            class="max-w-22 min-w-20 grow-0"
                                            @update:model-value="
                                                ($event) =>
                                                    updateLaw({ law: row.original, count: $event === null ? 0 : $event })
                                            "
                                        />
                                    </template>
                                </UTable>
                            </template>
                        </UAccordion>
                    </dl>
                </div>
            </div>
        </div>

        <USeparator :label="$t('common.result')" class="mb-2" />

        <PenaltyStats :summary="summary" :reduction="reduction" />

        <div class="my-2 flex flex-row items-center gap-2 text-sm">
            <p class="font-semibold">
                {{ $t('common.reduction') }}
            </p>
            <USlider v-model="reduction" size="sm" :min="0" :max="maxLeeway" :step="1" />
            <p class="w-12">{{ reduction }}%</p>
        </div>

        <div class="overflow-x-auto">
            <PenaltySummaryTable
                v-if="lawBooks && lawBooks.length > 0"
                :law-books="lawBooks"
                :selected-laws="selectedPenalties"
                :reduction="reduction"
            />
        </div>

        <UFieldGroup class="mt-2 inline-flex w-full">
            <UButton class="flex-1" icon="i-mdi-content-copy" :label="$t('common.copy')" @click="() => copySummary()" />
            <UButton trailing-icon="i-mdi-clear-outline" color="error" :label="$t('common.reset')" @click="reset()" />
        </UFieldGroup>
    </div>
</template>
