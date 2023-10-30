<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { watchDebounced } from '@vueuse/core';
import { CheckIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Divider from '~/components/partials/elements/Divider.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import * as google_protobuf_timestamp_pb from '~~/gen/ts/google/protobuf/timestamp';
import { TimeclockEntry } from '~~/gen/ts/resources/jobs/timeclock';
import { User } from '~~/gen/ts/resources/users/users';
import { TimeclockListEntriesRequest, TimeclockListEntriesResponse } from '~~/gen/ts/services/jobs/jobs';
import ListEntry from './ListEntry.vue';
import Stats from './Stats.vue';

const { $grpc } = useNuxtApp();

const query = ref<{
    user_ids?: User[];
    from?: string;
    to?: string;
}>({});
const offset = ref(0n);

const { data, pending, refresh, error } = useLazyAsyncData(`jobs-timeclock-${offset.value}`, () => listTimeclockEntries());

async function listTimeclockEntries(): Promise<TimeclockListEntriesResponse> {
    return new Promise(async (res, rej) => {
        try {
            const req: TimeclockListEntriesRequest = {
                pagination: {
                    offset: offset.value,
                },
                userIds: query.value.user_ids?.map((u) => u.userId) ?? [],
            };
            if (query.value.from) {
                req.from = {
                    timestamp: google_protobuf_timestamp_pb.Timestamp.fromDate(fromString(query.value.from)!),
                };
            }
            if (query.value.to) {
                req.to = {
                    timestamp: google_protobuf_timestamp_pb.Timestamp.fromDate(fromString(query.value.to)!),
                };
            }

            const call = $grpc.getJobsClient().timeclockListEntries(req);
            const { response } = await call;

            return res(response);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

type GroupedTimeClockEntries = { date: Date; key: string; entries: TimeclockEntry[] }[];
const grouped = computed(() => {
    const groups: GroupedTimeClockEntries = [];
    data.value?.entries.map((e) => {
        const date = toDate(e.date);
        const idx = groups.findIndex((g) => g.key === date.toString());
        if (idx === -1) {
            groups.push({
                date: date,
                entries: [e],
                key: date.toString(),
            });
        } else {
            groups[idx].entries.push(e);
        }
    });

    return groups;
});

const entriesChars = ref<User[]>([]);
const queryTargets = ref<string>('');

const searchNameInput = ref<HTMLInputElement | null>(null);
function focusSearch(): void {
    if (searchNameInput.value) {
        searchNameInput.value.focus();
    }
}

watch(offset, async () => refresh());
watchDebounced(query.value, async () => refresh(), { debounce: 600, maxWait: 1400 });

async function listColleagues(): Promise<User[]> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getJobsClient().colleaguesList({
                pagination: {
                    offset: offset.value,
                },
                searchName: queryTargets.value,
            });
            const { response } = await call;

            return res(response.users);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

function charsGetDisplayValue(chars: User[]): string {
    const cs: string[] = [];
    chars.forEach((c) => cs.push(`${c?.firstname} ${c?.lastname}`));

    return cs.join(', ');
}

watchDebounced(
    queryTargets,
    async () => {
        if (can('JobsService.TimeclockListEntries.Access.All')) {
            entriesChars.value = await listColleagues();
            if (query.value.user_ids) entriesChars.value.unshift(...query.value.user_ids);
        }
    },
    {
        debounce: 600,
        maxWait: 1400,
    },
);

onMounted(async () => {
    if (can('JobsService.TimeclockListEntries.Access.All')) {
        entriesChars.value = await listColleagues();
    }
});
</script>

<template>
    <div class="py-2 pb-14">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="refresh()">
                        <div class="flex flex-row gap-4 mx-auto">
                            <div v-if="can('JobsService.TimeclockListEntries.Access.All')" class="flex-1 form-control">
                                <label for="searchName" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.search') }}
                                    {{ $t('common.colleague', 1) }}
                                </label>
                                <div class="relative flex items-center mt-2">
                                    <Combobox as="div" v-model="query.user_ids" class="w-full" multiple nullable>
                                        <div class="relative">
                                            <ComboboxButton as="div">
                                                <ComboboxInput
                                                    autocomplete="off"
                                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    @change="queryTargets = $event.target.value"
                                                    :display-value="
                                                        (chars: any) => (chars ? charsGetDisplayValue(chars) : $t('common.na'))
                                                    "
                                                    :placeholder="$t('common.target')"
                                                    @focusin="focusTablet(true)"
                                                    @focusout="focusTablet(false)"
                                                />
                                            </ComboboxButton>

                                            <ComboboxOptions
                                                v-if="entriesChars.length > 0"
                                                class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                                            >
                                                <ComboboxOption
                                                    v-for="char in entriesChars"
                                                    :key="char.identifier"
                                                    :value="char"
                                                    as="char"
                                                    v-slot="{ active, selected }"
                                                >
                                                    <li
                                                        :class="[
                                                            'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                            active ? 'bg-primary-500' : '',
                                                        ]"
                                                    >
                                                        <span :class="['block truncate', selected && 'font-semibold']">
                                                            {{ char.firstname }} {{ char.lastname }}
                                                        </span>

                                                        <span
                                                            v-if="selected"
                                                            :class="[
                                                                active ? 'text-neutral' : 'text-primary-500',
                                                                'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                            ]"
                                                        >
                                                            <CheckIcon class="w-5 h-5" aria-hidden="true" />
                                                        </span>
                                                    </li>
                                                </ComboboxOption>
                                            </ComboboxOptions>
                                        </div>
                                    </Combobox>
                                </div>
                            </div>
                            <div class="flex-1 form-control">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.time_range') }}: {{ $t('common.from') }}
                                </label>
                                <div class="relative flex items-center mt-2">
                                    <input
                                        v-model="query.from"
                                        type="date"
                                        name="search"
                                        :placeholder="`${$t('common.time_range')} ${$t('common.from')}`"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    />
                                </div>
                            </div>
                            <div class="flex-1 form-control">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral"
                                    >{{ $t('common.time_range') }}:
                                    {{ $t('common.to') }}
                                </label>
                                <div class="relative flex items-center mt-2">
                                    <input
                                        v-model="query.to"
                                        type="date"
                                        name="search"
                                        :placeholder="`${$t('common.time_range')} ${$t('common.to')}`"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    />
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
            <div class="flow-root mt-2">
                <div class="mx-0 -my-2 overflow-x-auto">
                    <div class="inline-block min-w-full py-2 align-middle px-1">
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
                                            scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0"
                                        >
                                            {{ $t('common.date') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.name') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.time') }}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-base-800">
                                    <template v-for="group in grouped">
                                        <ListEntry
                                            v-for="(entry, idx) in group.entries"
                                            :key="entry.userId + toDate(entry.date).toString()"
                                            :entry="entry"
                                            class="transition-colors hover:bg-neutral/5"
                                            :first="idx === 0 ? group.date : undefined"
                                        />
                                    </template>
                                </tbody>
                                <thead>
                                    <tr>
                                        <th
                                            scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0"
                                        >
                                            {{ $t('common.date') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.name') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.time') }}
                                        </th>
                                    </tr>
                                </thead>
                            </table>

                            <TablePagination
                                :pagination="data?.pagination"
                                @offset-change="offset = $event"
                                :refresh="refresh"
                            />
                        </div>
                    </div>
                </div>
            </div>
            <div v-if="data && data.stats" class="flow-root mb-4">
                <div class="sm:flex sm:items-center">
                    <div class="sm:flex-auto">
                        <Divider :label="$t('components.jobs.timeclock.Stats.title')" />
                        <Stats :stats="data.stats" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
