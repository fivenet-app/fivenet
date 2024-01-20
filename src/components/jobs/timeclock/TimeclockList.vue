<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { watchDebounced } from '@vueuse/core';
import { CalendarIcon, CheckIcon, ChevronLeftIcon, ChevronRightIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericDivider from '~/components/partials/elements/GenericDivider.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import * as googleProtobufTimestamp from '~~/gen/ts/google/protobuf/timestamp';
import { TimeclockEntry } from '~~/gen/ts/resources/jobs/timeclock';
import { User } from '~~/gen/ts/resources/users/users';
import TimeclockListEntry from '~/components/jobs/timeclock/TimeclockListEntry.vue';
import TimeclockStatsBlock from '~/components/jobs/timeclock/TimeclockStatsBlock.vue';
import { useJobsStore } from '~/store/jobs';
import { dateToDateString } from '~/utils/time';
import type { ListTimeclockRequest, ListTimeclockResponse } from '~~/gen/ts/services/jobs/timeclock';

const { $grpc } = useNuxtApp();

const canAccessAll = can('JobsTimeclockService.ListTimeclock.Access.All');

const now = new Date();
const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
const currentDay = ref(new Date(today.getFullYear(), today.getMonth(), today.getDate()));

const futureDay = ref(new Date(currentDay.value.getFullYear(), currentDay.value.getMonth(), currentDay.value.getDate() + 1));
const previousDay = ref(new Date(currentDay.value.getFullYear(), currentDay.value.getMonth(), currentDay.value.getDate() - 1));

const perDay = ref(canAccessAll);

const query = ref<{
    user_ids?: User[];
    from?: string;
    to?: string;
    perDay: boolean;
}>({
    from: dateToDateString(currentDay.value),
    to: canAccessAll ? dateToDateString(previousDay.value) : undefined,
    perDay: perDay.value,
});
const offset = ref(0n);

const { data, pending, refresh, error } = useLazyAsyncData(`jobs-timeclock-${offset.value}`, () => listTimeclockEntries());

async function listTimeclockEntries(): Promise<ListTimeclockResponse> {
    try {
        const req: ListTimeclockRequest = {
            pagination: {
                offset: offset.value,
            },
            userIds: query.value.user_ids?.map((u) => u.userId) ?? [],
        };
        if (query.value.perDay !== undefined) {
            req.perDay = query.value.perDay;
        }
        if (query.value.from) {
            req.from = {
                timestamp: googleProtobufTimestamp.Timestamp.fromDate(fromString(query.value.from)!),
            };
        }
        if (query.value.to) {
            req.to = {
                timestamp: googleProtobufTimestamp.Timestamp.fromDate(fromString(query.value.to)!),
            };
        }

        const call = $grpc.getJobsTimeclockClient().listTimeclock(req);
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

type GroupedTimeClockEntries = { date: Date; key: string; entries: TimeclockEntry[] }[];
const grouped = computed(() => {
    const groups: GroupedTimeClockEntries = [];
    data.value?.entries.forEach((e) => {
        const date = toDate(e.date);
        const idx = groups.findIndex((g) => g.key === date.toString());
        if (idx === -1) {
            groups.push({
                date,
                entries: [e],
                key: date.toString(),
            });
        } else {
            groups[idx].entries.push(e);
        }
    });

    return groups;
});

const queryTargets = ref<string>('');

const searchNameInput = ref<HTMLInputElement | null>(null);
function focusSearch(): void {
    if (searchNameInput.value) {
        searchNameInput.value.focus();
    }
}

watch(offset, async () => refresh());
watchDebounced(
    query.value,
    async () => {
        if (canAccessAll) {
            if (query.value.user_ids !== undefined && query.value.user_ids.length > 0) {
                if (perDay.value) {
                    perDay.value = false;
                    query.value.to = undefined;
                }
            } else {
                perDay.value = true;
            }
        } else {
            perDay.value = false;
        }

        return refresh();
    },
    { debounce: 600, maxWait: 1400 },
);

const jobsStore = useJobsStore();
const { data: colleagues, refresh: refreshColleagues } = useLazyAsyncData(
    `jobs-colleagues-0-${queryTargets.value}`,
    () =>
        jobsStore.listColleagues({
            pagination: { offset: 0n },
            searchName: queryTargets.value,
        }),
    {
        immediate: false,
    },
);

function charsGetDisplayValue(chars: User[]): string {
    const cs: string[] = [];
    chars.forEach((c) => cs.push(`${c?.firstname} ${c?.lastname} (${c?.dateofbirth})`));

    return cs.join(', ');
}

watchDebounced(
    queryTargets,
    async () => {
        if (canAccessAll) {
            await refreshColleagues();
            if (query.value.user_ids) colleagues.value?.users.unshift(...query.value.user_ids);
        }
    },
    {
        debounce: 600,
        maxWait: 1400,
    },
);

onMounted(async () => {
    if (canAccessAll) {
        await refreshColleagues();
    }
});

function dayForward(): void {
    currentDay.value.setDate(currentDay.value.getDate() + 1);
    currentDay.value = new Date(currentDay.value);

    updateDates();
}

function dayBackwards(): void {
    currentDay.value.setDate(currentDay.value.getDate() - 1);
    currentDay.value = new Date(currentDay.value);

    updateDates();
}

function updateDates(): void {
    futureDay.value.setTime(
        new Date(currentDay.value.getFullYear(), currentDay.value.getMonth(), currentDay.value.getDate() + 1).getTime(),
    );
    futureDay.value = new Date(futureDay.value);
    previousDay.value.setTime(
        new Date(currentDay.value.getFullYear(), currentDay.value.getMonth(), currentDay.value.getDate() - 1).getTime(),
    );
    previousDay.value = new Date(previousDay.value);

    query.value.from = dateToDateString(currentDay.value);
    query.value.to = dateToDateString(previousDay.value);
}
</script>

<template>
    <div class="py-2 pb-14">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="refresh()">
                        <div class="mx-auto flex flex-row gap-4">
                            <div v-if="canAccessAll" class="form-control flex-1">
                                <label for="searchName" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.search') }}
                                    {{ $t('common.colleague', 1) }}
                                </label>
                                <div class="relative mt-2 flex items-center">
                                    <Combobox v-model="query.user_ids" as="div" class="w-full" multiple nullable>
                                        <div class="relative">
                                            <ComboboxButton as="div">
                                                <ComboboxInput
                                                    autocomplete="off"
                                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    :display-value="
                                                        (chars: any) => (chars ? charsGetDisplayValue(chars) : $t('common.na'))
                                                    "
                                                    :placeholder="$t('common.target')"
                                                    @change="queryTargets = $event.target.value"
                                                    @focusin="focusTablet(true)"
                                                    @focusout="focusTablet(false)"
                                                />
                                            </ComboboxButton>

                                            <ComboboxOptions
                                                v-if="colleagues?.users && colleagues?.users?.length > 0"
                                                class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                            >
                                                <ComboboxOption
                                                    v-for="char in colleagues?.users"
                                                    :key="char.identifier"
                                                    v-slot="{ active, selected }"
                                                    :value="char"
                                                    as="char"
                                                >
                                                    <li
                                                        :class="[
                                                            'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                            active ? 'bg-primary-500' : '',
                                                        ]"
                                                    >
                                                        <span :class="['block truncate', selected && 'font-semibold']">
                                                            {{ char.firstname }} {{ char.lastname }} ({{ char?.dateofbirth }})
                                                        </span>

                                                        <span
                                                            v-if="selected"
                                                            :class="[
                                                                active ? 'text-neutral' : 'text-primary-500',
                                                                'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                            ]"
                                                        >
                                                            <CheckIcon class="h-5 w-5" aria-hidden="true" />
                                                        </span>
                                                    </li>
                                                </ComboboxOption>
                                            </ComboboxOptions>
                                        </div>
                                    </Combobox>
                                </div>
                            </div>
                            <div class="form-control flex-1">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                    <template v-if="perDay"> {{ $t('common.date') }}: </template>
                                    <template v-else>
                                        {{ $t('common.time_range') }}:
                                        {{ $t('common.from') }}
                                    </template>
                                </label>
                                <div class="relative mt-2 flex items-center">
                                    <input
                                        v-model="query.from"
                                        type="date"
                                        name="search"
                                        :placeholder="`${$t('common.time_range')} ${$t('common.from')}`"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                            </div>
                            <div v-if="!perDay" class="form-control flex-1">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.time_range') }}:
                                    {{ $t('common.to') }}
                                </label>
                                <div class="relative mt-2 flex items-center">
                                    <input
                                        v-model="query.to"
                                        type="date"
                                        name="search"
                                        :placeholder="`${$t('common.time_range')} ${$t('common.to')}`"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                            </div>
                        </div>
                        <div v-if="perDay" class="mx-auto flex flex-row gap-4 pt-2">
                            <div class="form-control flex-1">
                                <button
                                    type="button"
                                    :disabled="futureDay > today"
                                    :class="[
                                        futureDay > today
                                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                            : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                                        'relative inline-flex w-full cursor-pointer place-content-start items-center rounded-md px-3 py-2 text-sm font-semibold text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2',
                                    ]"
                                    @click="dayForward()"
                                >
                                    <ChevronLeftIcon class="h-5 w-5" />
                                    {{ $t('common.forward') }} - {{ $d(futureDay, 'date') }}
                                </button>
                            </div>
                            <div class="form-control flex-initial">
                                <button
                                    type="button"
                                    disabled
                                    class="disabled relative inline-flex inline-flex w-full cursor-pointer place-content-end items-center items-center rounded-md bg-base-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-base-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-500"
                                >
                                    <CalendarIcon class="mr-1 h-5 w-5" />
                                    {{ $d(currentDay, 'date') }}
                                </button>
                            </div>
                            <div class="form-control flex-1">
                                <button
                                    type="button"
                                    class="relative inline-flex w-full cursor-pointer place-content-end items-center rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                                    @click="dayBackwards()"
                                >
                                    {{ $d(previousDay, 'date') }} - {{ $t('common.previous') }}
                                    <ChevronRightIcon class="h-5 w-5" />
                                </button>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
            <div class="mt-2 flow-root">
                <div class="-my-2 mx-0 overflow-x-auto">
                    <div class="inline-block min-w-full px-1 py-2 align-middle">
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.timeclock', 2)])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.timeclock', 2)])"
                            :retry="refresh"
                        />
                        <DataNoDataBlock
                            v-else-if="data && data.entries && data.entries.length === 0"
                            :focus="focusSearch"
                            :message="$t('components.citizens.citizens_list.no_citizens')"
                        />
                        <div v-else>
                            <table class="min-w-full divide-y divide-base-600">
                                <thead>
                                    <tr>
                                        <th
                                            v-if="!perDay"
                                            scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-1"
                                        >
                                            {{ $t('common.date') }}
                                        </th>
                                        <th v-else scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.name') }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.time') }}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-base-800">
                                    <template v-for="group in grouped" :key="group.key">
                                        <TimeclockListEntry
                                            v-for="(entry, idx) in group.entries"
                                            :key="entry.userId + toDate(entry.date).toString()"
                                            :entry="entry"
                                            :first="idx === 0 ? group.date : undefined"
                                            :show-date="!perDay"
                                        />
                                    </template>
                                </tbody>
                                <thead>
                                    <tr>
                                        <th
                                            v-if="!perDay"
                                            scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-1"
                                        >
                                            {{ $t('common.date') }}
                                        </th>
                                        <th v-else scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.name') }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.time') }}
                                        </th>
                                    </tr>
                                </thead>
                            </table>

                            <TablePagination
                                :pagination="data?.pagination"
                                :refresh="refresh"
                                @offset-change="offset = $event"
                            />
                        </div>
                    </div>
                </div>
            </div>
            <div v-if="data && data.stats" class="mb-4 flow-root">
                <div class="sm:flex sm:items-center">
                    <div class="sm:flex-auto">
                        <GenericDivider :label="$t('components.jobs.timeclock.Stats.title')" />
                        <TimeclockStatsBlock :stats="data.stats" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
