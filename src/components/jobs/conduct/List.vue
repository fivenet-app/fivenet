<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions, Switch } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { useConfirmDialog, watchDebounced } from '@vueuse/core';
import { CheckIcon } from 'mdi-vue3';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import { ConductEntry, ConductType } from '~~/gen/ts/resources/jobs/conduct';
import { User } from '~~/gen/ts/resources/users/users';
import { ConductListEntriesResponse } from '~~/gen/ts/services/jobs/jobs';
import CreateOrUpdateModal from '~/components/jobs/conduct/CreateOrUpdateModal.vue';
import ListEntry from '~/components/jobs/conduct/ListEntry.vue';
import { useJobsStore } from '~/store/jobs';

const { $grpc } = useNuxtApp();

const query = ref<{ types: ConductType[]; showExpired?: boolean; user_ids?: User[] }>({
    types: [],
    user_ids: [],
    showExpired: false,
});
const offset = ref(0n);

const { data, pending, refresh, error } = useLazyAsyncData(`jobs-conduct-${offset}`, () => listConductEntries());

async function listConductEntries(): Promise<ConductListEntriesResponse> {
    try {
        const call = $grpc.getJobsClient().conductListEntries({
            pagination: {
                offset: offset.value,
            },
            types: [],
            userIds: query.value.user_ids?.map((u) => u.userId) ?? [],
            showExpired: query.value.showExpired,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function deleteConductEntry(id: bigint): Promise<void> {
    try {
        const call = $grpc.getJobsClient().conductDeleteEntry({ id });
        await call;

        refresh();
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const queryTypes = ref('');

const entriesChars = ref<User[]>([]);
const queryTargets = ref<string>('');

const searchNameInput = ref<HTMLInputElement | null>(null);
function focusSearch(): void {
    if (searchNameInput.value) {
        searchNameInput.value.focus();
    }
}

watch(offset, async () => refresh());
watchDebounced(query.value, () => refresh(), { debounce: 600, maxWait: 1400 });

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

function updateEntryInPlace(entry: ConductEntry): void {
    if (data.value === null) {
        refresh();
        return;
    }

    const idx = data.value.entries.findIndex((e) => e.id === entry.id);
    if (idx !== undefined && idx > -1) {
        data.value.entries[idx] = entry;
    }
}

watchDebounced(
    queryTargets,
    async () => {
        await refreshColleagues();
        if (query.value.user_ids) colleagues.value?.users.unshift(...query.value.user_ids);
    },
    {
        debounce: 600,
        maxWait: 1400,
    },
);

onMounted(async () => {
    await refreshColleagues();
});

const open = ref(false);
const selectedEntry = ref<ConductEntry | undefined>();

const { isRevealed, reveal, confirm, cancel, onConfirm } = useConfirmDialog();

onConfirm(async (id) => deleteConductEntry(id));
</script>

<template>
    <div class="py-2 pb-14">
        <ConfirmDialog
            v-if="selectedEntry !== undefined"
            :open="isRevealed"
            :cancel="cancel"
            :confirm="() => confirm(selectedEntry!.id)"
        />

        <CreateOrUpdateModal
            :open="open"
            :entry="selectedEntry"
            @close="open = false"
            @created="data?.entries.unshift($event)"
            @update="updateEntryInPlace($event)"
        />

        <div class="px-1 sm:px-2 lg:px-4">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="">
                        <div class="flex flex-row gap-4 mx-auto">
                            <div class="flex-1 form-control">
                                <label for="searchName" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.search') }}
                                    {{ $t('common.target') }}
                                </label>
                                <div class="relative flex items-center mt-2">
                                    <Combobox v-model="query.user_ids" as="div" class="w-full mt-2" multiple nullable>
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
                                                v-if="entriesChars.length > 0"
                                                class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                                            >
                                                <ComboboxOption
                                                    v-for="char in entriesChars"
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
                                <label for="types" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.search') }}
                                    {{ $t('common.type') }}
                                </label>
                                <div class="relative flex items-center mt-2">
                                    <Combobox v-model="query.types" as="div" class="w-full mt-2" multiple nullable>
                                        <div class="relative">
                                            <ComboboxButton as="div">
                                                <ComboboxInput
                                                    autocomplete="off"
                                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    :display-value="
                                                        (cTypes: any) =>
                                                            cTypes ? (cTypes as ConductType[]).join(', ') : $t('common.na')
                                                    "
                                                    :placeholder="$t('common.type')"
                                                    @change="queryTypes = $event.target.value"
                                                    @focusin="focusTablet(true)"
                                                    @focusout="focusTablet(false)"
                                                />
                                            </ComboboxButton>

                                            <ComboboxOptions
                                                class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                                            >
                                                <ComboboxOption
                                                    v-for="cType in ConductType"
                                                    :key="cType.valueOf()"
                                                    v-slot="{ active, selected }"
                                                    :value="cType"
                                                    as="char"
                                                >
                                                    <li
                                                        :class="[
                                                            'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                            active ? 'bg-primary-500' : '',
                                                        ]"
                                                    >
                                                        <span :class="['block truncate', selected && 'font-semibold']">
                                                            {{ cType }}
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
                            <div class="flex-initial form-control">
                                <label for="show_expired" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('components.jobs.conduct.List.show_expired') }}
                                </label>
                                <div class="relative flex items-center mt-3">
                                    <Switch
                                        v-model="query.showExpired"
                                        :class="[
                                            query.showExpired ? 'bg-info-600' : 'bg-base-700',
                                            'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2',
                                        ]"
                                    >
                                        <span class="sr-only">
                                            {{ $t('components.jobs.conduct.List.show_expired') }}
                                        </span>
                                        <span
                                            aria-hidden="true"
                                            :class="[
                                                query.showExpired ? 'translate-x-5' : 'translate-x-0',
                                                'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-neutral ring-0 transition duration-200 ease-in-out',
                                            ]"
                                        />
                                    </Switch>
                                </div>
                            </div>
                            <div class="flex-initial form-control">
                                <label for="create" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.create') }}
                                </label>
                                <div class="relative flex items-center mt-3">
                                    <div v-if="can('JobsService.ConductCreateEntry')" class="flex-initial form-control">
                                        <button
                                            type="button"
                                            class="inline-flex px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                                            @click="
                                                selectedEntry = undefined;
                                                open = true;
                                            "
                                        >
                                            {{ $t('common.create') }}
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
            <div class="flow-root mt-2">
                <div class="mx-0 -my-2 overflow-x-auto">
                    <div class="inline-block min-w-full py-2 align-middle px-1">
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.conduct_register')])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.conduct_register')])"
                            :retry="refresh"
                        />
                        <DataNoDataBlock
                            v-else-if="data?.entries.length === 0"
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
                                            {{ $t('common.created_at') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.expires_at') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.type') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.description') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.target') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.creator') }}
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
                                    <ListEntry
                                        v-for="conduct in data?.entries"
                                        :key="conduct.id.toString()"
                                        :conduct="conduct"
                                        class="transition-colors hover:bg-neutral/5"
                                        @selected="
                                            selectedEntry = conduct;
                                            open = true;
                                        "
                                        @delete="
                                            selectedEntry = conduct;
                                            reveal();
                                        "
                                    />
                                </tbody>
                                <thead>
                                    <tr>
                                        <th
                                            scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0"
                                        >
                                            {{ $t('common.created_at') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.expires_at') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.type') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.description') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.target') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.creator') }}
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
