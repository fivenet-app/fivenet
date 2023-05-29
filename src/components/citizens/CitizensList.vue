<script lang="ts" setup>
import { Switch } from '@headlessui/vue';
import { MagnifyingGlassIcon } from '@heroicons/vue/20/solid';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { watchDebounced } from '@vueuse/core';
import { ref } from 'vue';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import TablePagination from '~/components/partials/TablePagination.vue';
import { PaginationResponse } from '~~/gen/ts/resources/common/database/database';
import { User } from '~~/gen/ts/resources/users/users';
import { ListCitizensRequest } from '~~/gen/ts/services/citizenstore/citizenstore';
import CitizenListEntry from './CitizensListEntry.vue';

const { $grpc } = useNuxtApp();

const query = ref<{ name: string; phoneNumber?: string; wanted?: boolean }>({
    name: '',
});
const pagination = ref<PaginationResponse>();
const offset = ref(0);

const {
    data: users,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`citizens-${offset.value}-${query.value.name}-${query.value.wanted}-${query.value.phoneNumber}`, () =>
    listCitizens()
);

async function listCitizens(): Promise<Array<User>> {
    return new Promise(async (res, rej) => {
        try {
            const req: ListCitizensRequest = {
                pagination: {
                    offset: offset.value,
                },
                searchName: query.value.name,
            };
            if (query.value.wanted) {
                req.wanted = query.value.wanted;
            }
            if (query.value.phoneNumber) {
                req.phoneNumber = query.value.phoneNumber;
            }

            const call = $grpc.getCitizenStoreClient().listCitizens(req);
            const { response } = await call;

            pagination.value = response.pagination;
            return res(response.users);
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
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
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
                            <div class="flex-1 form-control" v-can="'CitizenStoreService.ListCitizens.Fields.PhoneNumber'">
                                <label for="searchPhone" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.search') }}
                                    {{ $t('common.phone') }}
                                </label>
                                <div class="relative flex items-center mt-2">
                                    <input
                                        v-model="query.phoneNumber"
                                        type="tel"
                                        name="searchPhone"
                                        :placeholder="`${$t('common.phone')} ${$t('common.number')}`"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    />
                                </div>
                            </div>
                            <div
                                class="flex-initial form-control"
                                v-can="'CitizenStoreService.ListCitizens.Fields.UserProps.Wanted'"
                            >
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral"
                                    >{{ $t('components.citizens.citizens_list.only_wanted') }}
                                </label>
                                <div class="relative flex items-center mt-3">
                                    <Switch
                                        v-model="query.wanted"
                                        :class="[
                                            query.wanted ? 'bg-error-500' : 'bg-base-700',
                                            'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2',
                                        ]"
                                    >
                                        <span class="sr-only">Wanted</span>
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
                    </form>
                </div>
            </div>
            <div class="flow-root mt-2">
                <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.citizen', 2)])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.citizen', 2)])"
                            :retry="refresh"
                        />
                        <button
                            v-else-if="users && users.length === 0"
                            type="button"
                            @click="focusSearch()"
                            class="relative block w-full p-12 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
                        >
                            <MagnifyingGlassIcon class="w-12 h-12 mx-auto text-neutral" />
                            <span class="block mt-2 text-sm font-semibold text-gray-300">
                                {{ $t('components.citizens.citizens_list.no_citizens') }}
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
                                            {{ $t('common.name') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.job', 1) }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.sex') }}
                                        </th>
                                        <th
                                            v-can="'CitizenStoreService.ListCitizens.Fields.PhoneNumber'"
                                            scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral"
                                        >
                                            {{ $t('common.phone') }}
                                            {{ $t('common.number') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.date_of_birth') }}
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
                                    <CitizenListEntry
                                        v-for="user in users"
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
                                            v-can="'CitizenStoreService.ListCitizens.Fields.PhoneNumber'"
                                            scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral"
                                        >
                                            {{ $t('common.phone') }}
                                            {{ $t('common.number') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.date_of_birth') }}
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

                            <TablePagination :pagination="pagination" @offset-change="offset = $event" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
