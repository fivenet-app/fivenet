<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiChevronDown } from '@mdi/js';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { useClipboard } from '@vueuse/core';
import ListEntry from '~/components/penaltycalculator/ListEntry.vue';
import Stats from '~/components/penaltycalculator/Stats.vue';
import { useNotificationsStore } from '~/store/notifications';
import { PenaltiesSummary, PenaltyCategory, SelectedPenalty } from '~/utils/penalty';
import SummaryTable from './SummaryTable.vue';

const { $grpc } = useNuxtApp();

const { t, d } = useI18n();

const clipboard = useClipboard();
const notifications = useNotificationsStore();

const { data: lawBooks, pending, refresh, error } = useLazyAsyncData(`lawbooks`, () => listLawBooks());

async function listLawBooks(): Promise<PenaltyCategory[]> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCompletorClient().listLawBooks({});
            const { response } = await call;

            return res(response.books);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const queryPenalities = ref<string>('');
const filteredLawBooks = ref<Array<PenaltyCategory>>([]);
const selectedPenalties = ref<Array<SelectedPenalty>>([]);

const summary = ref<PenaltiesSummary>({
    fine: 0n,
    detentionTime: 0n,
    stvoPoints: 0n,
    count: 0n,
});

async function applyQuery(): Promise<void> {
    if (!lawBooks.value) return;

    let newLawBooks = structuredClone(toRaw(lawBooks.value));

    filteredLawBooks.value = newLawBooks!.map((ps) => {
        const penalties = ps.laws.map((p) => {
            const show =
                p.name.toLowerCase().includes(queryPenalities.value.toLowerCase()) ||
                p.description.toLowerCase().includes(queryPenalities.value.toLowerCase())
                    ? true
                    : false;
            return {
                ...p,
                show,
            };
        });

        const show = penalties.some((p) => p.show === true);

        return {
            ...ps,
            penalties,
            show,
        };
    });
}

function getNameForLawBookId(id: bigint): string | undefined {
    return lawBooks.value?.filter((b) => b.id === id)[0].name;
}

watch(lawBooks, async () => applyQuery());
watch(queryPenalities, async () => applyQuery());

