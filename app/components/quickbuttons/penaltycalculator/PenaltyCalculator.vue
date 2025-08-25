<script lang="ts" setup>
import { UButton, UIcon, USelect, UTooltip } from '#components';
import type { TableColumn } from '@nuxt/ui';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import PenaltyStats from '~/components/quickbuttons/penaltycalculator/PenaltyStats.vue';
import PenaltySummaryTable from '~/components/quickbuttons/penaltycalculator/PenaltySummaryTable.vue';
import { useCompletorStore } from '~/stores/completor';
import type { Law } from '~~/gen/ts/resources/laws/laws';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const completorStore = useCompletorStore();
const notifications = useNotificationsStore();

const { t, d, n } = useI18n();

const { data: lawBooks, status, refresh, error } = useLazyAsyncData(`lawbooks`, () => completorStore.listLawBooks());

export type SelectedPenalty = {
    law: Law;
    count: number;
};

export type PenaltiesSummary = {
    selectedPenalties: SelectedPenalty[];
    fine: number;
    detentionTime: number;
    stvoPoints: number;
    count: number;
};

const querySearchRaw = ref('');
const querySearch = computed(() => querySearchRaw.value.trim().toLowerCase());

const state = useState<PenaltiesSummary>('quickButton:penaltyCalculator:summary', () => ({
    selectedPenalties: [],
    fine: 0,
    detentionTime: 0,
    stvoPoints: 0,
    count: 0,
}));

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

const reduction = ref<number>(0);
const leeway = computed(() => reduction.value / 100);

function getNameForLawBookId(id: number): string | undefined {
    return lawBooks.value?.filter((b) => b.id === id)[0]?.name;
}

function calculate(e: SelectedPenalty): void {
    const idx = state.value.selectedPenalties.findIndex(
        (v) => v.law.lawbookId === e.law.lawbookId && v.law.name === e.law.name,
    );

    let count = e.count;
    if (idx > -1) {
        const existing = state.value.selectedPenalties.at(idx)!;
        state.value.selectedPenalties[idx] = e;

        if (existing.count !== e.count) {
            count = e.count - existing.count;
        }

        // If the selected penalty count is 0, remove it from the list
        if (e.count === 0) {
            state.value.selectedPenalties.splice(idx, 1);
        }
    } else if (e.count !== 0) {
        state.value.selectedPenalties.push(e);
    }

    if (e.law.fine) {
        state.value.fine += count * e.law.fine;
    }
    if (e.law.detentionTime) {
        state.value.detentionTime += count * e.law.detentionTime;
    }
    if (e.law.stvoPoints) {
        state.value.stvoPoints += count * e.law.stvoPoints;
    }
    state.value.count = state.value.count + count;
}

async function copySummary(): Promise<void> {
    let text =
        t('components.penaltycalculator.title') +
        ` (` +
        d(new Date(), 'long') +
        `)

${t('common.fine')}: ${n(state.value.fine, 'currency')}${
            leeway.value > 0 && state.value.fine > 0 ? ` ($-${(state.value.fine * leeway.value).toFixed(0)})` : ''
        }
${t('common.detention_time')}: ${state.value.detentionTime} ${t('common.time_ago.month', state.value.detentionTime)}${
            leeway.value > 0 && state.value.detentionTime > 0
                ? ` (-${(state.value.detentionTime * leeway.value).toFixed(0)} ${t('common.time_ago.month', (state.value.detentionTime * leeway.value).toFixed(0))})`
                : ''
        }
${t('common.traffic_infraction_points', 2)}: ${state.value.stvoPoints}${
            leeway.value > 0 && state.value.stvoPoints > 0
                ? ` (-${(state.value.stvoPoints * leeway.value).toFixed(0)} ${t('common.points', (state.value.stvoPoints * leeway.value).toFixed(0))})`
                : ''
        }
${t('common.reduction')}: ${reduction.value}%
${t('common.total_count')}: ${state.value.count}
`;

    if (state.value.selectedPenalties.length > 0) {
        text += `
${t('common.crime', state.value.selectedPenalties.length)}:
`;

        state.value.selectedPenalties.forEach((v) => {
            text += `* ${getNameForLawBookId(v.law.lawbookId)} - ${v.law.name} (${v.count}x)
`;
        });
    }

    notifications.add({
        title: { key: 'notifications.penaltycalculator.title', parameters: {} },
        description: { key: 'notifications.penaltycalculator.content', parameters: {} },
        type: NotificationType.INFO,
    });

    return copyToClipboardWrapper(text);
}

function reset(): void {
    querySearchRaw.value = '';
    state.value.selectedPenalties = [];

    state.value.count = 0;
    state.value.detentionTime = 0;
    state.value.fine = 0;
    state.value.stvoPoints = 0;
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
                cell: ({ row }) => $n(row.original.fine ?? 0, 'currency'),
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
                cell: ({ row }) =>
                    h(USelect, {
                        modelValue: state.value.selectedPenalties.find((p) => p.law.id === row.original.id)?.count ?? 0,
                        name: 'count',
                        items: Array.from(Array(7).keys()).map((v) => ({ value: v, label: v.toString() })),
                        'update:modelValue': ($event) => calculate({ law: row.original, count: parseInt($event) }),
                    }),
            },
        ] as TableColumn<Law>[],
);
</script>

<template>
    <div class="py-2">
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
                            :placeholder="$t('common.filter')"
                            :ui="{ trailing: 'pe-1' }"
                        >
                            <template #trailing>
                                <UButton
                                    v-show="querySearchRaw !== ''"
                                    color="neutral"
                                    variant="link"
                                    icon="i-mdi-close"
                                    aria-controls="search"
                                    @click="querySearchRaw = ''"
                                />
                            </template>
                        </UInput>
                    </UFormField>

                    <dl class="mt-4">
                        <UAccordion multiple :items="filteredLawBooks">
                            <template #content="{ item: lawBook }">
                                <UTable
                                    :columns="columns"
                                    :data="lawBook.book.laws"
                                    :empty="$t('common.not_found', [$t('common.law', 2)])"
                                    :pagination-options="{ manualPagination: true }"
                                    :sorting-options="{ manualSorting: true }"
                                />
                            </template>
                        </UAccordion>
                    </dl>
                </div>
            </div>
        </div>

        <USeparator :label="$t('common.result')" />

        <div class="flow-root">
            <div class="overflow-x-auto sm:-mx-6 lg:-mx-8">
                <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                    <div class="text-xl">
                        <PenaltyStats :summary="state" :reduction="reduction" />

                        <div class="my-2 flex flex-row items-center gap-2 text-sm">
                            <p class="font-semibold">
                                {{ $t('common.reduction') }}
                            </p>
                            <USlider v-model="reduction" size="sm" :min="0" :max="25" :step="1" />
                            <p class="w-12">{{ reduction }}%</p>
                        </div>

                        <div>
                            <PenaltySummaryTable
                                v-if="lawBooks && lawBooks.length > 0"
                                :law-books="lawBooks"
                                :selected-laws="state.selectedPenalties"
                                :reduction="reduction"
                            />
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <UButtonGroup class="mt-2 inline-flex w-full">
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" icon="i-mdi-content-copy" @click="copySummary()">
                    {{ $t('common.copy') }}
                </UButton>
                <UButton trailing-icon="i-mdi-clear-outline" color="error" @click="reset()">
                    {{ $t('common.reset') }}
                </UButton>
            </UButtonGroup>
        </UButtonGroup>
    </div>
</template>
