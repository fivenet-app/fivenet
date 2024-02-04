<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { ChevronDownIcon, GavelIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import LawListEntry from '~/components/penaltycalculator/LawListEntry.vue';
import PenaltyStats from '~/components/penaltycalculator/PenaltyStats.vue';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import { Law } from '~~/gen/ts/resources/laws/laws';
import PenaltySummaryTable from '~/components/penaltycalculator/PenaltySummaryTable.vue';
import GenericDivider from '~/components/partials/elements/GenericDivider.vue';

const completorStore = useCompletorStore();
const notifications = useNotificatorStore();

const { t, d } = useI18n();

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

const rawQuery = ref('');
const query = computed(() => rawQuery.value.toLowerCase());
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
                .filter((p) => p.name.toLowerCase().includes(query.value) || p.description?.toLowerCase().includes(query.value))
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
        .filter((books) => books.laws.length > 0),
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

    notifications.dispatchNotification({
        title: { key: 'notifications.penaltycalculator.title', parameters: {} },
        content: { key: 'notifications.penaltycalculator.content', parameters: {} },
        type: 'info',
    });

    return copyToClipboardWrapper(text);
}

function reset(): void {
    rawQuery.value = '';
    selectedPenalties.value = [];
}
</script>

<template>
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="relative">
                <h3 class="text-2xl font-semibold leading-6">
                    {{ $t('components.penaltycalculator.title') }}
                </h3>
            </div>
            <div class="pb-4 sm:flex sm:items-center">
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
                        :icon="GavelIcon"
                        :type="`${$t('common.citizen', 1)} ${$t('common.activity')}`"
                        class="mt-5"
                    />

                    <div v-else class="divide-y divide-neutral/10">
                        <div class="mt-5">
                            <input
                                v-model="rawQuery"
                                type="text"
                                name="search"
                                :placeholder="$t('common.filter')"
                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </div>
                        <dl class="mt-5 space-y-2 divide-y divide-neutral/10">
                            <Disclosure
                                v-for="lawBook in filteredLawBooks"
                                v-slot="{ open }"
                                :key="`${lawBook.id}-${query}`"
                                as="div"
                                class="pt-3"
                                :default-open="query.length > 0"
                            >
                                <dt>
                                    <DisclosureButton class="flex w-full items-start justify-between text-left text-neutral">
                                        <span class="text-base font-semibold leading-7">
                                            {{ lawBook.name }}
                                            <span v-if="lawBook.description">
                                                {{ ' - ' + lawBook.description }}
                                            </span>
                                        </span>
                                        <span class="ml-6 flex h-7 items-center">
                                            <ChevronDownIcon
                                                :class="[open ? 'upsidedown' : '', 'h-5 w-5 transition-transform']"
                                                aria-hidden="true"
                                            />
                                        </span>
                                    </DisclosureButton>
                                </dt>
                                <DisclosurePanel as="dd" class="mt-2 px-4">
                                    <div class="mt-2 flow-root">
                                        <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                                            <div class="inline-block min-w-full align-middle sm:px-6 lg:px-8">
                                                <table class="min-w-full divide-y divide-base-600">
                                                    <thead>
                                                        <tr>
                                                            <th
                                                                scope="col"
                                                                class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-1"
                                                            >
                                                                {{ $t('common.crime') }}
                                                            </th>
                                                            <th
                                                                scope="col"
                                                                class="px-2 py-3.5 text-left text-sm font-semibold text-neutral"
                                                            >
                                                                {{ $t('common.fine') }}
                                                            </th>
                                                            <th
                                                                scope="col"
                                                                class="px-2 py-3.5 text-left text-sm font-semibold text-neutral"
                                                            >
                                                                {{ $t('common.detention_time') }}
                                                            </th>
                                                            <th
                                                                scope="col"
                                                                class="px-2 py-3.5 text-left text-sm font-semibold text-neutral"
                                                            >
                                                                {{ $t('common.traffic_infraction_points', 2) }}
                                                            </th>
                                                            <th
                                                                scope="col"
                                                                class="px-2 py-3.5 text-left text-sm font-semibold text-neutral"
                                                            >
                                                                {{ $t('common.description') }}
                                                            </th>
                                                            <th
                                                                scope="col"
                                                                class="relative py-3.5 pl-3 pr-4 text-right text-sm font-semibold text-neutral sm:pr-0"
                                                            >
                                                                {{ $t('common.count') }}
                                                            </th>
                                                        </tr>
                                                    </thead>
                                                    <tbody class="divide-y divide-base-800">
                                                        <LawListEntry
                                                            v-for="law in lawBook.laws"
                                                            :key="law.id"
                                                            :law="law"
                                                            :count="
                                                                selectedPenalties.find((p) => p.law.id === law.id)?.count ?? 0
                                                            "
                                                            @selected="calculate($event)"
                                                        />
                                                    </tbody>
                                                </table>
                                            </div>
                                        </div>
                                    </div>
                                </DisclosurePanel>
                            </Disclosure>
                        </dl>
                    </div>
                </div>
            </div>
            <GenericDivider :label="$t('common.result')" />
            <div class="mt-2 flow-root">
                <div class="overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                        <div class="text-xl text-neutral">
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
            <div class="mt-2 flow-root">
                <div class="flex items-center gap-1">
                    <button
                        type="button"
                        class="flex-1 rounded-md bg-info-700 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-info-600"
                        @click="copyToClipboard()"
                    >
                        {{ $t('common.copy') }}
                    </button>
                    <button
                        type="button"
                        class="rounded-md bg-error-700 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-error-600"
                        @click="reset()"
                    >
                        {{ $t('common.reset') }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