function calculate(e: SelectedPenalty): void {
    const idx = selectedPenalties.value.findIndex((v) => v.law.lawbookId === e.law.lawbookId && v.law.name === e.law.name);
    let count = e.count;
    if (idx > -1) {
        const existing = selectedPenalties.value.at(idx)!;
        selectedPenalties.value[idx] = e;
        if (existing.count != e.count) {
            count = e.count - existing.count;
        }
    } else {
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

${t('components.penaltycalculator.fine')}: $${summary.value.fine}
${t('components.penaltycalculator.detention_time')}: ${summary.value.detentionTime} ${t(
            'common.time_ago.month',
            summary.value.detentionTime.toString()
        )}
${t('common.traffic_infraction_points', 2)}: ${summary.value.stvoPoints}
${t('common.total_count')}: ${summary.value.count}
`;

    if (selectedPenalties.value.length > 0) {
        text += `
${t('components.penaltycalculator.crime', selectedPenalties.value.length)}:
`;

        selectedPenalties.value.forEach((v) => {
            text += `* ${getNameForLawBookId(v.law.lawbookId)} - ${v.law.name} (${v.count}x)
`;
        });
    }

    notifications.dispatchNotification({
        title: { key: 'notifications.penaltycalculator.title', parameters: [] },
        content: { key: 'notifications.penaltycalculator.content', parameters: [] },
        type: 'info',
    });

    return clipboard.copy(text);
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
            <div class="sm:flex sm:items-center pb-4">
                <div class="sm:flex-auto">
                    <div class="divide-y divide-white/10">
                        <div class="mt-5">
                            <input
                                v-model="queryPenalities"
                                type="text"
                                name="search"
                                :placeholder="$t('common.filter')"
                                class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            />
                        </div>
                        <dl class="mt-5 space-y-2 divide-y divide-white/10">
                            <Disclosure
                                as="div"
                                v-for="lawBook in filteredLawBooks"
                                :key="lawBook.id.toString()"
                                class="pt-3"
                                v-slot="{ open }"
                                v-show="lawBook.show"
                            >
                                <dt>
                                    <DisclosureButton class="flex w-full items-start justify-between text-left text-white">
                                        <span class="text-base font-semibold leading-7">
                                            {{ lawBook.name }}
                                            <span v-if="lawBook.description">
                                                {{ ' - ' + lawBook.description }}
                                            </span>
                                        </span>
                                        <span class="ml-6 flex h-7 items-center">
                                            <SvgIcon
                                                :class="[open ? 'upsidedown' : '', 'h-6 w-6 transition-transform']"
                                                aria-hidden="true"
                                                type="mdi"
                                                :path="mdiChevronDown"
                                            />
                                        </span>
                                    </DisclosureButton>
                                </dt>
                                <DisclosurePanel as="dd" class="mt-2 px-4">
                                    <div class="flow-root mt-2">
                                        <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                                            <div class="inline-block min-w-full align-middle sm:px-6 lg:px-8">
                                                <table class="min-w-full divide-y divide-base-600">
                                                    <thead>
                                                        <tr>
                                                            <th
                                                                scope="col"
                                                                class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0"
                                                            >
                                                                {{ $t('components.penaltycalculator.crime') }}
                                                            </th>
                                                            <th
                                                                scope="col"
                                                                class="py-3.5 px-2 text-left text-sm font-semibold text-neutral"
                                                            >
                                                                {{ $t('components.penaltycalculator.fine') }}
                                                            </th>
                                                            <th
                                                                scope="col"
                                                                class="py-3.5 px-2 text-left text-sm font-semibold text-neutral"
                                                            >
                                                                {{ $t('components.penaltycalculator.detention_time') }}
                                                            </th>
                                                            <th
                                                                scope="col"
                                                                class="py-3.5 px-2 text-left text-sm font-semibold text-neutral"
                                                            >
                                                                {{ $t('common.traffic_infraction_points', 2) }}
                                                            </th>
                                                            <th
                                                                scope="col"
                                                                class="py-3.5 px-2 text-left text-sm font-semibold text-neutral"
                                                            >
                                                                {{ $t('common.other') }}
                                                            </th>
                                                            <th
                                                                scope="col"
                                                                class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral"
                                                            >
                                                                {{ $t('common.count') }}
                                                            </th>
                                                        </tr>
                                                    </thead>
                                                    <tbody class="divide-y divide-base-800">
                                                        <ListEntry
                                                            v-for="law in lawBook.laws"
                                                            :key="law.id.toString()"
                                                            :law="law"
                                                            @selected="calculate($event)"
                                                            v-show="law.show === undefined || law.show"
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
            <div class="relative">
                <div class="absolute inset-0 flex items-center" aria-hidden="true">
                    <div class="w-full border-t border-gray-300" />
                </div>
                <div class="relative flex justify-center">
                    <span class="bg-white px-3 text-base font-semibold leading-6 text-gray-900">
                        {{ $t('common.result') }}
                    </span>
                </div>
            </div>
            <div class="flow-root mt-2">
                <div class="overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                        <div class="text-neutral text-xl">
                            <Stats :summary="summary" />
                            <div class="mt-4">
                                <SummaryTable
                                    v-if="lawBooks && lawBooks.length > 0"
                                    :law-books="lawBooks"
                                    :selected-laws="selectedPenalties"
                                />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div class="flow-root mt-2">
                <div class="flex items-center">
                    <button
                        type="button"
                        @click="copyToClipboard()"
                        class="flex-1 rounded-md bg-info-700 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-info-600"
                    >
                        {{ $t('common.copy') }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
