<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { watchDebounced } from '@vueuse/core';
import { CheckIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import { useCompletorStore } from '~/store/completor';
import { UserShort } from '~~/gen/ts/resources/users/users';
import { ViewAuditLogRequest, ViewAuditLogResponse } from '~~/gen/ts/services/rector/rector';
import AuditLogEntry from '~/components/rector/audit/AuditLogEntry.vue';
import GenericTable from '~/components/partials/elements/GenericTable.vue';

const { $grpc } = useNuxtApp();

const completorStore = useCompletorStore();

const query = ref<{
    from: string;
    to: string;
    method: string;
    service: string;
    search: string;
}>({
    from: '',
    to: '',
    method: '',
    service: '',
    search: '',
});
const offset = ref(0n);

async function viewAuditLog(): Promise<ViewAuditLogResponse> {
    const req: ViewAuditLogRequest = {
        pagination: {
            offset: offset.value,
        },
        userIds: [],
    };

    if (selectedCitizens.value.length > 0) {
        const users: number[] = [];
        selectedCitizens.value?.forEach((v) => users.push(v.userId));
        req.userIds = users;
    }

    if (query.value.from !== '') {
        req.from = toTimestamp(fromString(query.value.from)!);
    }
    if (query.value.from !== '') {
        req.to = toTimestamp(fromString(query.value.to)!);
    }

    if (query.value.method !== '') {
        req.method = query.value.method;
    }
    if (query.value.service !== '') {
        req.service = query.value.service;
    }

    if (query.value.search !== '') {
        req.search = query.value.search;
    }

    try {
        const call = $grpc.getRectorClient().viewAuditLog(req);
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const { data, pending, refresh, error } = useLazyAsyncData(`rector-audit-${offset}`, () => viewAuditLog());

const entriesCitizens = ref<UserShort[]>([]);
const queryCitizens = ref('');
const selectedCitizens = ref<UserShort[]>([]);

async function findChars(): Promise<void> {
    if (queryCitizens.value === '') {
        return;
    }

    entriesCitizens.value = await completorStore.completeCitizens({
        search: queryCitizens.value,
    });
    entriesCitizens.value.unshift(...selectedCitizens.value);
}

const searchInput = ref<HTMLInputElement | null>(null);
function focusSearch(): void {
    if (searchInput.value) {
        searchInput.value.focus();
    }
}

function charsGetDisplayValue(chars: UserShort[]): string {
    const cs: string[] = [];
    chars.forEach((c) => cs.push(`${c?.firstname} ${c?.lastname} (${c?.dateofbirth})`));

    return cs.join(', ');
}

watch(offset, async () => refresh());
watchDebounced(query.value, async () => refresh(), {
    debounce: 600,
    maxWait: 1400,
});
watchDebounced(selectedCitizens.value, async () => refresh(), {
    debounce: 600,
    maxWait: 1400,
});
watchDebounced(queryCitizens, async () => await findChars(), {
    debounce: 600,
    maxWait: 1400,
});
</script>

<template>
    <div class="py-2">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="refresh()">
                        <div class="mx-auto flex flex-row gap-4">
                            <div class="form-control flex-1">
                                <label for="from" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.time_range') }}: {{ $t('common.from') }}
                                </label>
                                <div class="relative mt-2 flex items-center">
                                    <input
                                        v-model="query.from"
                                        type="datetime-local"
                                        name="from"
                                        :placeholder="`${$t('common.time_range')} ${$t('common.from')}`"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                            </div>
                            <div class="form-control flex-1">
                                <label for="to" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.time_range') }}: {{ $t('common.to') }}
                                </label>
                                <div class="relative mt-2 flex items-center">
                                    <input
                                        v-model="query.to"
                                        type="datetime-local"
                                        name="to"
                                        :placeholder="`${$t('common.time_range')} ${$t('common.to')}`"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                            </div>
                            <div class="form-control flex-1">
                                <label for="users" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.user', 2) }}
                                </label>
                                <div class="relative mt-2 items-center">
                                    <Combobox v-model="selectedCitizens" as="div" multiple nullable>
                                        <div class="relative">
                                            <ComboboxButton as="div">
                                                <ComboboxInput
                                                    autocomplete="off"
                                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    :display-value="
                                                        (chars: any) => (chars ? charsGetDisplayValue(chars) : $t('common.na'))
                                                    "
                                                    :placeholder="$t('common.user', 2)"
                                                    @change="queryCitizens = $event.target.value"
                                                    @focusin="focusTablet(true)"
                                                    @focusout="focusTablet(false)"
                                                />
                                            </ComboboxButton>

                                            <ComboboxOptions
                                                v-if="entriesCitizens.length > 0"
                                                class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                            >
                                                <ComboboxOption
                                                    v-for="char in entriesCitizens"
                                                    :key="char?.userId"
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
                                                            {{ char?.firstname }}
                                                            {{ char?.lastname }}
                                                            ({{ char?.dateofbirth }})
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
                                <label for="service" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.service') }}
                                </label>
                                <div class="relative mt-2 flex items-center">
                                    <input
                                        v-model="query.service"
                                        type="text"
                                        name="service"
                                        :placeholder="$t('common.service')"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                            </div>
                            <div class="form-control flex-1">
                                <label for="method" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.method') }}
                                </label>
                                <div class="relative mt-2 flex items-center">
                                    <input
                                        v-model="query.method"
                                        type="text"
                                        name="method"
                                        :placeholder="$t('common.method')"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                            </div>
                            <div class="form-control flex-1">
                                <label for="data" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.data') }}
                                </label>
                                <div class="relative mt-2 flex items-center">
                                    <input
                                        ref="searchInput"
                                        v-model="query.search"
                                        type="text"
                                        name="data"
                                        :placeholder="$t('common.search')"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
            <div class="mt-2 flow-root">
                <div class="-my-2 mx-0 overflow-x-auto">
                    <div class="inline-block min-w-full px-1 py-2 align-middle">
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.audit_log', 2)])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.audit_log', 2)])"
                            :retry="refresh"
                        />
                        <DataNoDataBlock
                            v-else-if="data?.logs.length === 0"
                            :type="$t('common.audit_log', 2)"
                            :focus="focusSearch"
                        />

                        <template v-else>
                            <GenericTable>
                                <template #thead>
                                    <tr>
                                        <th
                                            scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-1"
                                        >
                                            {{ $t('common.id') }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.time') }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.user', 1) }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.service') }}/{{ $t('common.method') }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.state') }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.data') }}
                                        </th>
                                        <th
                                            scope="col"
                                            class="relative py-3.5 pl-3 pr-4 text-right text-sm font-semibold text-neutral sm:pr-0"
                                        >
                                            {{ $t('common.action', 2) }}
                                        </th>
                                    </tr>
                                </template>
                                <template #tbody>
                                    <AuditLogEntry v-for="log in data?.logs" :key="log.id" :log="log" />
                                </template>
                            </GenericTable>

                            <TablePagination
                                :pagination="data?.pagination"
                                :refresh="refresh"
                                @offset-change="offset = $event"
                            />
                        </template>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
