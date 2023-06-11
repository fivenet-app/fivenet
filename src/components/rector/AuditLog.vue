<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { CheckIcon } from '@heroicons/vue/20/solid';
import { MagnifyingGlassIcon } from '@heroicons/vue/24/solid';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { watchDebounced } from '@vueuse/core';
import DataErrorBlock from '~/components//partials/DataErrorBlock.vue';
import TablePagination from '~/components//partials/TablePagination.vue';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import * as google_protobuf_timestamp_pb from '~~/gen/ts/google/protobuf/timestamp';
import { PaginationResponse } from '~~/gen/ts/resources/common/database/database';
import { AuditEntry } from '~~/gen/ts/resources/rector/audit';
import { UserShort } from '~~/gen/ts/resources/users/users';
import { ViewAuditLogRequest } from '~~/gen/ts/services/rector/rector';
import AuditLogEntry from './AuditLogEntry.vue';

const { $grpc } = useNuxtApp();

const query = ref<{ from: string; to: string; search: string }>({ from: '', to: '', search: '' });
const pagination = ref<PaginationResponse>();
const offset = ref(0n);

async function getAuditLog(): Promise<Array<AuditEntry>> {
    return new Promise(async (res, rej) => {
        const req: ViewAuditLogRequest = {
            pagination: {
                offset: offset.value,
            },
            userIds: [],
        };
        const users = new Array<number>();
        selectedChars.value?.forEach((v) => users.push(v.userId));
        req.userIds = users;

        if (query.value.search !== '') {
            req.search = query.value.search;
        }

        if (query.value.from !== '') {
            req.from = {
                timestamp: google_protobuf_timestamp_pb.Timestamp.fromDate(fromString(query.value.from)!),
            };
        }
        if (query.value.from !== '') {
            req.to = {
                timestamp: google_protobuf_timestamp_pb.Timestamp.fromDate(fromString(query.value.to)!),
            };
        }

        try {
            const call = $grpc.getRectorClient().viewAuditLog(req);
            const { response } = await call;

            pagination.value = response.pagination;
            return res(response.logs);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const { data: logs, pending, refresh, error } = useLazyAsyncData(`rector-audit-${offset}`, () => getAuditLog());

const entriesChars = ref<UserShort[]>([]);
const queryChar = ref('');
const selectedChars = ref<undefined | UserShort[]>([]);

async function findChars(): Promise<void> {
    if (queryChar.value === '') {
        return;
    }

    const call = $grpc.getCompletorClient().completeCitizens({
        search: queryChar.value,
    });
    const { response } = await call;

    entriesChars.value = response.users;
}

const searchInput = ref<HTMLInputElement | null>(null);
function focusSearch(): void {
    if (searchInput.value) {
        searchInput.value.focus();
    }
}

function charsGetDisplayValue(chars: UserShort[]): string {
    const cs = new Array<string>();
    chars.forEach((c) => cs.push(`${c?.firstname} ${c?.lastname}`));

    return cs.join(', ');
}

watch(offset, async () => refresh());
watchDebounced(query.value, async () => refresh(), {
    debounce: 600,
    maxWait: 1400,
});
watchDebounced(queryChar, async () => await findChars(), {
    debounce: 600,
    maxWait: 1400,
});
</script>

<template>
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="refresh()">
                        <div class="flex flex-row gap-4 mx-auto">
                            <div class="flex-1 form-control">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.time_range') }}: {{ $t('common.from') }}
                                </label>
                                <div class="relative flex items-center mt-2">
                                    <input
                                        v-model="query.from"
                                        ref="searchInput"
                                        type="datetime-local"
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
                                        ref="searchInput"
                                        type="datetime-local"
                                        name="search"
                                        :placeholder="`${$t('common.time_range')} ${$t('common.to')}`"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    />
                                </div>
                            </div>
                            <div class="flex-1 form-control">
                                <label for="users" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.user', 2) }}
                                </label>
                                <div class="relative items-center mt-2">
                                    <Combobox as="div" v-model="selectedChars" multiple nullable>
                                        <div class="relative">
                                            <ComboboxButton as="div">
                                                <ComboboxInput
                                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    @change="queryChar = $event.target.value"
                                                    :display-value="(chars: any) => chars ? charsGetDisplayValue(chars) : 'N/A'"
                                                    :placeholder="$t('common.user', 2)"
                                                />
                                            </ComboboxButton>

                                            <ComboboxOptions
                                                v-if="entriesChars.length > 0"
                                                class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm"
                                            >
                                                <ComboboxOption
                                                    v-for="char in entriesChars"
                                                    :key="char?.identifier"
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
                                                            {{ char?.firstname }}
                                                            {{ char?.lastname }}
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
                        </div>
                    </form>
                </div>
            </div>
            <div class="flow-root mt-2">
                <div class="-mx-4 -my-2 overflow-x-scroll sm:-mx-6 lg:-mx-8">
                    <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.audit_log', 2)])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.audit_log', 2)])"
                            :retry="refresh"
                        />
                        <button
                            v-else-if="logs && logs.length === 0"
                            type="button"
                            @click="focusSearch"
                            class="relative block w-full p-12 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
                        >
                            <MagnifyingGlassIcon class="w-12 h-12 mx-auto text-neutral" />
                            <span class="block mt-2 text-sm font-semibold text-gray-300">
                                {{ $t('common.not_found', [$t('common.audit_log', 2)]) }}
                            </span>
                        </button>
                        <div v-else>
                            <table class="min-w-full divide-y divide-base-600">
                                <thead>
                                    <tr>
                                        <th
                                            scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0"
                                        >
                                            {{ $t('common.id') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.time') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.user', 1) }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.service') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.state') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.data') }}
                                        </th>
                                        <th
                                            scope="col"
                                            class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral"
                                        >
                                            {{ $t('common.action', 2) }}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-base-800">
                                    <AuditLogEntry
                                        v-for="log in logs"
                                        :key="log.id?.toString()"
                                        :log="log"
                                        class="transition-colors hover:bg-neutral/5"
                                    />
                                </tbody>
                                <thead>
                                    <tr>
                                        <th
                                            scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0"
                                        >
                                            {{ $t('common.id') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.time') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.user', 1) }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.service') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.state') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.data') }}
                                        </th>
                                        <th
                                            scope="col"
                                            class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral"
                                        >
                                            {{ $t('common.action', 2) }}
                                        </th>
                                    </tr>
                                </thead>
                            </table>

                            <TablePagination :pagination="pagination" @offset-change="offset = $event" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
