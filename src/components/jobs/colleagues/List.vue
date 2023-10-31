<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { watchDebounced } from '@vueuse/core';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import { ColleaguesListResponse } from '~~/gen/ts/services/jobs/jobs';
import ListEntry from './ListEntry.vue';

const { $grpc } = useNuxtApp();

const query = ref<{ name: string }>({
    name: '',
});
const offset = ref(0n);

const { data, pending, refresh, error } = useLazyAsyncData(`jobs-colleagues-${offset.value}-${query.value.name}`, () =>
    listColleagues(),
);

async function listColleagues(): Promise<ColleaguesListResponse> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getJobsClient().colleaguesList({
                pagination: {
                    offset: offset.value,
                },
                searchName: query.value.name,
            });
            const { response } = await call;

            return res(response);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            throw e;
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
                                    {{ $t('common.colleague', 1) }}
                                </label>
                                <div class="relative flex items-center mt-2">
                                    <input
                                        v-model="query.name"
                                        ref="searchNameInput"
                                        type="text"
                                        name="searchName"
                                        :placeholder="$t('common.name')"
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
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.colleague', 2)])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.colleague', 2)])"
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
                                            {{ $t('common.rank', 1) }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.phone_number') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.date_of_birth') }}
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
                                            {{ $t('common.rank', 1) }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.phone_number') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.date_of_birth') }}
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
        </div>
    </div>
</template>
