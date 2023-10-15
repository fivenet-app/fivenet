<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel, Switch } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { watchDebounced } from '@vueuse/core';
import { useRouteHash } from '@vueuse/router';
import { vMaska } from 'maska';
import { ChevronDownIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import { ListCitizensRequest, ListCitizensResponse } from '~~/gen/ts/services/citizenstore/citizenstore';
import ListEntry from './ListEntry.vue';

const { $grpc } = useNuxtApp();

const query = ref<{ name?: string; phoneNumber?: string; wanted?: boolean; trafficPoints?: number; dateofbirth?: string }>({});
const offset = ref(0n);

const hash = useRouteHash();
if (hash.value !== undefined && hash.value !== null) {
    query.value = unmarshalHashToObject(hash.value as string);
}

const { data, pending, refresh, error } = useLazyAsyncData(
    `citizens-${offset.value}-${query.value.name}-${query.value.wanted}-${query.value.phoneNumber}`,
    () => {
        hash.value = marshalObjectToHash(query.value);
        return listCitizens();
    },
);

async function listCitizens(): Promise<ListCitizensResponse> {
    return new Promise(async (res, rej) => {
        try {
            const req: ListCitizensRequest = {
                pagination: {
                    offset: offset.value,
                },
                searchName: query.value.name ?? '',
            };
            if (query.value.wanted) {
                req.wanted = query.value.wanted;
            }
            if (query.value.phoneNumber) {
                req.phoneNumber = query.value.phoneNumber;
            }
            if (query.value.trafficPoints) {
                req.trafficPoints = BigInt(query.value.trafficPoints?.toString() ?? '0');
            }
            if (query.value.dateofbirth) {
                req.dateofbirth = query.value.dateofbirth;
            }

            const call = $grpc.getCitizenStoreClient().listCitizens(req);
            const { response } = await call;

            return res(response);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const searchNameInput = ref<HTMLInputElement | null>(null);
function focusSearch(): void {
    if (searchNameInput.value) {
        searchNameInput.value.focus();
    }
}

watch(offset, async () => refresh());
watchDebounced(query.value, () => refresh(), { debounce: 600, maxWait: 1400 });
</script>

<template>
    <div class="py-2 pb-14">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="refresh()">
                        <div class="flex flex-row gap-4 mx-auto">
                            <div class="flex-1 form-control">
                                <label for="searchName" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.search') }}
                                    {{ $t('common.citizen', 1) }}
                                </label>
                                <div class="relative flex items-center mt-2">
                                    <input
                                        v-model="query.name"
                                        ref="searchNameInput"
                                        type="text"
                                        name="searchName"
                                        :placeholder="`${$t('common.citizen', 1)} ${$t('common.name')}`"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    />
                                </div>
                            </div>
                            <div class="flex-1 form-control">
                                <label for="dateofbirth" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.search') }}
                                    {{ $t('common.date_of_birth') }}
                                </label>
                                <div class="relative flex items-center mt-2">
                                    <input
                                        v-model="query.dateofbirth"
                                        type="text"
                                        name="dateofbirth"
                                        v-maska
                                        data-maska="##.##.####"
                                        :placeholder="`${$t('common.date_of_birth')} (DD.MM.YYYY)`"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    />
                                </div>
                            </div>
                            <div
                                class="flex-initial form-control"
                                v-if="can('CitizenStoreService.ListCitizens.Fields.UserProps.Wanted')"
                            >
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('components.citizens.citizens_list.only_wanted') }}
                                </label>
                                <div class="relative flex items-center mt-3">
                                    <Switch
                                        v-model="query.wanted"
                                        :class="[
                                            query.wanted ? 'bg-error-500' : 'bg-base-700',
                                            'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2',
                                        ]"
                                    >
                                        <span class="sr-only">
                                            {{ $t('components.citizens.citizens_list.only_wanted') }}
                                        </span>
                                        <span
                                            aria-hidden="true"
                                            :class="[
                                                query.wanted ? 'translate-x-5' : 'translate-x-0',
                                                'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-neutral ring-0 transition duration-200 ease-in-out',
                                            ]"
                                        />
                                    </Switch>
                                </div>
                            </div>
                        </div>
                        <Disclosure as="div" class="pt-2" v-slot="{ open }">
                            <DisclosureButton class="flex w-full items-start justify-between text-left text-neutral">
                                <span class="text-base-200 leading-7">{{ $t('common.advanced_search') }}</span>
                                <span class="ml-6 flex h-7 items-center">
                                    <ChevronDownIcon
                                        :class="[open ? 'upsidedown' : '', 'h-6 w-6 transition-transform']"
                                        aria-hidden="true"
                                    />
                                </span>
                            </DisclosureButton>
                            <DisclosurePanel class="mt-2 pr-4">
                                <div class="flex flex-row gap-2">
                                    <div
                                        class="flex-1 form-control"
                                        v-if="can('CitizenStoreService.ListCitizens.Fields.PhoneNumber')"
                                    >
                                        <label for="searchPhone" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.search') }}
                                            {{ $t('common.phone_number') }}
                                        </label>
                                        <div class="relative flex items-center mt-2">
                                            <input
                                                v-model="query.phoneNumber"
                                                type="tel"
                                                name="searchPhone"
                                                :placeholder="$t('common.phone_number')"
                                                class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            />
                                        </div>
                                    </div>
                                    <div class="flex-1 form-control">
                                        <label for="trafficPoints" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.search') }}
                                            {{ $t('common.traffic_infraction_points', 2) }}
                                        </label>
                                        <div class="relative flex items-center mt-2">
                                            <input
                                                v-model="query.trafficPoints"
                                                type="number"
                                                name="trafficPoints"
                                                :placeholder="`${$t('common.traffic_infraction_points')}`"
                                                class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            />
                                        </div>
                                    </div>
                                </div>
                            </DisclosurePanel>
                        </Disclosure>
                    </form>
                </div>
            </div>
            <div class="flow-root mt-2">
                <div class="mx-0 -my-2 overflow-x-auto">
                    <div class="inline-block min-w-full py-2 align-middle px-1">
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.citizen', 2)])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.citizen', 2)])"
                            :retry="refresh"
                        />
                        <DataNoDataBlock
                            v-else-if="data?.users.length === 0"
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
                                            {{ $t('common.name') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.job', 1) }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.sex') }}
                                        </th>
                                        <th
                                            v-if="can('CitizenStoreService.ListCitizens.Fields.PhoneNumber')"
                                            scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral"
                                        >
                                            {{ $t('common.phone_number') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.date_of_birth') }}
                                        </th>
                                        <th
                                            v-if="
                                                can('CitizenStoreService.ListCitizens.Fields.UserProps.TrafficInfractionPoints')
                                            "
                                            scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral"
                                        >
                                            {{ $t('common.traffic_infraction_points') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.height') }}
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
                                        v-for="user in data?.users"
                                        :key="user.userId"
                                        :user="user"
                                        class="transition-colors hover:bg-neutral/5"
                                    />
                                </tbody>
                                <thead>
                                    <tr>
                                        <th
                                            scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0"
                                        >
                                            {{ $t('common.name') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.job', 1) }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.sex') }}
                                        </th>
                                        <th
                                            v-if="can('CitizenStoreService.ListCitizens.Fields.PhoneNumber')"
                                            scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral"
                                        >
                                            {{ $t('common.phone_number') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.date_of_birth') }}
                                        </th>
                                        <th
                                            v-if="
                                                can('CitizenStoreService.ListCitizens.Fields.UserProps.TrafficInfractionPoints')
                                            "
                                            scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral"
                                        >
                                            {{ $t('common.traffic_infraction_points') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.height') }}
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

                            <TablePagination :pagination="data?.pagination" @offset-change="offset = $event" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
