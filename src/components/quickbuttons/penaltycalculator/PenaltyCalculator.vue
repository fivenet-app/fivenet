<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import PenaltyStats from '~/components/quickbuttons/penaltycalculator/PenaltyStats.vue';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import { Law } from '~~/gen/ts/resources/laws/laws';
import PenaltySummaryTable from '~/components/quickbuttons/penaltycalculator/PenaltySummaryTable.vue';

const completorStore = useCompletorStore();
const notifications = useNotificatorStore();

const { t, d } = useI18n();

const props = withDefaults(
    defineProps<{
        loadLaws: boolean;
    }>(),
    {
        loadLaws: true,
    },
);

const { data: lawBooks, pending, refresh, error } = useLazyAsyncData(`lawbooks`, () => completorStore.listLawBooks());

export type SelectedPenalty = {
    law: Law;
    count: number;
};

export type PenaltiesSummary = {
    fine: number;
    detentionTime: number;
    stvoPoints: number;
    count: number;
};

const querySearchRaw = ref('');
const querySearch = computed(() => querySearchRaw.value.trim().toLowerCase());

const selectedPenalties = ref<SelectedPenalty[]>([]);

const summary = ref<PenaltiesSummary>({
    fine: 0,
    detentionTime: 0,
    stvoPoints: 0,
    count: 0,
});

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

function getNameForLawBookId(id: string): string | undefined {
    return lawBooks.value?.filter((b) => b.id === id)[0].name;
}

function calculate(e: SelectedPenalty): void {
    const idx = selectedPenalties.value.findIndex((v) => v.law.lawbookId === e.law.lawbookId && v.law.name === e.law.name);

    let count = e.count;
    if (idx > -1) {
        const existing = selectedPenalties.value.at(idx)!;
        selectedPenalties.value[idx] = e;
        if (existing.count !== e.count) {
            count = e.count - existing.count;
        }
        // If the selected penalty count is 0, remove it from the list
        if (e.count === 0) {
            selectedPenalties.value.splice(idx, 1);
        }
    } else if (e.count !== 0) {
        selectedPenalties.value.push(e);
    }

    if (e.law.fine) {
        summary.value.fine += count * e.law.fine;
    }
    if (e.law.detentionTime) {
        summary.value.detentionTime += count * e.law.detentionTime;
    }
    if (e.law.stvoPoints) {
        summary.value.stvoPoints += count * e.law.stvoPoints;
    }
    summary.value.count = summary.value.count + count;
}

async function copyToClipboard(): Promise<void> {
    let text =
        t('components.penaltycalculator.title') +
        ` (` +
        d(new Date(), 'long') +
        `)

${t('common.fine')}: $${summary.value.fine}
${t('common.detention_time')}: ${summary.value.detentionTime} ${t(
            'common.time_ago.month',
            summary.value.detentionTime.toString(),
        )}
${t('common.traffic_infraction_points', 2)}: ${summary.value.stvoPoints}
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
        type: 'info',
    });

    return copyToClipboardWrapper(text);
}

function reset(): void {
    querySearchRaw.value = '';
    selectedPenalties.value = [];

    summary.value.count = 0;
    summary.value.detentionTime = 0;
    summary.value.fine = 0;
    summary.value.stvoPoints = 0;
}

watch(props, () => {
    if (lawBooks.value === null && pending.value) {
        refresh();
    }
});

const columns = [
    {
        key: 'name',
        label: t('common.law'),
    },
    {
        key: 'fine',
        label: t('common.fine'),
    },
    {
        key: 'detentionTime',
        label: t('common.detention_time'),
    },
    {
        key: 'stvoPoints',
        label: t('common.traffic_infraction_points', 2),
    },
    {
        key: 'description',
        label: t('common.description'),
    },
    {
        key: 'count',
        label: t('common.count'),
    },
];
</script>

<template>
    <div class="py-2">
        <div class="pb-2 sm:flex sm:items-center">
            <div class="sm:flex-auto">
                <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.law', 2)])" class="mt-5" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.not_found', [$t('common.law', 2)])"
                    :retry="refresh"
                    class="mt-5"
                />
                <DataNoDataBlock
                    v-else-if="lawBooks && lawBooks.length === 0"
                    icon="i-mdi-gavel"
                    :type="`${$t('common.citizen', 1)} ${$t('common.activity')}`"
                    class="mt-5"
                />

                <div v-else>
                    <UFormGroup name="search">
                        <UInput
                            v-model="querySearchRaw"
                            type="text"
                            name="search"
                            :placeholder="$t('common.filter')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                    </UFormGroup>

                    <dl class="mt-4">
                        <UAccordion multiple :items="filteredLawBooks">
                            <template #item="{ item: lawBook }">
                                <div class="max-w-full">
                                    <UTable :columns="columns" :rows="lawBook.book.laws">
                                        <template #name-data="{ row: law }">
                                            <p class="whitespace-pre-line text-gray-900 dark:text-gray-300">
                                                {{ law.name }}
                                            </p>
                                        </template>
                                        <template #description-data="{ row: law }">
                                            <p
                                                class="line-clamp-2 w-full max-w-sm whitespace-normal break-all hover:line-clamp-none"
                                            >
                                                {{ law.description }}
                                            </p>
                                        </template>
                                        <template #count-data="{ row: law }">
                                            <USelect
                                                name="count"
                                                :options="Array.from(Array(7).keys())"
                                                :model-value="selectedPenalties.find((p) => p.law.id === law.id)?.count ?? 0"
                                                @change="calculate({ law: law, count: parseInt($event) })"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            />
                                        </template>
                                    </UTable>
                                </div>
                            </template>
                        </UAccordion>
                    </dl>
                </div>
            </div>
        </div>

        <UDivider :label="$t('common.result')" />

        <div class="flow-root">
            <div class="overflow-x-auto sm:-mx-6 lg:-mx-8">
                <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                    <div class="text-xl">
                        <PenaltyStats :summary="summary" />

                        <div class="mt-4">
                            <PenaltySummaryTable
                                v-if="lawBooks && lawBooks.length > 0"
                                :law-books="lawBooks"
                                :selected-laws="selectedPenalties"
                            />
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <UButtonGroup class="mt-2 inline-flex w-full">
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" @click="copyToClipboard()">
                    {{ $t('common.copy') }}
                </UButton>
                <UButton color="red" @click="reset()">
                    {{ $t('common.reset') }}
                </UButton>
            </UButtonGroup>
        </UButtonGroup>
    </div>
</template>
