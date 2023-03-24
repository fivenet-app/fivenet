<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { User } from '@arpanet/gen/resources/users/users_pb';
import { OrderBy, PaginationRequest } from '@arpanet/gen/resources/common/database/database_pb';
import { watchDebounced } from '@vueuse/core'
import { getCitizenStoreClient } from '../../grpc/grpc';
import { FindUsersRequest } from '@arpanet/gen/services/citizenstore/citizenstore_pb';
import TablePagination from '../partials/TablePagination.vue';
import CitizenListEntry from './CitizensListEntry.vue';
import { Switch } from '@headlessui/vue';
import { MagnifyingGlassIcon } from '@heroicons/vue/20/solid';

const search = ref<{ name: string, wanted: boolean }>({ name: '', wanted: false });
const orderBys = ref<Array<OrderBy>>([]);
const offset = ref(0);
const totalCount = ref(0);
const listEnd = ref(0);
const users = ref<Array<User>>([]);

function findUsers(pos: number) {
    if (pos < 0) pos = 0;

    const req = new FindUsersRequest();
    req.setPagination((new PaginationRequest()).setOffset(pos));
    req.setSearchname(search.value.name);
    req.setWanted(search.value.wanted);
    req.setOrderbyList(orderBys.value);

    getCitizenStoreClient().
        findUsers(req, null).
        then((resp) => {
            const pag = resp.getPagination();
            if (pag !== undefined) {
                totalCount.value = pag.getTotalCount();
                offset.value = pag.getOffset();
                listEnd.value = pag.getEnd();
            }
            users.value = resp.getUsersList();
        });
}

function toggleOrderBy(column: string): void {
    const index = orderBys.value.findIndex((o: OrderBy) => {
        return o.getColumn() == column;
    });
    let orderBy: OrderBy;
    if (index > -1) {
        //@ts-ignore I just checked if it exists, so it should exist
        orderBy = orderBys.at(index);
        if (orderBy.getDesc()) {
            orderBys.value.splice(index);
        }
        else {
            orderBy.setDesc(true);
        }
    }
    else {
        orderBy = new OrderBy();
        orderBy.setColumn(column);
        orderBy.setDesc(false);
        orderBys.value.push(orderBy);
    }
    findUsers(offset.value);
}

const searchInput = ref<HTMLInputElement | null>(null);

function focusSearch(): void {
    if (searchInput.value) {
        searchInput.value.focus();
    }
}

watchDebounced(search.value, () => findUsers(0), { debounce: 650, maxWait: 1500 });

onMounted(() => {
    findUsers(0);
});
</script>

<template>
    <div class="py-2 sm:min-h-[67rem]">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="findUsers(0)">
                        <div class="flex flex-row gap-4 mx-auto">
                            <div class="flex-1 form-control">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral">Search</label>
                                <div class="relative mt-2 flex items-center">
                                    <input v-model="search.name" ref="focusSearch" type="text" name="search"
                                        id="search" placeholder="Citizen Name"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral shadow-sm placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                                </div>
                            </div>
                            <div class="flex-initial form-control">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral">Only
                                    Wanted</label>
                                <div class="relative mt-3 flex items-center">
                                    <Switch v-model="search.wanted"
                                        :class="[search.wanted ? 'bg-error-500' : 'bg-base-700', 'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2']">
                                        <span class="sr-only">Wanted</span>
                                        <span aria-hidden="true"
                                            :class="[search.wanted ? 'translate-x-5' : 'translate-x-0', 'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-neutral shadow ring-0 transition duration-200 ease-in-out']" />
                                    </Switch>
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
            <div class="mt-2 flow-root">
                <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                        <button v-if="users.length == 0" type="button" @click="focusSearch()"
                            class="relative block w-full rounded-lg border-2 border-dashed border-base-300 p-12 text-center hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2">
                            <MagnifyingGlassIcon class="text-neutral mx-auto h-12 w-12" />
                            <span class="mt-2 block text-sm font-semibold text-gray-300">Use the search field
                                above to search or update your query</span>
                        </button>
                        <div v-else>
                            <table class="min-w-full divide-y divide-base-600">
                                <thead>
                                    <tr>
                                        <th v-on:click="toggleOrderBy('firstname')" scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0">Name
                                        </th>
                                        <th v-on:click="toggleOrderBy('job')" scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">Job
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">Sex
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">Date
                                            of
                                            Birth
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Height
                                        </th>
                                        <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral">
                                            Action
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-base-800">
                                    <CitizenListEntry v-for="user in users" :key="user.getUserId()" :user="user" class="hover:bg-neutral/5 transition-colors" />
                                </tbody>
                                <thead>
                                    <tr>
                                        <th v-on:click="toggleOrderBy('firstname')" scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0">Name
                                        </th>
                                        <th v-on:click="toggleOrderBy('job')" scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">Job
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">Sex
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">Date
                                            of
                                            Birth
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Height
                                        </th>
                                        <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral">
                                            Action
                                        </th>
                                    </tr>
                                </thead>
                            </table>

                            <TablePagination :offset="offset" :entries="users.length" :end="listEnd" :total="totalCount"
                                :callback="findUsers" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
