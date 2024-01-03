<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { watchDebounced } from '@vueuse/core';
import { CheckIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import * as googleProtobufTimestamp from '~~/gen/ts/google/protobuf/timestamp';
import { User } from '~~/gen/ts/resources/users/users';
import { RequestsListEntriesRequest, RequestsListEntriesResponse } from '~~/gen/ts/services/jobs/jobs';
import ListEntry from '~/components/jobs/requests/ListEntry.vue';
import { useJobsStore } from '~/store/jobs';

const { $grpc } = useNuxtApp();

const query = ref<{
    user_ids?: User[];
    from?: string;
    to?: string;
}>({});
const offset = ref(0n);

const { data, pending, refresh, error } = useLazyAsyncData(`jobs-requests-${offset.value}`, () => listRequests());

async function listRequests(): Promise<RequestsListEntriesResponse> {
    try {
        const req: RequestsListEntriesRequest = {
            pagination: {
                offset: offset.value,
            },
            userIds: query.value.user_ids?.map((u) => u.userId) ?? [],
        };
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

        const call = $grpc.getJobsClient().requestsListEntries(req);
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const queryTargets = ref<string>('');

const searchNameInput = ref<HTMLInputElement | null>(null);
function focusSearch(): void {
    if (searchNameInput.value) {
        searchNameInput.value.focus();
    }
}

watch(offset, async () => refresh());
watchDebounced(query.value, async () => refresh(), { debounce: 600, maxWait: 1400 });

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
    chars.forEach((c) => cs.push(`${c?.firstname} ${c?.lastname}`));

    return cs.join(', ');
}

watchDebounced(
    queryTargets,
    async () => {
        if (can('JobsService.RequestsListEntries.Access.All')) {
            await refreshColleagues();
        }
    },
    {
        debounce: 600,
        maxWait: 1400,
    },
);

onMounted(async () => {
    if (can('JobsService.RequestsListEntries.Access.All')) {
        await refreshColleagues();
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
                            <div v-if="can('JobsService.RequestsListEntries.Access.All')" class="flex-1 form-control">
                                <label for="searchName" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.search') }}
                                    {{ $t('common.colleague', 1) }}
                                </label>
                                <div class="relative flex items-center mt-2">
                                    <Combobox v-model="query.user_ids" as="div" class="w-full" multiple nullable>
                                        <div class="relative">
                                            <ComboboxButton as="div">
                                                <ComboboxInput
                                                    autocomplete="off"
                                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
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
                                                v-if="colleagues?.users && colleagues.users.length > 0"
                                                class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                                            >
                                                <ComboboxOption
                                                    v-for="char in colleagues.users"
                                                    v-slot="{ active, selected }"
                                                    :key="char.identifier"
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
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
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
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
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
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.request', 2)])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.request', 2)])"
                            :retry="refresh"
                        />
                        <DataNoDataBlock
                            v-else-if="data && data.entries?.length === 0"
                            :focus="focusSearch"
                            :message="$t('components.citizens.citizens_list.no_citizens')"
                        />
                        <div v-else>
                            <ul role="list" class="flex flex-col">
                                <ListEntry v-for="request in data?.entries" :key="request.id" :request="request" />
                            </ul>

                            <TablePagination
                                :pagination="data?.pagination"
                                :refresh="refresh"
                                @offset-change="offset = $event"
                            />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
